package channel

import (
	"stand-meeting-notice/src/utils"
)

type Channel interface {
	Send(meetingInfo *utils.MeetingInfo) error
}
