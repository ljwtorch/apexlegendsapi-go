# apexlegendsapi-go

使用 Go 调用 [Apex Legends Status API](https://apexlegendsapi.com/) 的客户端包，覆盖玩家统计、比赛记录、排行榜、地图轮换、猎杀者门槛、服务器状态和名称转 UID。

该服务是第三方**非官方 API**，由 Apex Legends Status 提供，与 EA／Respawn 官方服务无隶属关系。接口的可用性、开放权限和数据内容由上游决定。

## 安装

要求 Go **1.27.2 或更高版本**；客户端仅使用标准库。

```bash
go get github.com/ljwtorch/apexlegendsapi-go
```

导入时包名为 `apex`：

```go
import apex "github.com/ljwtorch/apexlegendsapi-go"
```

## 快速开始

在[开发者门户](https://apexlegendsapi.com/auth)申请 API Key，设置环境变量：

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
		Player:   "ORIGIN_ACCOUNT_NAME",
		Platform: apex.PlatformPC,
		Options:  apex.PlayerOptions{Version: 5},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(response.Body))
}
```

以上代码也位于 [examples/basic/main.go](examples/basic/main.go)。在项目根目录执行 `go run ./examples/basic` 可运行；先将示例玩家名替换为真实 Origin 账号名。Steam 玩家也使用其绑定的 Origin 名称。重复查询同一玩家时建议使用 UID。

示例打印的是上游原始 JSON。字段含义和脱敏响应示例见 [PlayerByName 响应文档](docs/player-by-name.md)；其他接口的响应说明见下方接口表中的文档链接。

## 接口覆盖

所有方法首个参数都是 `context.Context`，返回 `(*apex.Response, error)`。

| 方法 | 请求参数 | 上游接口 |
|---|---|---|
| [`PlayerByName`](docs/player-by-name.md) | `PlayerByNameRequest` | `/bridge`，按名称查询统计 |
| [`PlayerByUID`](docs/player-by-uid.md) | `PlayerByUIDRequest` | `/bridge`，按 UID 查询统计 |
| [`Games`](docs/games.md) | `GamesRequest` | `/games`，查询已收集的比赛记录 |
| [`LegacyHistory`](docs/legacy-history.md) | `LegacyHistoryRequest` | `/bridge?history=1`，管理旧版比赛追踪 |
| [`Leaderboard`](docs/leaderboard.md) | `LeaderboardRequest` | `/leaderboard`，排行榜或追踪器 key 列表 |
| [`MapRotation`](docs/map-rotation.md) | `*MapRotationOptions`，可传 `nil` | `/maprotation`，地图轮换 |
| [`Predator`](docs/predator.md) | 无 | `/predator`，猎杀者门槛和大师人数 |
| [`Servers`](docs/servers.md) | 无 | `/servers`，服务器状态 |
| [`OriginUID`](docs/origin.md) | `player string` | `/origin`，Origin 名称查 UID |
| [`NameToUID`](docs/name-to-uid.md) | `NameToUIDRequest` | `/nametouid`，按平台和名称查 UID |

完整参数、错误码及平台支持见 [APEX_LEGENDS_API_ZH.md](APEX_LEGENDS_API_ZH.md)。每个导出类型和方法均有 Go 文档注释；[examples](examples) 目录包含每个接口的可直接运行示例。

接口响应文档说明返回字段和解码方式。

### 可运行示例

先设置 `APEX_API_KEY`，再执行对应目录：

| 接口 | 命令 | 示例源码 |
|---|---|---|
| `PlayerByName` | `go run ./examples/player-by-name` | [main.go](examples/player-by-name/main.go) |
| `PlayerByUID` | `go run ./examples/player-by-uid` | [main.go](examples/player-by-uid/main.go) |
| `Games` | `go run ./examples/games` | [main.go](examples/games/main.go) |
| `LegacyHistory` | `go run ./examples/legacy-history` | [main.go](examples/legacy-history/main.go) |
| `Leaderboard` | `go run ./examples/leaderboard` | [main.go](examples/leaderboard/main.go) |
| `MapRotation` | `go run ./examples/map-rotation` | [main.go](examples/map-rotation/main.go) |
| `Predator` | `go run ./examples/predator` | [main.go](examples/predator/main.go) |
| `Servers` | `go run ./examples/servers` | [main.go](examples/servers/main.go) |
| `OriginUID` | `go run ./examples/origin` | [main.go](examples/origin/main.go) |
| `NameToUID` | `go run ./examples/name-to-uid` | [main.go](examples/name-to-uid/main.go) |

## 使用指南

详细用法见 [使用指南](docs/usage.md)：

- [玩家查询选项](docs/usage.md#玩家查询选项)
- [比赛记录、排行榜及地图轮换](docs/usage.md#比赛记录排行榜及地图轮换)
- [响应与 JSON 解码](docs/usage.md#响应与-json-解码)
- [错误处理](docs/usage.md#错误处理)
- [客户端配置](docs/usage.md#客户端配置)

## 上游限制与接入要求

- 默认所有接口合计每秒 5 次请求。SDK 不内置限流或自动重试，由应用按 Key 统一协调；`X-Current-Rate` 响应头可用于观察当前速率。发生 429 时由应用决定等待策略。
- 新版比赛记录需要白名单，官网称目前不向新用户开放，每小时最多查询 5 名不同玩家。应用需要每 4 分钟请求一次目标玩家的 `/bridge` 来收集记录，再使用 `/games` 查询；SDK 不启动后台任务。
- 旧版比赛记录及排行榜目前也不向新用户开放。SDK 支持这些接口不代表能获得访问权限。
- Webhook 目前仅支持比赛记录，需要在开发者门户配置。接收端接受 POST 请求，User-Agent 为 `ApexAPI Webhook/0.1`，须在 3 秒内接受请求。SDK 不提供 HTTP 接收服务器。
- 展示排行榜数据时，除白标授权外，须添加指向 [Apex Legends Status](https://apexlegendsstatus.com/) 的可点击链接，文案为 **Data provided by Apex Legends Status**。展示服务器状态时须提供该站链接或注明 **Data from apexlegendsstatus.com**。

## 开发与验证

```bash
go test -race ./...
go vet ./...
```

测试通过本地模拟服务器和模拟 HTTP Transport 验证全部接口，无需 API Key，也不会请求真实服务。没有 `Output` 注释的 Go 示例只参与编译，不在测试中调用远程 API。

## 许可证

[Apache-2.0](LICENSE)
