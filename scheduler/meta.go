package scheduler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/minio/minio/logs"
	"github.com/minio/rpc/json2"
)

type metaClient struct {
	key    string
	token  string
	server string
}

func NewMetaClient(key, token, server string) *metaClient {
	return &metaClient{
		key:    key,
		token:  token,
		server: server,
	}
}

func (m *metaClient) Rebuild(id int64, object string) (data *RebuildResp, err error) {
	var req struct {
		DatasetID int64  `json:"dataset_id"`
		Object    string `json:"object"`
	}
	req.DatasetID = id
	req.Object = object
	data = &RebuildResp{}
	err = m.postReq("meta.DatasetRebuild", []interface{}{req}, data)
	return data, err
}

func (m *metaClient) Backup(name string, wallet string, data ...*FileData) (id int64, err error) {
	if len(data) == 0 {
		return 0, errors.New("data is required")
	}
	err = m.postReq("meta.Backup", []interface{}{name, data, wallet}, &id)
	return
}

func (m *metaClient) BackupDealStatus(id int64) (deal *DatasetDeal, err error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	deal = &DatasetDeal{}
	err = m.postReq("meta.DatasetDeal", []interface{}{id}, deal)
	return
}

func (m *metaClient) postReq(method string, args interface{}, dest interface{}) (err error) {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return errors.New("dest is not a pointer")
	}
	if dest == nil {
		return errors.New("dest not be nil")
	}
	b, err := json2.EncodeClientRequest(method, args)
	if err != nil {
		return
	}

	response, err := m.doPostReq("application/json;charset=utf-8", bytes.NewReader(b))
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	res := commonResp{
		Data: dest,
	}
	if err = json2.DecodeClientResponse(response.Body, &res); err != nil {
		return
	}
	if res.Code != "success" {
		return errors.New(res.Message)
	}
	return
}

func (m *metaClient) doPostReq(contentType string, body io.Reader) (resp *http.Response, err error) {
	server := m.server
	if !strings.HasSuffix(server, "/rpc/v0") {
		server, _ = url.JoinPath(server, "/rpc/v0")
	}
	client := new(http.Client)
	req, err := http.NewRequest(http.MethodPost, server, body)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("api-key", m.key)
	req.Header.Set("api-token", m.token)
	return client.Do(req)
}

type commonResp struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type FileData struct {
	SourceName  string `json:"source_name"`
	DataSize    int64  `json:"data_size"`
	IsDirectory bool   `json:"is_directory"`
	DownloadURL string `json:"download_url"`
}

type DatasetDeal struct {
	DatasetID     int64       `json:"dataset_id"`
	DatasetName   string      `json:"dataset_name"`
	DatasetStatus string      `json:"dataset_status"`
	Status        int         `json:"status"`
	PayloadCID    string      `json:"payload_cid"`
	PayloadURL    string      `json:"payload_url"`
	TaskID        string      `json:"task_id"`
	TaskName      string      `json:"task_name"`
	FileDescList  []*FileDesc `json:"file_desc_list"`
}

type FileDesc struct {
	Uuid           string
	SourceFileName string
	SourceFilePath string
	SourceFileMd5  string
	SourceFileSize int64
	CarFileName    string
	CarFilePath    string
	CarFileMd5     string
	CarFileUrl     string
	CarFileSize    int64
	PayloadCid     string
	PieceCid       string
	SourceId       *int
	Deals          []*DealInfo
}

type DealInfo struct {
	DealId     int
	DealCid    string
	MinerFid   string
	StartEpoch int
	Cost       string
}

type RebuildResp struct {
	Status     int      `json:"status"`
	PayloadCID string   `json:"payload_cid"`
	PayloadURL string   `json:"payload_url"`
	Providers  []string `json:"providers"`
	DueAt      int64    `json:"due_at"`
	CreatedAt  int64    `json:"created_at"`
}
