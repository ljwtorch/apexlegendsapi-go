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

| 方法 | 请求参数 | 上游接口 | 响应文档 |
|---|---|---|---|
| [`PlayerByName`](docs/player-by-name.md) | `PlayerByNameRequest` | `/bridge`，按名称查询统计 | [返回值](docs/player-by-name.md) |
| [`PlayerByUID`](docs/player-by-uid.md) | `PlayerByUIDRequest` | `/bridge`，按 UID 查询统计 | [返回值](docs/player-by-uid.md) |
| [`Games`](docs/games.md) | `GamesRequest` | `/games`，查询已收集的比赛记录 | [返回值](docs/games.md) |
| [`LegacyHistory`](docs/legacy-history.md) | `LegacyHistoryRequest` | `/bridge?history=1`，管理旧版比赛追踪 | [返回值](docs/legacy-history.md) |
| [`Leaderboard`](docs/leaderboard.md) | `LeaderboardRequest` | `/leaderboard`，排行榜或追踪器 key 列表 | [返回值](docs/leaderboard.md) |
| [`MapRotation`](docs/map-rotation.md) | `*MapRotationOptions`，可传 `nil` | `/maprotation`，地图轮换 | [返回值](docs/map-rotation.md) |
| [`Predator`](docs/predator.md) | 无 | `/predator`，猎杀者门槛和大师人数 | [返回值](docs/predator.md) |
| [`Servers`](docs/servers.md) | 无 | `/servers`，服务器状态 | [返回值](docs/servers.md) |
| [`OriginUID`](docs/origin.md) | `player string` | `/origin`，Origin 名称查 UID | [返回值](docs/origin.md) |
| [`NameToUID`](docs/name-to-uid.md) | `NameToUIDRequest` | `/nametouid`，按平台和名称查 UID | [返回值](docs/name-to-uid.md) |

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

### 玩家查询选项

```go
clubs := false
response, err := client.PlayerByUID(ctx, apex.PlayerByUIDRequest{
	UID:      "PLAYER_UID",
	Platform: apex.PlatformSwitch,
	Options: apex.PlayerOptions{
		Version:         5,
		EnableClubsBeta: &clubs,
		SkipRank:        true,
		Merge:           true,
		RemoveMerged:    true,
	},
})
```

`Version=0` 省略版本参数，使用上游默认值 5。`EnableClubsBeta=nil` 表示不传，指针可显式传入 `true` 或 `false`。`SkipRank`、`Merge`、`RemoveMerged` 为存在即生效的参数，只有设为 `true` 才会发送。SDK 将 UID 作为字符串传递。

按名称查询及 `NameToUID` 支持 `PlatformPC`、`PlatformPS4`、`PlatformX1`；按 UID 查询及旧版记录还支持 `PlatformSwitch`；排行榜另外支持 `PlatformAny`。移动平台不受支持。

### 比赛记录、排行榜及地图轮换

```go
start := int64(1700000000)
limit := 10
response, err := client.Games(ctx, apex.GamesRequest{
	UID:   "PLAYER_UID",
	Mode:  apex.ModeBattleRoyale,
	Start: &start,
	Limit: &limit,
})
```

`GamesRequest.Start`、`End` 和 `Limit` 使用指针，`nil` 省略参数，非空指针按实际值发送（包括零）。模式可选 `ModeBattleRoyale`、`ModeArenas`、`ModeUnknown`，空字符串表示不筛选。

旧版记录通过 `LegacyHistoryRequest{UID: "PLAYER_UID", Platform: apex.PlatformPC, Action: apex.HistoryInfo}` 调用；操作常量还包括 `HistoryGet`、`HistoryAdd`、`HistoryDelete`。按官网参数表，所有操作均要求 UID 和平台；新增和删除会修改上游追踪列表。

排行榜使用 `LeaderboardRequest{Platform: apex.PlatformAny, Legend: "Global", Key: "rankScore"}`。省略 `Key` 获取指定传奇的可用追踪器 key；传奇名称按官网要求首字母大写。

完整地图轮换使用 `client.MapRotation(ctx, &apex.MapRotationOptions{Version: 2})`。传 `nil` 或版本 0 会省略版本参数，遵循上游默认的版本 1，仅返回大逃杀普通模式。

## 响应与 JSON 解码

`Response` 提供 `StatusCode`、`Header` 和 `Body json.RawMessage`。SDK 完整保留响应正文，不假定它一定是对象，也不预先定义官网未公开的响应字段。

```go
// 接续一次成功的接口调用。
fmt.Println(response.StatusCode)
fmt.Println(response.Header.Get("X-Current-Rate"))
fmt.Println(string(response.Body))

var data any
if err := response.Decode(&data); err != nil {
	log.Fatal(err)
}
```

`Decode` 也支持调用方按实际数据定义的结构体；解码到 `any` 或 `map[string]any` 时，数字使用 `json.Number`，避免转换为 `float64` 造成精度丢失。空正文、格式错误或多个 JSON 值会返回解码错误。HTTP 调用成功只表示收到 2xx 响应；JSON 格式检查在 `Decode` 时进行。

## 错误处理

```go
response, err := client.Servers(ctx)
if err != nil {
	var apiErr *apex.APIError
	switch {
	case errors.As(err, &apiErr):
		fmt.Printf("HTTP %d, 当前速率 %s\n", apiErr.StatusCode, apiErr.Header.Get("X-Current-Rate"))
	case errors.Is(err, context.Canceled):
		fmt.Println("请求已取消")
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Println("请求超时")
	default:
		fmt.Println(err)
	}
	return
}
fmt.Println(string(response.Body))
```

此片段还需要导入标准库 `errors`。全部非 2xx 响应均返回 `*apex.APIError`，包括官网未列出的状态码；其中包含 `StatusCode`、`Header` 和原始 `Body`。收到错误 HTTP 响应时，返回的 `Response` 也可读取；参数校验及建立连接失败时通常为 `nil`。正文读取失败可能同时包含 HTTP 错误和读取错误，仍可使用 `errors.As`／`errors.Is` 判断。

错误字符串省略响应正文和请求 URL；需要查看正文时请显式读取 `apiErr.Body`。原始正文和响应头可能包含服务方回显的数据，记录日志前应自行处理。网络错误、取消和超时保留原始错误链。

## 客户端配置

```go
client, err := apex.NewClient(os.Getenv("APEX_API_KEY"),
	apex.WithHTTPClient(&http.Client{Timeout: 15 * time.Second}),
	apex.WithBaseURL("https://api.apexlegendsstatus.com"),
	apex.WithAuthMode(apex.AuthHeader),
)
```

默认请求超时为 30 秒。提供自定义 `http.Client` 后，采用它自己的超时和重定向策略。`WithBaseURL` 支持 HTTP(S) 地址和路径前缀，适用于代理和本地测试，不接受凭据、查询字符串或片段。

默认发送 `Authorization: YOUR_API_KEY`，不添加 `Bearer` 前缀。`AuthQuery` 改用 `auth` 查询参数。客户端可以并发复用；不要在使用期间修改传入的 HTTP 客户端配置。

默认发送 `Accept: */*`，与官网 curl 示例一致。上游无扩展名路由会先对 PHP 脚本进行内容协商；指定 `Accept: application/json` 可能在执行 API 前返回 `406 Not Acceptable`。通配请求头兼容该路由，响应仍可通过 `Response.Decode` 解码为 JSON。

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
