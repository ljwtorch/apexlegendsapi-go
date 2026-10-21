# PlayerByUID：按 UID 查询玩家统计

- **SDK 方法：** Client.PlayerByUID
- **HTTP 接口：** GET /bridge
- **官网文档：** <https://apexlegendsapi.com/#query-by-uid>
- **字段核对：** 2026-10-21
- **可运行示例：** [examples/player-by-uid/main.go](../examples/player-by-uid/main.go)

## 请求

必填参数是 uid 和 platform。支持 PC、PS4、X1 和 SWITCH；重复查询同一玩家时推荐使用 UID，可避免改名造成查询失败。其余可选参数与 PlayerByName 相同，详见 [按名称查询文档](player-by-name.md)。

## 返回值

返回结构与 [PlayerByName](player-by-name.md) 相同：顶层通常包含 global、realtime、legends、total 和 processingTime。详细字段含义、可选值和脱敏 JSON 请参阅 [PlayerByName 响应字段表](player-by-name.md#返回值)。

| 路径 | 类型 | 含义 | 来源 |
|---|---|---|---|
| global | object | 玩家账号、段位、战斗通行证和徽章 | 样例观察 |
| realtime | object | 在线和当前游戏状态 | 样例观察 |
| legends | object | 当前传奇及全部传奇追踪器 | 样例观察 |
| total | object | tracker key 汇总值 | 样例观察 |
| processingTime | number | 上游处理耗时，通常以秒为单位 | 样例观察 |

SWITCH 查询只改变请求平台，不改变响应字段约定；实际字段仍可能因玩家数据和上游版本而缺失或为 null。

## 脱敏示例

~~~json
{
  "global": {
    "name": "PLAYER_NAME",
    "uid": "PLAYER_UID",
    "platform": "SWITCH",
    "rank": {
      "rankName": "Unranked",
      "rankScore": 0
    }
  },
  "realtime": {},
  "legends": {},
  "total": {},
  "processingTime": 0
}
~~~

## Go 解码

~~~go
response, err := client.PlayerByUID(ctx, apex.PlayerByUIDRequest{
    UID:      "PLAYER_UID",
    Platform: apex.PlatformPC,
})
if err != nil {
    return err
}

var data any
if err := response.Decode(&data); err != nil {
    return err
}
~~~

## 限制与错误

官网错误码、鉴权、速率限制和 skipRank、merge 等参数说明见 [总览文档](../APEX_LEGENDS_API_ZH.md)。不要把样例字段当作长期稳定的 Go 结构体契约。
