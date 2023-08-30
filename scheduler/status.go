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

// restore to filesystem
const (
	StatusRebuildRestoreFailed = iota + 150
	StatusRebuildRestoreReady
	StatusRebuildRestoring
	StatusRebuildRestored
)

var RebuildStatusMsg = map[int]string{
	StatusRebuildReady:             "ready",
	StatusRebuildSuccess:           "completed",
	StatusRebuildCarDownloadFailed: "download failed",
	StatusRebuildCarDownloadReady:  "download ready",
	StatusRebuildCarDownloading:    "downloading",
	StatusRebuildCarDownloaded:     "downloaded",
	StatusRebuildRetrieveFailed:    "retrieve failed",
	StatusRebuildRetrieveReady:     "retrieve ready",
	StatusRebuildRetrieving:        "retrieving",
	StatusRebuildRetrieved:         "retrieved",
	StatusRebuildCarRestoreFailed:  "car restore failed",
	StatusRebuildCarRestoreReady:   "car restore ready",
	StatusRebuildCarRestoring:      "car restoring",
	StatusRebuildCarRestored:       "car restored",
	StatusRebuildStoreFailed:       "store failed",
	StatusRebuildStoreReady:        "store ready",
	StatusRebuildStoring:           "storing",
	StatusRebuildStored:            "stored",
	StatusRebuildRestoreFailed:     "restore failed",
	StatusRebuildRestoreReady:      "restore ready",
	StatusRebuildRestoring:         "restoring",
	StatusRebuildRestored:          "restored",
}

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
