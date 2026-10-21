# Leaderboard：排行榜

- **SDK 方法：** Client.Leaderboard
- **HTTP 接口：** GET /leaderboard
- **官网文档：** <https://apexlegendsapi.com/#leaderboards>
- **字段核对：** 2026-10-21
- **可运行示例：** [examples/leaderboard/main.go](../examples/leaderboard/main.go)

## 请求

platform 和 legend 必填，key 可选。legend 首字母大写；使用 Global 可查询全局追踪器、等级、rankScore 和 arenaScore。省略 key 时请求可用的 tracker key 列表。

## 返回值

官网明确说明：指定 tracker 时返回数据库中前 500 名玩家的 UID、名称、平台和值；数据每 6 小时更新。省略 key 时返回该传奇可用的 key 列表。

| 返回内容 | 含义 | 来源 |
|---|---|---|
| 排行榜元素的 UID | 玩家稳定 UID | 官方文档 |
| 排行榜元素的名称 | 玩家名称 | 官方文档 |
| 排行榜元素的平台 | 玩家所属平台 | 官方文档 |
| 排行榜元素的值 | 指定 tracker 的统计值 | 官方文档 |
| 未传 key 的结果 | 指定传奇可用 tracker key 列表 | 官方文档 |
| 排名数量 | 通常最多 500 条 | 官方文档 |

官网没有发布固定的顶层字段名或数组包装方式，以下是脱敏的语义示意：

~~~json
[
  {
    "uid": "PLAYER_UID",
    "name": "PLAYER_NAME",
    "platform": "PC",
    "value": 0
  }
]
~~~

示意中的字段名来自官网对返回内容的描述，实际响应仍以服务端 JSON 为准。

## Go 解码

~~~go
response, err := client.Leaderboard(ctx, apex.LeaderboardRequest{
    Platform: apex.PlatformAny,
    Legend:   "Global",
    Key:      "rankScore",
})
if err != nil {
    return err
}

var entries []map[string]any
if err := response.Decode(&entries); err != nil {
    return err
}
~~~

## 限制与错误

该 API 目前不向新用户开放。展示数据时，除非拥有白标授权，必须链接 <https://apexlegendsstatus.com> 并显示 Data provided by Apex Legends Status。完整限制见 [总览文档](../APEX_LEGENDS_API_ZH.md)。
