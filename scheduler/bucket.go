package scheduler

import (
	"gorm.io/gorm"
)

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
