package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"stand-meeting-notice/src/logger"
	"strconv"
	"time"

	"stand-meeting-notice/src/config"
)

// DataFile 数据文件结构
type DataFile struct {
	Date        map[string]*DateInfo `json:"date"`        // 站会日期信息
	LastUpdated int64                `json:"lastUpdated"` // 最后更新日期
	Persons     []string             `json:"persons"`     // 站会主持人列表
}

// DateInfo 日期信息结构
type DateInfo struct {
	Person int   `json:"person"` // 站会主持人索引，-1表示节假日
	Next   int64 `json:"next"`   // 下一个站会日期
}

// LoadDataFile 加载数据文件
func LoadDataFile(dataFilePath string) (*DataFile, error) {
	// 检查文件是否存在
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("data file does not exist")
	}

	// 读取文件内容
	data, err := os.ReadFile(dataFilePath)
	if err != nil {
		return nil, err
	}

	// 解析 JSON
	var dataFile DataFile
	if err := json.Unmarshal(data, &dataFile); err != nil {
		return nil, err
	}

	return &dataFile, nil
}

// SaveDataFile 保存数据文件
func SaveDataFile(dataFilePath string, dataFile *DataFile) error {
	// 序列化 JSON
	data, err := json.Marshal(dataFile)
	if err != nil {
		return err
	}

	// 保存到文件
	return os.WriteFile(dataFilePath, data, 0644)
}

// NeedUpdateDataFile 检查是否需要更新数据文件
func NeedUpdateDataFile(dataFile *DataFile) bool {
	// 如果数据文件不存在，需要更新
	if dataFile == nil {
		return true
	}

	// 检查最后更新日期是否是一个月前
	lastUpdatedTime := IntDateToTime(dataFile.LastUpdated)
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	return !lastUpdatedTime.After(oneMonthAgo)
}

// GetDefaultDataFilePath 获取默认数据文件路径
func GetDefaultDataFilePath(cfgDataFile string) string {
	// 如果配置中指定了数据文件路径，使用配置中的路径
	if cfgDataFile != "" {
		return cfgDataFile
	}

	// 默认使用 data.json
	return filepath.Join(".", "data.json")
}

// FindCurrentMeetingInfo 查找当前站会信息
func FindCurrentMeetingInfo(dataFile *DataFile, currentDate int64) (*MeetingInfo, bool) {
	// 转换当前日期为字符串
	currentDateStr := fmt.Sprintf("%d", currentDate)

	// 查找当前日期的信息
	dateInfo, exists := dataFile.Date[currentDateStr]
	if !exists {
		return nil, false
	}

	// 如果是节假日，尝试通过next找到下一个工作日
	if dateInfo.Person == -1 {
		nextDateStr := fmt.Sprintf("%d", dateInfo.Next)
		nextDateInfo, nextExists := dataFile.Date[nextDateStr]
		if !nextExists || nextDateInfo.Person == -1 {
			return nil, false
		}

		// 获取当前主持人（应该是上一个工作日的主持人）
		// 这里需要遍历找到上一个工作日
		var prevDateInfo *DateInfo
		for _, dInfo := range dataFile.Date {
			if dInfo.Person != -1 {
				prevDateInfo = dInfo
				break
			}
		}

		if prevDateInfo == nil || prevDateInfo.Person == -1 {
			return nil, false
		}

		// 获取下一个主持人
		nextPerson := dataFile.Persons[nextDateInfo.Person]

		return &MeetingInfo{
			CurrentPerson: dataFile.Persons[prevDateInfo.Person],
			NextDate:      dateInfo.Next,
			NextPerson:    nextPerson,
		}, true
	}

	// 正常工作日
	currentPerson := dataFile.Persons[dateInfo.Person]

	// 查找下一个主持人
	nextDateStr := fmt.Sprintf("%d", dateInfo.Next)
	nextDateInfo, nextExists := dataFile.Date[nextDateStr]
	if !nextExists || nextDateInfo.Person == -1 {
		// 找不到下一个工作日，只返回当前主持人
		return &MeetingInfo{
			CurrentPerson: currentPerson,
			NextDate:      0,
			NextPerson:    "",
		}, true
	}

	nextPerson := dataFile.Persons[nextDateInfo.Person]

	return &MeetingInfo{
		CurrentPerson: currentPerson,
		NextDate:      dateInfo.Next,
		NextPerson:    nextPerson,
	}, true
}

// GenerateDataFile 生成数据文件
func GenerateDataFile(cfg *config.Config, currentDate int64, currentYear int, oldDateFile *DataFile) (*DataFile, error) {

	// 确定开始日期
	var startDate int64
	var startPerson int
	var persons []string

	if oldDateFile != nil {
		info, exists := oldDateFile.Date[strconv.FormatInt(currentDate, 10)]
		if exists {
			if info.Person == -1 {
				// 节假日，获取下一个工作日
				info, exists = oldDateFile.Date[strconv.FormatInt(currentDate, 10)]
				if exists {
					// 这次不可能是节假日了
					startDate = currentDate
					startPerson = info.Person
					persons = oldDateFile.Persons
				}
			} else {
				// 工作日
				startDate = currentDate
				startPerson = info.Person
				persons = oldDateFile.Persons
			}
		}
	}

	if startDate == 0 {
		startDate = cfg.FirstData
		startPerson = cfg.FirstIndex
		persons = cfg.Persons
	}

	startTime := IntDateToTime(startDate)

	// 确保开始日期是当年或上一年的日期
	startYear := startTime.Year()
	if startYear < currentYear-1 || startYear > currentYear {
		// 如果开始日期不是当年或上一年，调整为当年的1月1日
		logger.LogError("Incorrect start date", nil)
		os.Exit(1)
	}

	// 获取节假日数据
	var holidays []int64
	var holidayMap HolidayMap
	var err error

	if startYear < currentYear {
		holidays, err = FetchHolidaysFromApi(currentYear-1, currentYear, currentYear+1)
	} else {
		holidays, err = FetchHolidaysFromApi(currentYear, currentYear+1)
	}
	if err != nil {
		return nil, err
	}
	holidayMap = int64ArrToMap(holidays)
	maxHoliday := holidays[0]
	for _, holiday := range holidays {
		if holiday > maxHoliday {
			maxHoliday = holiday
		}
	}
	// 确定结束日期
	endDate := maxHoliday
	endTime := IntDateToTime(endDate)

	// 遍历从开始日期到结束日期
	currentPerson := startPerson
	lasts := make([]*DateInfo, 0, 10)
	maxPerson := len(persons) - 1
	// 初始化数据文件结构
	dataFile := &DataFile{
		Date:        make(map[string]*DateInfo),
		LastUpdated: currentDate,
		Persons:     persons,
	}

	for d := startTime; d.Before(endTime); d = d.AddDate(0, 0, 1) {
		dateInt := TimeToIntDate(d)
		dateStr := strconv.FormatInt(dateInt, 10)

		temp := &DateInfo{
			Person: -1,
			Next:   -1,
		}
		dataFile.Date[dateStr] = temp

		if IsHoliday(dateInt, holidayMap) {
			// 节假日
			lasts = append(lasts, temp)
			continue
		}

		// 工作日
		temp.Person = currentPerson
		currentPerson += 1
		if currentPerson > maxPerson {
			currentPerson = 0
		}

		if len(lasts) > 0 {
			for _, last := range lasts {
				last.Next = dateInt
			}
		}
		lasts = []*DateInfo{temp}
	}

	return dataFile, nil
}
