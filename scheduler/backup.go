package scheduler

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/minio/minio/internal/auth"
	"github.com/minio/minio/internal/config"
	"github.com/minio/minio/internal/jwt"
	"github.com/minio/minio/logs"
	"github.com/minio/pkg/env"
	"github.com/robfig/cron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	MicroSecondPerDay          = 86400000000
	MicroSecondPerMinute       = 60000000
	FS3SourceId                = 3
	StatusBackupTaskCreated    = "Created"
	StatusRebuildTaskCompleted = "Completed"
	StatusBackupPlanEnabled    = "Enabled"
)

const (
	StatusCodeCreated = iota
	StatusCodeRunning
	StatusCodeCompleted
)

var StatusPlanMsg = map[int]string{
	StatusCodeCreated:   "created",
	StatusCodeRunning:   "running",
	StatusCodeCompleted: "completed",
}

const (
	StatusDisabled = iota
	StatusEnabled
)

const (
	StatusMsgEnabled  = "Enabled"
	StatusMsgDisabled = "Disabled"
)

var StatusMsgAble = map[int]string{
	StatusDisabled: StatusMsgDisabled,
	StatusEnabled:  StatusMsgEnabled,
}

var StatusCodeAble = map[string]int{
	StatusMsgDisabled: StatusDisabled,
	StatusMsgEnabled:  StatusEnabled,
}

func BackupScheduler() {
	c := cron.New()
	interval := "@every 2m"
	err := c.AddFunc(interval, func() {
		logs.GetLogger().Println("++++++++++ backup scheduler is running at " + time.Now().Format("2006-01-02 15:04:05") + " ++++++++++")
		err := BackupBucketScheduler()
		if err != nil {
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

func BackupBucketScheduler() error {
	//get all the running backup plans from db
	var plans []*PsqlBucketBackupPlan
	plan := PsqlBucketBackupPlan{Status: 1}
	if err := pdb.Debug().Model(plan).Where(plan).Find(&plans).Error; err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	for _, plan := range plans {
		interval := plan.Interval * 3600 * 24
		if plan.LastAt >= time.Now().Unix()-int64(interval) {
			logs.GetLogger().Infof("plan %s already in backup progress", plan.Name)
			continue
		}
		if err := backupPlan(plan); err != nil {
			logs.GetLogger().Error(err)
			continue
		}
	}
	return nil
}

func backupPlan(plan *PsqlBucketBackupPlan) error {
	if err := pdb.First(plan, plan.ID).Error; err != nil {
		return err
	}
	// update backup time
	pdb.Model(plan).Updates(PsqlBucketBackupPlan{LastAt: time.Now().Unix()})
	buckets := strings.Split(plan.Bucket, ",")
	for _, bucket := range buckets {
		if err := backup(bucket, "", plan); err != nil {
			logs.GetLogger().Error(err)
		}
	}
	return nil
}

func GetPsqlDb() (*gorm.DB, error) {
	host := config.GetUserConfig().PsqlHost
	user := config.GetUserConfig().PsqlUser
	password := config.GetUserConfig().PsqlPassword
	dbname := config.GetUserConfig().PsqlDbname
	port := config.GetUserConfig().PsqlPort
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return db, err
}

type PsqlBucketBackupPlan struct {
	ID             int `gorm:"primary_key"`
	UserAccessKey  string
	Name           string
	Bucket         string
	Interval       int
	ProviderRegion string
	Duration       int
	VerifiedDeal   bool
	FastRetrieval  bool
	Status         int
	LastAt         int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func backup(bucket, object string, plan *PsqlBucketBackupPlan) (err error) {
	filWallet := config.GetUserConfig().Fs3WalletAddress
	if filWallet == "" {
		return errors.New("Please provide a wallet address for sending deals")
	}
	fs3VolumeAddress := config.GetUserConfig().Fs3VolumeAddress
	sourceFilePath := filepath.Join(fs3VolumeAddress, bucket, object)
	fi, err := os.Stat(sourceFilePath)
	if err != nil {
		return
	}
	token, err := authenticateJWTUsers(env.Get(config.EnvRootUser, ""), env.Get(config.EnvRootPassword, ""), 24*time.Hour*7)
	if err != nil {
		logs.GetLogger().Error("authenticateURL error:", err)
		return
	}

	downloadURL, _ := url.JoinPath(env.Get("HOST_NAME", ""), "", "download", bucket, object)
	values := url.Values{}
	values.Set("token", token)
	if fi.IsDir() {
		values.Set("format", "zip")
	}
	downloadURL += "?" + values.Encode()
	client := NewMetaClient(env.Get("SWAN_KEY", ""), env.Get("SWAN_TOKEN", ""), env.Get("META_SERVER", ""))
	id, err := client.Backup(fmt.Sprintf("%s-%s-%s", bucket, fi.Name(), time.Now().Format("20060102150405")), filWallet, &FileData{
		SourceName:  fi.Name(),
		DataSize:    fi.Size(),
		IsDirectory: fi.IsDir(),
		DownloadURL: downloadURL,
	})
	if err != nil {
		logs.GetLogger().Error("backup error:", err)
		return
	}
	// save data
	return pdb.Create(&PsqlBucketObjectBackup{
		UserAccessKey:  plan.UserAccessKey,
		BucketName:     bucket,
		ObjectName:     object,
		IsDir:          fi.IsDir(),
		Size:           fi.Size(),
		DownloadURL:    downloadURL,
		Filepath:       sourceFilePath,
		MsID:           id,
		ProviderRegion: plan.ProviderRegion,
		Duration:       plan.Duration,
		VerifiedDeal:   plan.VerifiedDeal,
		FastRetrieval:  plan.FastRetrieval,
	}).Error
}

func authenticateJWTUsers(accessKey, secretKey string, expiry time.Duration) (string, error) {
	passedCredential, err := auth.CreateCredentials(accessKey, secretKey)
	if err != nil {
		return "", err
	}
	expiresAt := time.Now().UTC().Add(expiry)
	return authenticateJWTUsersWithCredentials(passedCredential, expiresAt)
}

func authenticateJWTUsersWithCredentials(credentials auth.Credentials, expiresAt time.Time) (string, error) {
	claims := jwt.NewMapClaims()
	claims.SetExpiry(expiresAt)
	claims.SetAccessKey(credentials.AccessKey)

	jwt := jwtgo.NewWithClaims(jwtgo.SigningMethodHS512, claims)
	return jwt.SignedString([]byte(credentials.SecretKey))
}

func BackupStatusParse(status int) int {
	if status >= 45 {
		return 1
	}
	if status > 0 && status%10 == 0 {
		return -1
	}
	return 0
}

func RebuildStatusParse(status int) int {
	if status == 153 || status == 1 {
		return 1
	}
	if status > 0 && status%10 == 0 {
		return -1
	}
	return 0
}
