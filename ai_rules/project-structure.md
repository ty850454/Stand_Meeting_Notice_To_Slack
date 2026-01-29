# 项目结构规范

## Go 项目结构

### 目录结构
project-name/
├── cmd/                 # 应用入口
│   └── server/          # 服务器入口
│       └── main.go      # 主函数
├── internal/            # 内部包
│   ├── module/          # 模块构建器
│   │   ├── nn/          # NN业务模块
│   │   ├── seckill/     # 秒杀业务模块
│   │   └── common/      # 通用模块
│   ├── domain/          # 领域对象
│   │   ├── module/      # 模块VO
│   │   └── [业务名]/     # 业务领域对象（BO、DTO）
│   ├── dataservice/     # 数据服务
│   │   └── impl/        # 数据服务实现
│   └── provider/        # 外部服务提供者
├── pkg/                 # 可导出的公共包
│   ├── utils/           # 工具函数
│   ├── config/          # 配置管理
│   └── logger/          # 日志工具
├── configs/             # 配置文件
├── scripts/             # 脚本文件
├── go.mod               # Go模块文件
├── go.sum               # 依赖校验文件
└── README.md            # 项目说明

### 包结构规范
- **cmd/**：应用入口点，包含main函数
- **internal/**：内部包，不对外暴露
  - **module/**：模块构建器，按业务划分
  - **domain/**：领域对象，包含VO、BO、DTO等
  - **dataservice/**：数据服务层
  - **provider/**：外部服务提供者
- **pkg/**：可被其他项目导入的公共包
- **configs/**：配置文件
- **scripts/**：构建、部署等脚本

### 命名规范
- **模块构建器**：[业务名]ModuleBuilder（如 NnFeedsModuleBuilder）
- **数据服务**：[业务名]DataService（如 NnRedPacketDataService）
- **模块VO**：[业务名]VO（如 NnRedPacketVO）
- **业务BO**：[业务名]BO（如 NnRoundFeatureBO）
- **包名**：使用小写字母，简短且有意义（如 `module`、`domain`）
- **文件名**：使用小写字母，单词间用下划线分隔（如 `http_server.go`）
- **目录名**：使用小写字母，单词间用下划线分隔（如 `data_service`）

### 模块划分
- **按业务领域划分**：每个业务领域有独立的模块目录
- **按功能分层**：清晰的分层结构，如模块层、领域层、数据服务层
- **公共代码提取**：将公共功能提取到 `pkg` 或 `common` 目录

### 依赖管理
- 使用 Go Modules 管理依赖
- 在 `go.mod` 文件中指定模块路径和依赖版本
- 依赖版本应明确指定，避免使用 `latest`
- 定期运行 `go mod tidy` 清理未使用的依赖

### 构建与部署
- 使用 `scripts` 目录存放构建和部署脚本
- 提供标准化的构建流程
- 支持容器化部署（如 Docker）

### 最佳实践
- **单一职责**：每个包和文件应只负责一个功能
- **接口分离**：使用接口定义依赖关系
- **依赖注入**：通过构造函数注入依赖
- **错误处理**：统一的错误处理机制
- **测试覆盖**：为核心功能编写单元测试
- **文档完善**：为公共包提供完整的文档

### 示例结构
```
aladdin-app/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── module/
│   │   ├── nn/
│   │   │   ├── nn_feeds_module_builder.go
│   │   │   └── nn_red_packet_module_builder.go
│   │   ├── seckill/
│   │   │   └── seckill_module_builder.go
│   │   └── common/
│   │       └── base_module_builder.go
│   ├── domain/
│   │   ├── module/
│   │   │   ├── nn_feeds_vo.go
│   │   │   └── nn_red_packet_vo.go
│   │   ├── nn/
│   │   │   ├── nn_round_feature_bo.go
│   │   │   └── fund_query_dto.go
│   │   └── seckill/
│   │       └── seckill_activity_bo.go
│   ├── dataservice/
│   │   ├── data_service.go
│   │   └── impl/
│   │       ├── nn_red_packet_data_service.go
│   │       └── seckill_data_service.go
│   └── provider/
│       ├── fund_provider.go
│       └── user_provider.go
├── pkg/
│   ├── utils/
│   │   └── context_util.go
│   ├── config/
│   │   └── config.go
│   └── logger/
│       └── logger.go
├── configs/
│   ├── config.yaml
│   └── config.prod.yaml
├── scripts/
│   ├── build.sh
│   └── deploy.sh
├── go.mod
├── go.sum
└── README.md
```
