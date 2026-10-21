# Games：查询比赛记录

- **SDK 方法：** Client.Games
- **HTTP 接口：** GET /games
- **官网文档：** <https://apexlegendsapi.com/#new-match-history>
- **字段核对：** 2026-10-21
- **可运行示例：** [examples/games/main.go](../examples/games/main.go)

## 请求

uid 必填；mode 可选 BATTLE_ROYALE、ARENAS、UNKNOWN；start 和 end 是 Unix 时间戳；limit 限制最多返回的比赛数量。新版比赛记录需要白名单，目前不向新用户开放；每小时最多查询 5 名不同玩家，并需要定期请求 /bridge 收集数据。

`uid` 是玩家账号 UID，不能填写玩家名（例如 `torchlee`）。需要按名称查询时，先通过玩家查询接口取得 UID。

## 返回值

返回已收集比赛记录的 JSON 数组。官网说明每个元素代表一场比赛事件，但没有公开元素字段的稳定 Schema，因此 SDK 保留每个元素的原始 JSON。

| 返回内容 | 含义 | 来源 |
|---|---|---|
| 数组元素 | 一场已收集的比赛记录 | 官方文档 |
| 数组顺序 | 由上游返回决定；不要假设永久按时间排序 | 上游未公开 |
| 空数组 [] | 当前没有符合 UID 和筛选条件的记录 | 语义推断，需以实际响应为准 |
| 元素内部字段 | 模式、时间、地图、战绩等可能随上游版本变化 | 上游未公开 |

示例只表示根类型，不代表固定字段：

~~~json
[]
~~~

## Go 解码

推荐先解码到 []map[string]any，再根据实际响应定义业务结构：

~~~go
response, err := client.Games(ctx, apex.GamesRequest{
    UID:  "PLAYER_UID",
    Mode: apex.ModeBattleRoyale,
})
if err != nil {
    return err
}

var games []map[string]any
if err := response.Decode(&games); err != nil {
    return err
}
~~~

## 限制与错误

白名单、收集周期和速率限制见 [总览文档](../APEX_LEGENDS_API_ZH.md)。400、403、404、429 和 500 的含义以官网错误码表为准；未公开的字段或错误正文必须以实际响应为准。

实测 `/games` 也会在 HTTP 404 中返回 `{"Error":"This API key doesn't exist or you are not allowed to access this data."}`，因此不能仅凭状态码判断玩家不存在。SDK 将此响应识别为 `ErrAccessDenied`，可用 `errors.Is(err, apex.ErrAccessDenied)` 判断，并保留 `APIError` 中的原始状态码和响应体。需要检查 API key 和比赛记录白名单权限；正确 UID 也不能绕过上游权限限制。
