package main

import (
	"flag"
	"os"

	"stand-meeting-notice/src/config"
	"stand-meeting-notice/src/logger"
	"stand-meeting-notice/src/utils"
)

const (
	// DefaultConfigFilePath 默认配置文件路径
	DefaultConfigFilePath = "config.json"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", DefaultConfigFilePath, "Config file path")
	testDate := flag.String("date", "", "Test date in yyyyMMdd format")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.LogError("failed to load config", err)
		os.Exit(1)
	}

	// 获取当前日期
	var currentDate int64
	var currentYear int

	if *testDate != "" {
		// 使用测试日期
		testTIme, err := utils.ParseStrToTime(*testDate)
		if err != nil {
			logger.LogError("invalid test date format", err)
			currentDate, currentYear = utils.GetCurrentIntDateAndYear()
		} else {
			currentDate = utils.TimeToIntDate(testTIme)
			currentYear = testTIme.Year()
			// 记录使用的测试日期
			logger.LogInfo("Using test date", "date", currentDate)
		}
	} else {
		// 使用当前日期
		currentDate, currentYear = utils.GetCurrentIntDateAndYear()
	}

	// 获取数据文件路径
	dataFilePath := utils.GetDefaultDataFilePath(cfg.DataFile)

	// 加载数据文件
	dataFile, err := utils.LoadDataFile(dataFilePath)

	// 检查是否需要更新数据文件
	if err != nil || utils.NeedUpdateDataFile(dataFile) {
		logger.LogInfo("Generating new data file...")
		// 生成新的数据文件
		newDataFile, err := utils.GenerateDataFile(cfg, currentDate, currentYear, dataFile)
		if err != nil {
			logger.LogError("failed to generate data file", err)
			os.Exit(1)
		}

		// 保存数据文件
		if err := utils.SaveDataFile(dataFilePath, newDataFile); err != nil {
			logger.LogError("failed to save data file", err)
			os.Exit(1)
		}

		logger.LogInfo("Data file generated successfully")
		dataFile = newDataFile
	}

	// 查找当前站会信息
	meetingInfo, found := utils.FindCurrentMeetingInfo(dataFile, currentDate)
	if !found {
		// 理论上不可能没有的
		logger.LogError("cannot find the data for the specified date", nil, "date", currentDate)
		os.Exit(1)
	}

	// 生成通知消息
	message := utils.GenerateNotificationMessage(meetingInfo)

	// 检查是否提供了 Slack Webhook URL
	if cfg.SlackUrl == "" {
		// 没有提供 Slack URL，只输出日志
		logger.LogInfo("Notification", "msg", message)
	} else {
		// 发送 Slack 通知
		if err := utils.SendSlackNotification(cfg.SlackUrl, message); err != nil {
			logger.LogError("failed to send slack notification", err)
			os.Exit(1)
		}

		// 记录成功信息
		logger.LogInfo("Meeting notification sent successfully", "mgs", message)
	}
}
