package scheduler

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
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

func (m *metaClient) Backup(name string, wallet string, data ...*FileData) (id int64, err error) {
	if len(data) == 0 {
		return 0, errors.New("data is required")
	}
	b, err := json2.EncodeClientRequest("meta.Backup", []interface{}{name, data, wallet})
	if err != nil {
		return
	}
	response, err := m.doPostReq("application/json;charset=utf-8", bytes.NewReader(b))
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	defer response.Body.Close()
	var res struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    int64  `json:"data"`
	}
	if err = json2.DecodeClientResponse(response.Body, &res); err != nil {
		return
	}
	if res.Code != "success" {
		return 0, errors.New(res.Message)
	}
	return res.Data, nil
}

func (m *metaClient) BackupDealStatus(id int64) (deal *DatasetDeal, err error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	b, err := json2.EncodeClientRequest("meta.DatasetDeal", []interface{}{id})
	if err != nil {
		return
	}
	response, err := m.doPostReq("application/json;charset=utf-8", bytes.NewReader(b))
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	defer response.Body.Close()
	var res struct {
		Code    string       `json:"code"`
		Message string       `json:"message"`
		Data    *DatasetDeal `json:"data,omitempty"`
	}
	if err = json2.DecodeClientResponse(response.Body, &res); err != nil {
		return
	}
	if res.Code != "success" {
		return nil, errors.New(res.Message)
	}
	return res.Data, nil
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
