package scheduler

const (
	StatusRebuildReady = iota
	StatusRebuildSuccess
)

// download from car
const (
	StatusRebuildCarDownloadFailed = iota + 110
	StatusRebuildCarDownloadReady
	StatusRebuildCarDownloading
	StatusRebuildCarDownloaded
)

// download from provider
const (
	StatusRebuildRetrieveFailed = iota + 120
	StatusRebuildRetrieveReady
	StatusRebuildRetrieving
	StatusRebuildRetrieved
)

// restore from car
const (
	StatusRebuildCarRestoreFailed = iota + 130
	StatusRebuildCarRestoreReady
	StatusRebuildCarRestoring
	StatusRebuildCarRestored
)

// upload
const (
	StatusRebuildStoreFailed = iota + 140
	StatusRebuildStoreReady
	StatusRebuildStoring
	StatusRebuildStored
)

const (
	StatusRebuildDownloadFailed = iota + 150
	StatusRebuildDownloadReady
	StatusRebuildDownloading
	StatusRebuildDownloaded
)

// restore to filesystem
const (
	StatusRebuildRestoreFailed = iota + 160
	StatusRebuildRestoreReady
	StatusRebuildRestoring
	StatusRebuildRestored
)

func BackupStatusMsg(backup *PsqlBucketObjectBackup) string {
	status := backup.Status
	if CanRebuild(backup) {
		return "completed"
	}
	if status > 0 && status%10 == 0 {
		return "failed"
	}
	if status == 0 || status == 11 {
		return "ready"
	}
	return "backing up"
}

func RebuildStatusMsg(rebuild *PsqlBucketObjectRebuild) string {
	status := rebuild.Status
	if status > 0 && status%10 == 0 {
		return "failed"
	}
	if status == 1 || status == StatusRebuildRestored {
		return "completed"
	}
	if status == 0 {
		return "ready"
	}
	return "rebuilding"
}
