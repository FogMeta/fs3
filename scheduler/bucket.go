package scheduler

import (
	"context"
	"errors"
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
		if err := ImportFromS3Bucket(restart); err != nil {
			logs.GetLogger().Error(err)
			return
		}
		restart = false
	})
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	c.Start()
}

func ImportFromS3Bucket(restart bool) (err error) {
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
	s3Client, err := minio.New(s3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3.AccessKeyID, s3.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	mc, err := minio.New("127.0.0.1:9000", &minio.Options{
		Creds: credentials.NewStaticV4(env.Get(config.EnvRootUser, ""), env.Get(config.EnvRootPassword, ""), ""),
	})
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	ctx := context.Background()
	ok, err := s3Client.BucketExists(ctx, s3.BucketName)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	if !ok {
		return errors.New("not found bucket")
	}
	// query target bucket objects
	mObjs := make(map[string]minio.ObjectInfo)
	for info := range mc.ListObjects(ctx, s3.TargetBucket, minio.ListObjectsOptions{
		Recursive: true,
	}) {
		mObjs[info.Key] = info
	}

	// query download bucket objects
	objects := s3Client.ListObjects(ctx, s3.BucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	var totalCnt, sucCnt int
	if s3.Status != statusS3Importing {
		if err = pdb.Model(s3).Updates(PsqlBucketImportS3{Status: statusS3Importing}).Error; err != nil {
			return
		}
	}
	for info := range objects {
		totalCnt++
		if co, ok := mObjs[info.Key]; ok && co.Size == info.Size {
			logs.GetLogger().Infof("%s %s already synced %s,skip", s3.BucketName, info.Key, s3.TargetBucket)
			continue
		}

		obj, err := s3Client.GetObject(ctx, s3.BucketName, info.Key, minio.GetObjectOptions{})
		if err != nil {
			logs.GetLogger().Errorf("%s %s get failed : %v", s3.BucketName, info.Key, err)
			continue
		}
		_, err = mc.PutObject(ctx, s3.TargetBucket, info.Key, obj, info.Size, minio.PutObjectOptions{})
		if err != nil {
			logs.GetLogger().Errorf("%s %s put failed : %v", s3.BucketName, info.Key, err)
			continue
		}
		logs.GetLogger().Infof("%s %s synced %s", s3.BucketName, info.Key, s3.TargetBucket)
		sucCnt++
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
