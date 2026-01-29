# 站会提醒服务

## 项目简介

站会提醒服务是一个自动发送站会主持人提醒的工具，基于 Go 语言开发。它会根据配置文件和节假日信息，计算出当前站会的主持人和下一个站会的主持人，并通过 Slack 发送通知。

## 功能特性

- 自动计算当前站会主持人
- 自动计算下一个站会主持人
- 支持节假日自动跳过
- 通过 Slack 发送通知
- 自动获取和更新节假日信息
- 完善的错误处理和日志记录

## 目录结构

```
stand-meeting-notice/
├── pkg/                 # 可导出的公共包
│   ├── utils/           # 工具函数
│   ├── config/          # 配置管理
│   └── logger/          # 日志工具
├── config.json          # 配置文件
├── go.mod               # Go模块文件
├── main.go              # 主函数
└── README.md            # 项目说明
```

## 配置文件

配置文件 `config.json` 包含以下字段：

- `persons`：站会主持人列表
- `first_data`：第一个站会的日期，格式为 yyyyMMdd 文本转数字
- `first_index`：第一个站会的主持人在 persons 列表中的索引，从 0 开始
- `slack_url`：如果配置url，则会发送通知到该url，否则输出到终端
- `data_file`：节假日信息文件名，多个项目的时候可以以此区分

示例配置：

```json
{
  "persons": [
    "张三",
    "李四",
    "王五"
  ],
  "first_data": 20260102,
  "first_index": 1,
  "slack_url": "https://hooks.slack.com/triggers/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",
  "data_file": "data.json"
}
```

## 程序参数

- `-config`：指定配置文件路径，默认为`config.json`，如果有多个通知项目，可以指定不同的配置文件

## 运行方法

1. 确保已安装 Go 1.11 或更高版本
2. 克隆项目到本地
3. 在项目根目录下运行：

```bash
go run main.go
```

## 部署说明

可以将此服务部署为定时任务，每天运行一次，例如使用 crontab（Linux）或任务计划程序（Windows）。

## 依赖

- Go 标准库
- 无第三方依赖

## 许可证

本项目采用 MIT 许可证。
