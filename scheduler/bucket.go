package scheduler

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio/internal/config"
	"github.com/minio/minio/logs"
	"github.com/minio/pkg/env"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

const (
	statusS3ImportFailed = iota - 1
	statusS3ImportReady
	statusS3Importing
	statusS3Imported
)

func ImportS3Scheduler() {
	c := cron.New()
	interval := "@every 1h"
	restart := true
	err := c.AddFunc(interval, func() {
		logs.GetLogger().Println("---------- import from s3 bucket scheduler is running at " + time.Now().Format("2006-01-02 15:04:05") + " ----------")
		if err := ImportFromS3Bucket(&restart); err != nil {
			logs.GetLogger().Error(err)
			return
		}
	})
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	c.Start()
}

func ImportFromS3Bucket(restartPtr *bool) (err error) {
	restart := false
	if restartPtr != nil {
		restart = *restartPtr
		if restart {
			defer func() {
				*restartPtr = false
			}()
		}
	}
	var s3 PsqlBucketImportS3
	mdb := pdb.Model(PsqlBucketImportS3{}).Where("status = ?", statusS3ImportReady)
	if restart {
		mdb = mdb.Or("status = ?", statusS3Importing).Order("status desc")
	}
	if err := mdb.First(&s3).Error; err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	// check import
	sc, err := minio.New(s3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3.AccessKeyID, s3.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	tc, err := minio.New("127.0.0.1:9000", &minio.Options{
		Creds: credentials.NewStaticV4(env.Get(config.EnvRootUser, ""), env.Get(config.EnvRootPassword, ""), ""),
	})
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	ctx := context.Background()
	ok, err := sc.BucketExists(ctx, s3.BucketName)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	if !ok {
		return errors.New("not found bucket")
	}
	// query target bucket objects
	mObjs := make(map[string]minio.ObjectInfo)
	for info := range tc.ListObjects(ctx, s3.TargetBucket, minio.ListObjectsOptions{
		Recursive: true,
	}) {
		mObjs[info.Key] = info
	}

	// query download bucket objects
	objects := sc.ListObjects(ctx, s3.BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	var totalCnt, sucCnt int
	if s3.Status != statusS3Importing {
		if err = pdb.Model(s3).Updates(PsqlBucketImportS3{Status: statusS3Importing}).Error; err != nil {
			return
		}
	}
	*restartPtr = false

	limit, err := env.GetInt(EnvSyncLimit, 0)
	if err != nil {
		logs.GetLogger().Error(err)
	}
	manager := NewSyncManager(limit, sc, tc, s3.BucketName, s3.TargetBucket)
	manager.Start(ctx)
	defer manager.Close()

	var objs []*minio.ObjectInfo
	for info := range objects {
		totalCnt++
		if info.Err != nil {
			logs.GetLogger().Infof("%s get object error: %v", s3.BucketName, info.Err)
			continue
		}
		if co, ok := mObjs[info.Key]; ok && co.Size == info.Size {
			logs.GetLogger().Infof("%s %s already synced %s,skip", s3.BucketName, info.Key, s3.TargetBucket)
			continue
		}

		obj := info
		objs = append(objs, &obj)
		manager.Send(&obj)
	}
	manager.WaitDone(totalCnt)

	for _, obj := range objs {
		if obj.Err == nil {
			sucCnt++
		}
	}

	if sucCnt == totalCnt {
		s3.Status = statusS3Imported
	} else if sucCnt == 0 {
		s3.Status = statusS3ImportFailed
	} else {
		return nil
	}
	return pdb.Model(s3).Updates(PsqlBucketImportS3{Status: s3.Status}).Error
}

type PsqlBucketImportS3 struct {
	gorm.Model
	AccessKeyID     string `gorm:"column:access_key_id"`
	SecretAccessKey string
	BucketName      string
	Endpoint        string
	Location        string
	TargetBucket    string
	UserAccessKey   string
	Status          int
	StatusMsg       string
	Progress        int
}

var pdb *gorm.DB

func Init() {
	var err error
	pdb, err = GetPsqlDb()
	if err != nil {
		panic(err)
	}
}

func GetPDB() *gorm.DB {
	return pdb
}

const (
	EnvSyncLimit = "SYNC_LIMIT"
)

type syncManager struct {
	limit        int
	sourceClient *minio.Client
	targetClient *minio.Client
	sourceBucket string
	targetBucket string
	ch           chan *minio.ObjectInfo
	exit         chan bool
	once         sync.Once
	cnt          atomic.Int32
}

func NewSyncManager(limit int, sourceClient, targetClient *minio.Client, sourceBucket, targetBucket string) *syncManager {
	if limit <= 0 {
		limit = 1
	}
	return &syncManager{
		limit:        limit,
		ch:           make(chan *minio.ObjectInfo, limit),
		sourceClient: sourceClient,
		targetClient: targetClient,
		sourceBucket: sourceBucket,
		targetBucket: targetBucket,
		exit:         make(chan bool),
	}
}

func (sm *syncManager) Send(info *minio.ObjectInfo) {
	sm.ch <- info
}

func (sm *syncManager) Start(ctx context.Context) {
	for i := 0; i < sm.limit; i++ {
		go sm.doSync(ctx)
	}
}

func (sm *syncManager) Close() {
	sm.once.Do(func() {
		close(sm.ch)
	})
}

func (sm *syncManager) Interrupt() {
	sm.exit <- true
}

func (sm *syncManager) WaitDone(total int) {
	for {
		if sm.cnt.Load() == int32(total) {
			return
		}
	}
}

func (sm *syncManager) doSync(ctx context.Context) {
	for {
		select {
		case info, ok := <-sm.ch:
			if info == nil && !ok {
				return
			}
			if info == nil {
				continue
			}
			obj, err := sm.sourceClient.GetObject(ctx, sm.sourceBucket, info.Key, minio.GetObjectOptions{})
			if err != nil {
				info.Err = err
				logs.GetLogger().Errorf("%s %s get failed : %v", sm.sourceBucket, info.Key, err)
				continue
			}
			_, err = sm.targetClient.PutObject(ctx, sm.targetBucket, info.Key, obj, info.Size, minio.PutObjectOptions{})
			if err != nil {
				info.Err = err
				logs.GetLogger().Errorf("%s %s put failed : %v", sm.targetBucket, info.Key, err)
				continue
			}
			logs.GetLogger().Infof("%s %s synced to %s", sm.sourceBucket, info.Key, sm.targetBucket)
		case <-ctx.Done():
			logs.GetLogger().Error(ctx.Err())
			return
		case <-sm.exit:
			logs.GetLogger().Info("interrupted")
			return
		}
	}
}

type PsqlBucketObjectBackup struct {
	gorm.Model
	UserAccessKey string
	BucketName    string
	ObjectName    string
	IsDir         bool
	Size          int64
	VersionID     string `gorm:"column:version_id"`
	DownloadURL   string `gorm:"column:download_url"`
	Filepath      string
	PayloadCID    string `gorm:"column:payload_cid"`
	PayloadURL    string `gorm:"column:payload_url"`
	MsID          int64  `gorm:"column:ms_id"`
	Status        int
	StatusMsg     string
}

type PsqlBucketObjectBackupSlice struct {
	ID         uint `gorm:"primarykey"`
	BackupID   uint `gorm:"column:backup_id"`
	FileName   string
	Size       int64
	PayloadCID string `gorm:"column:payload_cid"`
	PayloadURL string `gorm:"column:payload_url"`
	Status     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PsqlBucketObjectBackupSliceDeal struct {
	ID         uint   `gorm:"primarykey"`
	BackUpID   uint   `gorm:"column:backup_id"`
	PayloadCID string `gorm:"column:payload_cid"`
	MinerID    string `gorm:"column:miner_id"`
	DealID     int    `gorm:"column:deal_id"`
	DealCID    string `gorm:"column:deal_cid"`
	Cost       string
	Status     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func BackupSyncScheduler() {
	c := cron.New()
	interval := "@every 1m"
	err := c.AddFunc(interval, func() {
		logs.GetLogger().Println("---------- backup sync scheduler is running at " + time.Now().Format("2006-01-02 15:04:05") + " ----------")
		if err := BackupSync(); err != nil {
			logs.GetLogger().Error(err)
			return
		}
	})
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	c.Start()
}

func BackupSync() error {
	limit := 10
	statusDealActive := 45
	var id uint
	db := pdb.Model(PsqlBucketObjectBackup{}).Where("status < ?", statusDealActive)
	for {
		var backups []*PsqlBucketObjectBackup
		if id > 0 {
			db = pdb.Model(PsqlBucketObjectBackup{}).Where("id < ?", id).Where("status < ?", statusDealActive)
		}
		if err := db.Order("id desc").Limit(limit).Find(&backups).Error; err != nil {
			return err
		}
		for _, backup := range backups {
			if err := syncBackupInfo(backup); err != nil {
				logs.GetLogger().Error(err)
			}
			id = backup.ID
		}
		if len(backups) < limit {
			return nil
		}
	}
}

func syncBackupInfo(backup *PsqlBucketObjectBackup) error {
	if backup.MsID == 0 {
		return nil
	}
	client := NewMetaClient(env.Get("SWAN_KEY", ""), env.Get("SWAN_TOKEN", ""), env.Get("META_SERVER", ""))
	resp, err := client.BackupDealStatus(backup.MsID)
	if err != nil {
		return err
	}

	backup.Status = resp.Status
	if err = pdb.Model(backup).Updates(&PsqlBucketObjectBackup{
		PayloadCID: resp.PayloadCID,
		PayloadURL: resp.PayloadURL,
		Status:     resp.Status,
		StatusMsg:  resp.DatasetStatus,
	}).Error; err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	for _, fd := range resp.FileDescList {
		if err := syncBackupDetail(backup.ID, fd); err != nil {
			logs.GetLogger().Error(err)
		}
	}
	return nil
}

func syncBackupDetail(id uint, fd *FileDesc) error {
	bs := &PsqlBucketObjectBackupSlice{
		BackupID:   id,
		PayloadCID: fd.PayloadCid,
	}

	if err := pdb.Where(bs).First(bs).Error; err != nil {
		// insert
		bs.FileName = fd.CarFileName
		bs.PayloadURL = fd.CarFileUrl
		bs.Size = fd.SourceFileSize
		if err = pdb.Create(bs).Error; err != nil {
			return err
		}
	}
	// update
	if bs.PayloadURL != fd.CarFileUrl {
		if err := pdb.Model(bs).Updates(&PsqlBucketObjectBackupSlice{PayloadURL: fd.CarFileUrl}).Error; err != nil {
			return err
		}
	}

	// deals
	for _, deal := range fd.Deals {
		// query deal
		bsd := &PsqlBucketObjectBackupSliceDeal{
			BackUpID: bs.ID,
			MinerID:  deal.MinerFid,
		}
		if err := pdb.Where(bsd).First(bsd).Error; err != nil {
			bsd.PayloadCID = fd.PayloadCid
			bsd.DealID = deal.DealId
			bsd.DealCID = deal.DealCid
			bsd.Cost = deal.Cost
			if err = pdb.Create(bsd).Error; err != nil {
				logs.GetLogger().Error(err)
			}
		}

		// update
		if bsd.DealID != deal.DealId || bsd.DealCID != deal.DealCid {
			if err := pdb.Model(bsd).Updates(PsqlBucketObjectBackupSliceDeal{
				DealID:  deal.DealId,
				DealCID: deal.DealCid,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
