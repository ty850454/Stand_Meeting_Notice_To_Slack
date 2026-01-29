package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// HolidayAPIResponse 节假日API响应结构
type HolidayAPIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		List []struct {
			Date int64 `json:"date"`
		} `json:"list"`
		Page  int `json:"page"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"data"`
}

type HolidayMap map[int64]bool

// HolidaysData 节假日数据结构
type HolidaysData struct {
	Holidays    []int64 `json:"holidays"`    // 节假日列表
	LastUpdated int64   `json:"lastUpdated"` // 最后更新日期
}

func int64ArrToMap(holidays []int64) HolidayMap {
	holidayMap := make(HolidayMap, len(holidays))
	for _, holiday := range holidays {
		holidayMap[holiday] = true
	}
	return holidayMap
}

// FetchHolidaysFromApi 从API获取节假日列表
func FetchHolidaysFromApi(years ...int) ([]int64, error) {
	var builder strings.Builder
	length := len(years)
	for i, year := range years {
		builder.WriteString(strconv.Itoa(year))
		if i < length-1 {
			builder.WriteString(",")
		}
	}

	yearStr := builder.String()

	// 构建API URL
	url := fmt.Sprintf("https://api.apihubs.cn/holiday/get?field=date&year=%s&workday=2&size=3660", yearStr)

	// 发送 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var apiResp HolidayAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	// 检查响应状态
	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", apiResp.Msg)
	}

	// 提取节假日列表
	var holidays []int64
	for _, item := range apiResp.Data.List {
		holidays = append(holidays, item.Date)
	}

	return holidays, nil
}

// IsHoliday 检查指定日期是否是节假日
func IsHoliday(date int64, holidayMap HolidayMap) bool {
	_, ok := holidayMap[date]
	return ok
}
