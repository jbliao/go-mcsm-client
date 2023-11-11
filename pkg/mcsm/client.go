package mcsm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"strings"
)

type Client struct {
	client   *http.Client
	Url      *urlpkg.URL
	ApiKey   string
	SvcUuid  string
	InstUuid string
}

type Queries map[string]string

type RequestBody map[string]any

type Result struct {
	Status int
	Time   uint64
	Data   json.RawMessage
}

func NewClient(urlStr, apiKey, serviceUuid, instanceUuid string) (*Client, error) {
	url, err := urlpkg.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:   http.DefaultClient,
		Url:      url,
		ApiKey:   apiKey,
		SvcUuid:  serviceUuid,
		InstUuid: instanceUuid,
	}, nil
}

func (c *Client) prepareRequest(ctx context.Context, method, path string, queries *Queries, body io.Reader) (req *http.Request, err error) {
	tempUrl := *c.Url

	q := tempUrl.Query()
	for key, value := range *queries {
		q.Set(key, value)
	}

	q.Set("apikey", c.ApiKey)

	tempUrl.RawQuery = q.Encode()
	tempUrl.Path = path

	req, err = http.NewRequestWithContext(ctx, method, tempUrl.String(), body)
	if err != nil {
		return
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return
}

func (c *Client) call(ctx context.Context, method, path string, queries *Queries, content *RequestBody) (out json.RawMessage, err error) {

	var reader io.Reader = http.NoBody
	if content != nil {
		var jsonBytes []byte
		jsonBytes, err = json.Marshal(content)
		if err != nil {
			return
		}

		reader = bytes.NewReader(jsonBytes)
	}

	req, err := c.prepareRequest(ctx, method, path, queries, reader)
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result Result
	err = json.Unmarshal(bodyByte, &result)
	if err != nil {
		return
	}

	out = result.Data
	if result.Status != 200 {
		err = errors.New(string(out))
	}

	return
}

type FileListItem struct {
	Name       string
	Size       int
	ModifyTime string `json:"time"`
	Permission uint   `json:"mode"`
	Type       uint
}

type FileListData struct {
	Items []FileListItem
}

func (c *Client) ListFile(ctx context.Context, path, fileName string) (out *FileListData, err error) {
	resp, err := c.call(ctx, "GET", "/api/files/list", &Queries{
		"remote_uuid": c.SvcUuid,
		"uuid":        c.InstUuid,
		"target":      path,
		"page":        "0",
		"page_size":   "30",
		"file_name":   fileName,
	}, nil)

	if err != nil {
		return nil, err
	}

	out = &FileListData{}
	err = json.Unmarshal(resp, out)
	return
}

func (c *Client) ZipFiles(ctx context.Context, zipName string, targetFiles ...string) (err error) {
	_, err = c.call(ctx, "POST", "/api/files/compress",
		&Queries{
			"remote_uuid": c.SvcUuid,
			"uuid":        c.InstUuid,
		}, &RequestBody{
			"type":    1,
			"source":  zipName,
			"targets": targetFiles,
			"code":    "utf-8",
		})

	return
}

type FileStatus struct {
	InstanceFileTask int
	GlobalFileTask   int
	Platform         string
	IsGlobalInstance bool
}

func (c *Client) FileStatus(ctx context.Context) (status *FileStatus, err error) {
	result, err := c.call(ctx, "GET", "/api/files/status",
		&Queries{
			"remote_uuid": c.SvcUuid,
			"uuid":        c.InstUuid,
		}, nil)
	if err != nil {
		return
	}

	status = &FileStatus{}
	err = json.Unmarshal(result, status)
	return
}

type FileDownloadData struct {
	Password string
	Addr     string
}

func (c *Client) DownloadFile(ctx context.Context, fileName string) (reader io.Reader, err error) {
	resp, err := c.call(ctx, "GET", "/api/files/download",
		&Queries{
			"remote_uuid": c.SvcUuid,
			"uuid":        c.InstUuid,
			"file_name":   fileName,
		}, nil)
	if err != nil {
		return
	}

	var info FileDownloadData
	if err = json.Unmarshal(resp, &info); err != nil {
		return
	}

	tmp := strings.Split(fileName, "/")
	fileName = tmp[len(tmp)-1]
	fileResp, err := http.Get(fmt.Sprintf("http://%s/download/%s/%s", info.Addr, info.Password, fileName))
	if fileResp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(fileResp.Body)
		err = errors.New(string(msg))
	}

	return fileResp.Body, err
}

func (c *Client) SendCommand(ctx context.Context, command string) (err error) {
	_, err = c.call(ctx, "GET", "/api/protected_instance/command",
		&Queries{
			"remote_uuid": c.SvcUuid,
			"uuid":        c.InstUuid,
			"command":     command,
		}, nil)

	return
}
