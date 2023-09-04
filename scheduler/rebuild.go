package scheduler

import (
	"archive/zip"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio/internal/config"
	"github.com/minio/minio/logs"
	"github.com/minio/pkg/env"
	"github.com/robfig/cron"
)

var rebuildMap sync.Map

func RebuildSyncScheduler() {
	c := cron.New()
	interval := "@every 10m"
	err := c.AddFunc(interval, func() {
		logs.GetLogger().Println("---------- rebuild sync scheduler is running at " + time.Now().Format("2006-01-02 15:04:05") + " ----------")
		if err := RebuildSync(); err != nil {
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

func RebuildSync() error {
	limit := 10
	var id uint
	pdb := GetPDB()
	db := pdb.Model(PsqlBucketObjectRebuild{}).Where("status < ?", StatusRebuildRestored)
	for {
		var rebuilds []*PsqlBucketObjectRebuild
		if id > 0 {
			db = pdb.Model(PsqlBucketObjectRebuild{}).Where("id < ?", id).Where("status < ?", StatusRebuildRestored)
		}
		if err := db.Order("id desc").Limit(limit).Find(&rebuilds).Error; err != nil {
			return err
		}
		for _, rebuild := range rebuilds {
			if err := syncRebuildInfo(rebuild); err != nil {
				logs.GetLogger().Error(err)
			}
			id = rebuild.ID
		}
		if len(rebuilds) < limit {
			return nil
		}
	}
}

func syncRebuildInfo(rebuild *PsqlBucketObjectRebuild) (err error) {
	if _, ok := rebuildMap.Load(rebuild); ok {
		logs.GetLogger().Infof("%d is rebuilding", rebuild.ID)
		return nil
	}
	rebuildMap.Store(rebuild.ID, rebuild)
	var status int
	defer func() {
		rebuildMap.Delete(rebuild.ID)
		if status > 0 {
			updateRebuildStatus(rebuild, status)
		}
	}()
	mc := NewMetaClient(env.Get("SWAN_KEY", ""), env.Get("SWAN_TOKEN", ""), env.Get("META_SERVER", ""))
	data, err := mc.Rebuild(rebuild.MsID, rebuild.ObjectName)
	if err != nil {
		logs.GetLogger().Error("rebuild error:", err)
		return
	}
	logs.GetLogger().Infof("rebuild data resp: %#v", data)
	if data.Status == 1 {
		data.Status = StatusRebuildStored
	}
	rebuild.PayloadURL = data.PayloadURL
	rebuild.PayloadCID = data.PayloadCID
	rebuild.Status = data.Status
	rebuild.DueAt = data.DueAt
	pdb := GetPDB()
	if err = pdb.Model(rebuild).Updates(PsqlBucketObjectRebuild{
		Status:     data.Status,
		DueAt:      data.DueAt,
		PayloadCID: data.PayloadCID,
		PayloadURL: data.PayloadURL,
		Providers:  strings.Join(data.Providers, ","),
	}).Error; err != nil {
		logs.GetLogger().Error(err)
		return
	}
	if rebuild.Status < StatusRebuildStored {
		logs.GetLogger().Infof("rebuild status : %d, wait rebuild\n", rebuild.Status)
		return
	}

	logs.GetLogger().Info("rebuild status :", rebuild.Status)
	if data.PayloadURL == "" {
		return errors.New("invalid payload url")
	}

	// download file
	updateRebuildStatus(rebuild, StatusRebuildDownloading)
	targetPath, err := GetFile(rebuild)
	if err != nil {
		status = StatusRebuildDownloadFailed
		return
	}

	// restore to bucket
	updateRebuildStatus(rebuild, StatusRebuildRestoreReady)
	client, err := minio.New(env.Get("SERVER_ENDPOINT", "127.0.0.1:9000"), &minio.Options{
		Creds: credentials.NewStaticV4(env.Get(config.EnvRootUser, ""), env.Get(config.EnvRootPassword, ""), ""),
	})
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, rebuild.BucketName)
	if err != nil {
		return err
	}
	bucket := rebuild.BucketName
	if !exists || rebuild.ObjectName == "" {
		logs.GetLogger().Infof("bucket %s not exists, make it", bucket)
		// create bucket
		bucket = rebuild.BucketName + "-" + time.Now().Format("20060102150405")
		if err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			logs.GetLogger().Error("make new bucket failed: ", err)
			return
		}
	}

	// restore file
	updateRebuildStatus(rebuild, StatusRebuildRestoring)
	status = StatusRebuildRestored
	defer func() {
		if err != nil {
			status = StatusRebuildRestoreFailed
		}
	}()
	if !rebuild.IsDir {
		return restoreSingleFile(ctx, *client, bucket, targetPath, rebuild.ObjectName)
	}
	dirName := rebuild.ObjectName
	if exists && rebuild.ObjectName != "" {
		dirName += "-" + time.Now().Format("20060102150405")
	}
	err = filepath.WalkDir(targetPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			logs.GetLogger().Error(err)
			return err
		}
		if d.IsDir() {
			logs.GetLogger().Infof("path: %s is dir\n", path)
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		object := strings.TrimPrefix(path, targetPath+"/")
		if rebuild.ObjectName != "" {
			object = filepath.Join(dirName, object)
		}
		logs.GetLogger().Info("rebuild path: ", path, ", name: ", info.Name(), ", object: ", object, ", size: ", info.Size())
		_, err = client.PutObject(ctx, bucket, object, file, info.Size(), minio.PutObjectOptions{})
		if err != nil {
			logs.GetLogger().Errorf("put object %s failed: %v", object, err)
		}
		return err
	})
	return err
}

func restoreSingleFile(ctx context.Context, client minio.Client, bucket, path, object string) (err error) {
	if object == "" {
		return errors.New("invalid empty object")
	}
	name := filepath.Base(object)
	dir := filepath.Dir(object)
	ext := filepath.Ext(object)
	namePrefix := strings.TrimSuffix(name, ext)
	suffix := "-" + time.Now().Format("20060102150405") + ext
	object = namePrefix + suffix
	if dir != "." {
		object = filepath.Join(dir, object)
	}
	fi, err := os.Stat(path)
	if err != nil {
		return
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = client.PutObject(ctx, bucket, object, file, fi.Size(), minio.PutObjectOptions{})
	return err
}

func updateRebuildStatus(rebuild *PsqlBucketObjectRebuild, status int) error {
	return pdb.Model(rebuild).Updates(PsqlBucketObjectRebuild{Status: status}).Error
}

func putObject(ctx context.Context, client minio.Client, bucket, object, path string, size int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = client.PutObject(ctx, bucket, object, file, size, minio.PutObjectOptions{})
	return err
}

func GetFile(rebuild *PsqlBucketObjectRebuild) (path string, err error) {
	path, err = DownloadFile(rebuild)
	if err != nil {
		return
	}
	if !rebuild.IsDir {
		return
	}
	// unzip
	return unzipArchive(path)
}

func unzipArchive(zipPath string) (path string, err error) {
	name := strings.TrimSuffix(filepath.Base(zipPath), ".zip")
	path = filepath.Join(filepath.Dir(zipPath), name)
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return
	}
	os.MkdirAll(path, 0755)
	defer reader.Close()
	for _, file := range reader.File {
		info := file.FileInfo()
		dst := filepath.Join(path, info.Name())
		if info.IsDir() {
			os.MkdirAll(dst, 0755)
			continue
		}
		if err = copyZipFile(dst, file); err != nil {
			return
		}
	}
	return
}

func copyZipFile(path string, file *zip.File) error {
	dst, err := createFile(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	zr, err := file.Open()
	if err != nil {
		return err
	}
	defer zr.Close()
	_, err = io.Copy(dst, zr)
	if err != nil {
		return err
	}
	return nil
}

func DownloadFile(rebuild *PsqlBucketObjectRebuild) (path string, err error) {
	downloadURL := rebuild.PayloadURL
	logs.GetLogger().Infof("%d %s %s start downloading url:%s \n", rebuild.ID, rebuild.BucketName, rebuild.ObjectName, downloadURL)
	client := &http.Client{}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// path: /download/ms_id/name
	path = filepath.Join(env.Get("REBUILD_PATH", ""), strconv.Itoa(int(rebuild.MsID)), rebuild.BucketName, rebuild.ObjectName)
	if rebuild.IsDir {
		path += ".zip"
	}
	logs.GetLogger().Infof("%d %s %s download path: %s\n", rebuild.ID, rebuild.BucketName, rebuild.ObjectName, path)

	fi, err := os.Stat(path)
	if err == nil {
		size, _ := DownloadLinkSize(downloadURL)
		if size > 0 && fi.Size() == size {
			logs.GetLogger().Infof("%d %s %s  already downloaded, path: %s, no need download\n", rebuild.ID, rebuild.BucketName, rebuild.ObjectName, path)
			return
		}
	}

	file, err := createFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
	}
	logs.GetLogger().Infof("%d %s %s downloaded to path: %s\n", rebuild.ID, rebuild.BucketName, rebuild.ObjectName, path)
	return
}

func createFile(name string) (*os.File, error) {
	dir := filepath.Dir(name)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func DownloadLinkSize(url string) (size int64, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return size, err
	}
	contentLength := resp.Header.Get("Content-Length")
	return strconv.ParseInt(contentLength, 10, 64)
}
