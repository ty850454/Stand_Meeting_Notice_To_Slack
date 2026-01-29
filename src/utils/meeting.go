package utils

// MeetingInfo 站会信息结构
type MeetingInfo struct {
	CurrentPerson string // 当前站会主持人
	NextDate      int64  // 下一个站会日期
	NextPerson    string // 下一个站会主持人
}
