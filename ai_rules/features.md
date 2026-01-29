# 功能实现规范

## Go 实现规范

### 数据服务层
- 必须实现 DataService 接口
- 使用构造函数注入依赖
- Execute 方法的第一个参数是 InvocationContext
- Execute 方法的第二个参数是 map[string]interface{} 或自定义结构体
- 必须返回 (T, error) 类型

示例：
```
gopackage dataservice

import (
    "context"
)

// DataService 数据服务接口
type DataService[T any] interface {
    Execute(ctx context.Context, businessReq map[string]interface{}) (T, error)
}

// NnRedPacketDataService 红包数据服务
type NnRedPacketDataService struct {
    // 依赖注入的字段
}

// NewNnRedPacketDataService 创建红包数据服务实例
func NewNnRedPacketDataService() *NnRedPacketDataService {
    return &NnRedPacketDataService{}
}

// Execute 执行数据服务逻辑
func (s *NnRedPacketDataService) Execute(ctx context.Context, businessReq map[string]interface{}) ([]FundQueryDTO, error) {
    // 实现逻辑
    return nil, nil
}
```

### 模块构建器
- 必须实现 ModuleBuilder 接口
- 使用构造函数注入依赖
- 实现 Name()、Build()、Transform() 三个方法
- 通过 Context.GetBizResult() 获取数据服务结果

示例：
```
gopackage builder

import (
    "context"
)

// ModuleBuilder 模块构建器接口
type ModuleBuilder[T any] interface {
    Name() string
    Build(ctx context.Context) (T, error)
    Transform(ctx context.Context, data T) (T, error)
}

// NnRedPacketModuleBuilder 红包模块构建器
type NnRedPacketModuleBuilder struct {
    // 依赖注入的字段
    redPacketService *dataservice.NnRedPacketDataService
}

// NewNnRedPacketModuleBuilder 创建红包模块构建器实例
func NewNnRedPacketModuleBuilder(redPacketService *dataservice.NnRedPacketDataService) *NnRedPacketModuleBuilder {
    return &NnRedPacketModuleBuilder{
        redPacketService: redPacketService,
    }
}

// Name 返回模块名称
func (b *NnRedPacketModuleBuilder) Name() string {
    return "nnRedPacket"
}

// Build 构建模块数据
func (b *NnRedPacketModuleBuilder) Build(ctx context.Context) (NnRedPacketVO, error) {
    funds, err := b.redPacketService.Execute(ctx, nil)
    if err != nil {
        return NnRedPacketVO{}, err
    }
    // 构建逻辑
    return NnRedPacketVO{}, nil
}

// Transform 转换模块数据
func (b *NnRedPacketModuleBuilder) Transform(ctx context.Context, data NnRedPacketVO) (NnRedPacketVO, error) {
    // 转换逻辑
    return data, nil
}
```

### 依赖注入
- 使用构造函数进行依赖注入
- 避免使用全局变量
- 推荐使用接口而非具体实现

示例：
```
gopackage main

import (
    "context"
    "yourproject/builder"
    "yourproject/dataservice"
)

func main() {
    // 创建依赖实例
    redPacketService := dataservice.NewNnRedPacketDataService()
    
    // 注入依赖
    redPacketBuilder := builder.NewNnRedPacketModuleBuilder(redPacketService)
    
    // 使用构建器
    ctx := context.Background()
    result, err := redPacketBuilder.Build(ctx)
    if err != nil {
        // 错误处理
        return
    }
    
    // 后续处理
}
```

### 错误处理
- 必须在函数签名中返回 error
- 使用标准的错误处理模式：if err != nil
- 错误信息应清晰明了，包含上下文信息
- 推荐使用 errors.Wrap 或 fmt.Errorf 添加错误上下文

示例：
```
goimport (
    "errors"
    "fmt"
)

func (s *NnRedPacketDataService) Execute(ctx context.Context, businessReq map[string]interface{}) ([]FundQueryDTO, error) {
    if businessReq == nil {
        return nil, errors.New("businessReq cannot be nil")
    }
    
    // 业务逻辑
    funds, err := s.fetchFunds(ctx)
    if err != nil {
        return nil, fmt.Errorf("fetch funds failed: %w", err)
    }
    
    return funds, nil
}
```

### 上下文传递
- 使用 context.Context 传递上下文信息
- 避免在函数参数中传递过多的上下文信息
- 使用 context.WithValue 传递必要的上下文数据

示例：```gofunc (s *NnRedPacketDataService) Execute(ctx context.Context, businessReq map[string]interface{}) ([]FundQueryDTO, error) {
    // 从上下文获取信息
    userId, ok := ctx.Value("userId").(string)
    if !ok {
        return nil, errors.New("userId not found in context")
    }
    
    // 使用 userId 进行业务逻辑
    // ...
    
    return funds, nil
}
```

### 包结构规范
- 按功能划分包，如 dataservice、builder、model 等
- 包名使用小写字母，简短且有意义
- 避免包之间的循环依赖
- 保持包的职责单一

### 接口设计规范
- 接口应小而专注，通常只包含1-3个方法
- 接口命名使用名词或动名词，如 DataService、ModuleBuilder
- 方法命名使用动词或动词短语，如 Execute、Build
- 接口应定义在使用它的包中，而非实现它的包中