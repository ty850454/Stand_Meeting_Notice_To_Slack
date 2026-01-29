package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SlackNotification Slack通知结构
type SlackNotification struct {
	Msg string `json:"msg"`
}

// SendSlackNotification 发送Slack通知
func SendSlackNotification(webhookURL string, message string) error {
	// 创建通知对象
	notification := SlackNotification{
		Msg: message,
	}

	// 序列化通知对象为 JSON
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack notification failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// GenerateNotificationMessage 生成通知消息
func GenerateNotificationMessage(meetingInfo *MeetingInfo) string {
	if meetingInfo.NextDate != 0 && meetingInfo.NextPerson != "" {
		// 有下一个站会主持人
		nextDateStr := FormatIntDate(meetingInfo.NextDate)
		return fmt.Sprintf("今天由 %s 主持站会，下一次 %s 由 %s 主持", meetingInfo.CurrentPerson, nextDateStr, meetingInfo.NextPerson)
	}

	// 没有下一个站会主持人
	return fmt.Sprintf("今天由 %s 主持站会", meetingInfo.CurrentPerson)
}
