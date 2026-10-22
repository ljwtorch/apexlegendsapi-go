# apexlegendsapi-go

轻量、零第三方依赖的 [Apex Legends Status API](https://apexlegendsapi.com/) Go 客户端。

[![Go Reference](https://pkg.go.dev/badge/github.com/ljwtorch/apexlegendsapi-go.svg)](https://pkg.go.dev/github.com/ljwtorch/apexlegendsapi-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/ljwtorch/apexlegendsapi-go)](https://goreportcard.com/report/github.com/ljwtorch/apexlegendsapi-go)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache--2.0-blue.svg)](LICENSE)

[功能](#功能) · [快速开始](#快速开始) · [接口覆盖](#接口覆盖) · [使用指南](#使用指南) · [接入须知](#接入须知)

查询玩家统计、比赛记录、排行榜、地图轮换、猎杀者门槛和服务器状态，也可通过名称获取玩家 UID，适合构建数据看板、机器人和游戏工具。

> [!NOTE]
> 上游是由 Apex Legends Status 提供的第三方非官方 API，与 EA／Respawn 官方服务无隶属关系。接口可用性、访问权限和数据内容由上游决定。

## 功能

- **完整接口封装**：提供 10 个方法，使用类型化请求参数和平台常量，并在发送前校验参数。
- **保留原始数据**：返回 HTTP 状态码、响应头和原始正文，可按需解码为自定义结构体，不丢弃未知字段。
- **可控的请求生命周期**：所有调用支持 `context.Context`；客户端可跨 goroutine 复用，默认超时 30 秒。
- **可识别的错误**：通过 `APIError`、`errors.Is` 和 `errors.As` 处理上游错误、权限不足、玩家查找失败及超时。
- **灵活配置**：支持自定义 `http.Client`、API 地址，以及 Header 或查询参数鉴权；仅使用 Go 标准库。

## 快速开始

### 安装

需要 **Go 1.27.2 或更高版本**，以及在[开发者门户](https://apexlegendsapi.com/auth)申请的 API Key。

在自己的 Go 项目中安装；新项目请先执行 `go mod init example.com/apex-demo`：

```bash
go get github.com/ljwtorch/apexlegendsapi-go
```

### 查询玩家统计

设置环境变量，将以下代码保存为 `main.go`，并将 `YOUR_ORIGIN_NAME` 替换为真实账号名：

```bash
export APEX_API_KEY='YOUR_API_KEY'
```

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	apex "github.com/ljwtorch/apexlegendsapi-go"
)

func main() {
	client, err := apex.NewClient(os.Getenv("APEX_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.PlayerByName(ctx, apex.PlayerByNameRequest{
		Player:   "YOUR_ORIGIN_NAME",
		Platform: apex.PlatformPC,
		Options:  apex.PlayerOptions{Version: 5},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(response.Body))
}
```

运行 `go run .`，即可打印上游返回的原始 JSON。导入路径为 `github.com/ljwtorch/apexlegendsapi-go`，包名为 `apex`。

> [!TIP]
> PC 玩家使用 Origin 账号名；Steam 玩家也需要使用其绑定的 Origin 名称。重复查询同一玩家时，建议改用 `PlayerByUID`，避免改名影响查询。按名称查询支持 PC、PlayStation 和 Xbox；按 UID 查询还支持 Switch。

## 接口覆盖

所有方法首个参数均为 `context.Context`，返回 `(*apex.Response, error)`。点击方法名查看请求参数、响应字段和解码示例。

| 方法 | 功能 | 上游接口 | 示例源码 |
| --- | --- | --- | --- |
| [`PlayerByName`](docs/player-by-name.md) | 按名称查询玩家统计 | `/bridge` | [player-by-name](examples/player-by-name/main.go) |
| [`PlayerByUID`](docs/player-by-uid.md) | 按 UID 查询玩家统计 | `/bridge` | [player-by-uid](examples/player-by-uid/main.go) |
| [`Games`](docs/games.md) | 查询已收集的比赛记录 | `/games` | [games](examples/games/main.go) |
| [`LegacyHistory`](docs/legacy-history.md) | 查询、添加或删除旧版比赛追踪 | `/bridge?history=1` | [legacy-history](examples/legacy-history/main.go) |
| [`Leaderboard`](docs/leaderboard.md) | 查询排行榜或可用追踪器 key | `/leaderboard` | [leaderboard](examples/leaderboard/main.go) |
| [`MapRotation`](docs/map-rotation.md) | 查询当前及下一张地图 | `/maprotation` | [map-rotation](examples/map-rotation/main.go) |
| [`Predator`](docs/predator.md) | 查询猎杀者门槛及大师人数 | `/predator` | [predator](examples/predator/main.go) |
| [`Servers`](docs/servers.md) | 查询服务器状态 | `/servers` | [servers](examples/servers/main.go) |
| [`OriginUID`](docs/origin.md) | 通过 Origin 名称获取 UID | `/origin` | [origin](examples/origin/main.go) |
| [`NameToUID`](docs/name-to-uid.md) | 按平台和名称获取 UID | `/nametouid` | [name-to-uid](examples/name-to-uid/main.go) |

克隆仓库后，设置 `APEX_API_KEY`，按需替换示例源码中的玩家名、UID 和平台，在项目根目录运行：

```bash
go run ./examples/basic
go run ./examples/map-rotation
```

其他示例使用相同方式运行，例如 `go run ./examples/servers`。`LegacyHistory` 的 `HistoryAdd` 和 `HistoryDelete` 操作会修改上游追踪列表。

## 使用指南

### 读取响应

`Response` 包含 `StatusCode`、`Header` 和 `Body json.RawMessage`。成功调用后，可通过 `Decode` 解码为自定义结构体，或先使用 `any` 探索数据：

```go
var data any
if err := response.Decode(&data); err != nil {
    log.Fatal(err)
}
fmt.Println(data)
```

解码到 `any` 或 `map[string]any` 时，数字使用 `json.Number`，避免大整数精度丢失。完整 JSON 格式校验发生在 `Decode` 时；HTTP 调用成功不代表正文一定是有效 JSON。

### 配置与错误处理

默认使用 `Authorization: YOUR_API_KEY` 鉴权，不添加 `Bearer` 前缀。通过 `NewClient` 的选项可调整行为：

| 选项 | 用途 |
| --- | --- |
| `WithHTTPClient` | 自定义超时、Transport 和重定向策略 |
| `WithBaseURL` | 设置 API 地址，支持代理路径前缀和本地测试地址 |
| `WithAuthMode` | 使用 `AuthHeader`（默认）或 `AuthQuery` 鉴权 |

非 2xx 响应及包含非空 JSON `Error` 字符串的响应会返回 `*apex.APIError`，其中保留状态码、响应头和原始正文。使用 `errors.As` 读取详情，使用 `errors.Is` 判断 `ErrAccessDenied`、`ErrPlayerNotFound` 或 context 取消与超时错误。

更多用法见[使用指南](docs/usage.md)：[玩家查询选项](docs/usage.md#玩家查询选项)、[比赛记录与地图轮换](docs/usage.md#比赛记录排行榜及地图轮换)、[响应解码](docs/usage.md#响应与-json-解码)、[错误处理](docs/usage.md#错误处理)和[客户端配置](docs/usage.md#客户端配置)。

## 接入须知

> [!IMPORTANT]
> 比赛记录和排行榜有额外访问限制。上游目前不向新用户开放新版／旧版比赛记录和排行榜 API；新版比赛记录还需要白名单。SDK 提供接口封装不代表 API Key 已具备访问权限。

根据[上游文档](https://apexlegendsapi.com/)，接入时还需注意：

- **请求限流**：新建 Key 默认所有接口合计每秒 5 次请求，可通过响应头 `X-Current-Rate` 观察当前速率。SDK 不内置限流或自动重试，需由应用统一控制。
- **比赛数据收集**：新版比赛记录每小时最多查询 5 名不同玩家；需每 4 分钟请求一次目标玩家的 `/bridge` 收集数据，再通过 `/games` 查询。应用负责调度及 Webhook 接收。
- **地图轮换版本**：使用 `MapRotationOptions{Version: 2}` 获取所有模式；传 `nil` 或版本 0 会使用上游默认版本 1，仅返回大逃杀普通模式。
- **数据署名**：展示排行榜时，除白标授权外，须提供 [Data provided by Apex Legends Status](https://apexlegendsstatus.com/) 链接；展示服务器状态时，须提供该站链接或注明 `Data from apexlegendsstatus.com`。

接口权限、参数和通用错误码也可查阅 [API 中文参考](APEX_LEGENDS_API_ZH.md)，最新要求以上游文档为准。

## 本地开发

```bash
git clone https://github.com/ljwtorch/apexlegendsapi-go.git
cd apexlegendsapi-go
go test -race ./...
go vet ./...
```

源码位于根目录，接口使用文档位于 [docs](docs)，可运行示例位于 [examples](examples)。测试和静态检查无需 API Key；运行示例会请求真实上游服务。
