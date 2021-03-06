package ding_talk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DingTalkClient struct {
	gatewayUrl string
}

const (
	Text       = "text"       // 文本
	Link       = "link"       // 链接
	Markdown   = "markdown"   // https://open-doc.dingtalk.com/microapp/serverapi2/qf2nxq#-6
	ActionCard = "actionCard" // https://open-doc.dingtalk.com/microapp/serverapi2/qf2nxq#-7
	FeedCard   = "feedCard"   // https://open-doc.dingtalk.com/microapp/serverapi2/qf2nxq#-9
)

type RobotSendResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func NewClient(url string) *DingTalkClient {
	return &DingTalkClient{gatewayUrl: url}
}

func NewClients(urls []string) []DingTalkClient {
	clients := make([]DingTalkClient, 0, len(urls))
	for _, url := range urls {
		clients = append(clients, DingTalkClient{gatewayUrl: url})
	}
	return clients
}

func (d *DingTalkClient) Execute(params interface{}) (RobotSendResponse, error) {
	var (
		response RobotSendResponse
	)
	switch params.(type) {
	case TextMessage:
	case LinkMessage:
	case MarkdownMessage:
	case SingleActionCardMessage:
	case ActionCardMessage:
	case FeedCardMessage:
	default:
		log.Fatalf("叮叮通知参数验证失败：params: %+v", params)
	}
	// json Marshal
	data, err := json.Marshal(params)
	if err != nil {
		return response, errors.New(fmt.Sprintf("json序列化失败 %+v", err))
	}
	body := bytes.NewReader(data)

	// request ding ding
	request := &http.Request{}
	request, err = http.NewRequest(http.MethodPost, d.gatewayUrl, body)
	if err != nil {
		return response, errors.New(fmt.Sprintf("new request 失败 %+v", err))
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := http.Client{Timeout: 6 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		return response, errors.New(fmt.Sprintf("发送通知失败 %+v", err))
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, errors.New(fmt.Sprintf("读取response值失败 %+v", err))
	}

	err = json.Unmarshal(respData, &response)
	if err != nil {
		return response, errors.New(fmt.Sprintf("json反序列化response失败 data:%+v, err:%+v", string(respData), err))
	}
	return response, err
}
