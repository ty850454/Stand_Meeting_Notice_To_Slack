package channel

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stand-meeting-notice/src/config"
	"stand-meeting-notice/src/utils"
	"time"
)

type FeiShuChannel struct {
	Webhook     string
	Sign        bool
	Secret      string
	UserMapping *map[string]string
}

func NewFeiShuChannel(cfg *config.FeiShuConfig) *FeiShuChannel {
	return &FeiShuChannel{
		Webhook:     cfg.Webhook,
		Sign:        cfg.Sign,
		Secret:      cfg.Secret,
		UserMapping: cfg.UserMapping,
	}
}

func genSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

type FeiShuNotification struct {
	Timestamp string `json:"timestamp,omitempty"`
	Sign      string `json:"sign,omitempty"`
	MsgType   string `json:"msg_type"`
	Content   struct {
		Text string `json:"text"`
	} `json:"content"`
}

type FeiShuResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func sendFeiShuNotification(webhookURL string, message string, sign bool, secret string) error {
	println(message)
	notification := FeiShuNotification{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{Text: message},
	}

	if sign {
		timestamp := time.Now().Unix()
		signStr, err := genSign(secret, timestamp)
		if err != nil {
			return err
		}
		notification.Timestamp = fmt.Sprintf("%v", timestamp)
		notification.Sign = signStr
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
		return fmt.Errorf("fei shu notification failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var result FeiShuResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	if result.Code != 0 {
		return fmt.Errorf("fei shu notification failed with code: %d, msg: %s", result.Code, result.Msg)
	}

	return nil
}

func (channel *FeiShuChannel) Send(meetingInfo *utils.MeetingInfo) error {

	var currentPerson string
	if channel.UserMapping != nil {
		if userId, ok := (*channel.UserMapping)[meetingInfo.CurrentPerson]; ok {
			currentPerson = fmt.Sprintf("<at user_id=\"%s\">%s</at>", userId, meetingInfo.CurrentPerson)
		} else {
			currentPerson = meetingInfo.CurrentPerson
		}
	}

	var message string
	if meetingInfo.NextDate != 0 && meetingInfo.NextPerson != "" {
		// 有下一个站会主持人
		nextDateStr := utils.FormatIntDate(meetingInfo.NextDate)
		message = fmt.Sprintf("今天由 %s 主持站会，下一次 %s 由 %s 主持", currentPerson, nextDateStr, meetingInfo.NextPerson)
	} else {
		message = fmt.Sprintf("今天由 %s 主持站会", currentPerson)
	}

	return sendFeiShuNotification(channel.Webhook, message, channel.Sign, channel.Secret)
}
