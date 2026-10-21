# LegacyHistory：旧版比赛记录操作

- **SDK 方法：** Client.LegacyHistory
- **HTTP 接口：** GET /bridge?history=1
- **官网文档：** <https://apexlegendsapi.com/#legacy-match-history>
- **字段核对：** 2026-10-21
- **可运行示例：** [examples/legacy-history/main.go](../examples/legacy-history/main.go)

## 请求

必填参数是 uid、platform 和 action，并由 SDK 固定发送 history=1。action 可选：

| action | 含义 | 来源 |
|---|---|---|
| info | 返回当前追踪列表中的 UID 和平台 | 官方文档 |
| get | 返回指定 UID 的全部已追踪比赛事件 | 官方文档 |
| add | 将 UID 加入追踪列表 | 官方文档 |
| delete | 从追踪列表移除 UID | 官方文档 |

旧版 API 目前不向新用户开放；只有先加入追踪列表，才可能通过 get 获取比赛数据。

## 返回值

返回内容随 action 变化，官网没有公开稳定 JSON Schema：

| action | 返回业务含义 | 字段确定性 |
|---|---|---|
| info | 当前 API Key 正在追踪的玩家清单 | 官方明确用途，具体字段未公开 |
| get | 指定玩家的比赛事件集合 | 官方明确用途，事件字段未公开 |
| add | 添加操作的结果或状态 | 上游未公开 |
| delete | 删除操作的结果或状态 | 上游未公开 |

脱敏根类型示例：

~~~json
{
  "action": "info",
  "players": []
}
~~~

上例仅用于说明操作结果对象的概念；action、players 等字段名不是官网承诺的固定 Schema。

## Go 解码

~~~go
response, err := client.LegacyHistory(ctx, apex.LegacyHistoryRequest{
    UID:      "PLAYER_UID",
    Platform: apex.PlatformPC,
    Action:   apex.HistoryInfo,
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

add 和 delete 会修改上游追踪列表。不要在没有确认业务意图时自动调用这两个操作。鉴权、速率限制和错误码见 [总览文档](../APEX_LEGENDS_API_ZH.md)。
