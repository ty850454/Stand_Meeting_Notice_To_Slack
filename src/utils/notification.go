package utils

import (
	"fmt"
)

// GenerateSimpleMessage 生成通知消息
func GenerateSimpleMessage(meetingInfo *MeetingInfo) string {
	if meetingInfo.NextDate != 0 && meetingInfo.NextPerson != "" {
		// 有下一个站会主持人
		nextDateStr := FormatIntDate(meetingInfo.NextDate)
		return fmt.Sprintf("今天由 %s 主持站会，下一次 %s 由 %s 主持", meetingInfo.CurrentPerson, nextDateStr, meetingInfo.NextPerson)
	}

	// 没有下一个站会主持人
	return fmt.Sprintf("今天由 %s 主持站会", meetingInfo.CurrentPerson)
}
