# ABTest 基础库 / ABTest Base Library

一个基于 Go 语言开发的 AB 测试基础库，提供简单易用的 AB 测试分组功能。  
A Go-based AB testing base library that provides simple and easy-to-use AB testing grouping functionality.

## 原理 / Principle

AB 测试（也称为拆分测试）是一种比较两个或多个版本的方法，以确定哪个版本性能更好。本系统基于以下原理实现：  
AB testing (also known as split testing) is a method of comparing two or more versions to determine which performs better. This system is implemented based on the following principles:

1. **哈希分组**：使用设备ID和实验ID进行MD5哈希计算，确保同一设备在不同实验中分组一致  
   **Hash Grouping**: Use MD5 hash calculation with device ID and experiment ID to ensure consistent grouping of the same device across different experiments

2. **白名单机制**：支持为特定设备指定分组，用于内部测试和验证  
   **Whitelist Mechanism**: Support specifying groups for specific devices for internal testing and verification

3. **互斥分层**：确保同一实验的不同分层不会同时作用于同一用户  
   **Mutual Exclusion Layering**: Ensure that different layers of the same experiment do not act on the same user simultaneously

## 目的 / Purpose

- 为产品提供科学的数据驱动决策支持  
  Provide scientific data-driven decision support for products

- 降低新功能上线风险  
  Reduce the risk of launching new features

- 优化用户体验和产品性能  
  Optimize user experience and product performance

- 提供灵活可配置的AB测试框架  
  Provide a flexible and configurable AB testing framework

## 功能特性 / Features

- ✅ 基于设备ID的AB测试分组 / Device ID-based AB test grouping
- ✅ 白名单机制支持 / Whitelist mechanism support
- ✅ 多实验并行支持 / Multi-experiment parallel support
- ✅ 分层互斥确保测试独立性 / Layered mutual exclusion ensures test independence
- ✅ HTTP API 接口 / HTTP API interface
- ✅ 实时配置更新 / Real-time configuration updates

## 使用方法 / Usage

### 1. 启动服务 / Starting the Service

```bash
go run main.go


（1）查询AB测试分组 / Querying AB Test Groups
通过HTTP接口查询设备的分组情况：
Query device grouping through the HTTP interface:

bash
# 通过查询参数传递设备ID / Pass device ID via query parameters
curl "http://localhost:8080/abtest?device_id=test_device_123"

# 通过请求头传递设备ID / Pass device ID via request header
curl -H "device_id: test_device_123" http://localhost:8080/abtest


（2）响应格式 / Response Format
json
{
  "data": {
    "ab_id": "1-1,2-2"
  }
}
其中 ab_id 字段格式为 实验ID-分层ID，多个实验用逗号分隔。
The ab_id field format is experimentID-layerID, with multiple experiments separated by commas.


（3）配置AB测试 / Configuring AB Tests
在 main.go 的 initSampleConfig 函数中配置AB测试：
Configure AB tests in the initSampleConfig function in main.go:

go
configs := []*model.AbTestConfig{
    {
        AbId: 1, // 实验ID / Experiment ID
        AbConfig: `[
            {"layer_id": 1, "slot_ids": [0, 1, 2, 3, 4], "white_list": ["test_device_1"]},
            {"layer_id": 2, "slot_ids": [5, 6, 7, 8, 9], "white_list": []}
        ]`,
    },
}
配置字段说明：
Configuration field description:

layer_id: 分层ID / Layer ID

slot_ids: 该分层占用的槽位（0-999） / Slots occupied by this layer (0-999)

white_list: 白名单设备ID列表 / Whitelist device ID list


（4）添加新的AB测试 / Adding New AB Tests
在配置中添加新的实验 / Add a new experiment to the configuration

调用 model.ProcessABConfig(configs) 处理配置 / Call model.ProcessABConfig(configs) to process the configuration

服务会自动应用新配置 / The service will automatically apply the new configuration


（5）项目结构 / Project Structure
text
abtest/
├── model/         # 数据模型和核心逻辑 / Data models and core logic
│   ├── model.go   # ABTest结构体和配置处理 / ABTest structure and configuration processing
│   └── core.go    # 分组逻辑核心实现 / Core implementation of grouping logic
├── util/          # 工具函数 / Utility functions
│   └── util.go    # 通用工具函数 / General utility functions
├── main.go        # 程序入口和HTTP服务 / Program entry and HTTP service
└── go.mod         # 模块依赖管理 / Module dependency management