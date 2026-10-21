# 使用指南

本文介绍玩家查询选项、比赛记录与排行榜参数、地图轮换、响应解码、错误处理和客户端配置。安装与快速开始见 [README](../README.md)，各接口字段说明见 README 的 [接口覆盖](../README.md#接口覆盖)。

下列代码片段沿用快速开始中创建的 `client` 和 `ctx`；按片段需要导入对应标准库。

## 玩家查询选项

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

## 比赛记录、排行榜及地图轮换

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

`Decode` 也支持调用方按实际数据定义的结构体；解码到 `any` 或 `map[string]any` 时，数字使用 `json.Number`，避免转换为 `float64` 造成精度丢失。空正文、格式错误或多个 JSON 值会返回解码错误。SDK 会识别非空 JSON `Error` 字符串，即使 HTTP 状态为 2xx 也返回错误。调用成功仍不保证正文是有效 JSON，完整 JSON 格式检查在 `Decode` 时进行。

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

此片段还需要导入标准库 `errors`。全部非 2xx 响应均返回 `*apex.APIError`，包括官网未列出的状态码；2xx 响应中的非空 JSON `Error` 字符串也会返回该错误。其中包含实际的 `StatusCode`、`Header` 和原始 `Body`。收到错误 HTTP 响应时，返回的 `Response` 也可读取；参数校验及建立连接失败时通常为 `nil`。正文读取失败可能同时包含 HTTP 错误和读取错误，仍可使用 `errors.As`／`errors.Is` 判断。

已识别的上游权限错误可通过 `errors.Is(err, apex.ErrAccessDenied)` 判断；玩家查找失败可通过 `errors.Is(err, apex.ErrPlayerNotFound)` 判断。玩家查找失败也可能来自上游查找服务异常，不能仅凭该错误认定账号不存在。

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

