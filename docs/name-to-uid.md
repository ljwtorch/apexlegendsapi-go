# NameToUID：按平台和名称查询 UID

- **SDK 方法：** Client.NameToUID
- **HTTP 接口：** GET /nametouid
- **官网文档：** <https://apexlegendsapi.com/#name-to-uid>
- **字段核对：** 2026-10-21
- **可运行示例：** [examples/name-to-uid/main.go](../examples/name-to-uid/main.go)

## 请求

必填参数是 player 和 platform。官网支持 PS4（PlayStation）、X1（Xbox）和 PC（Origin）。

## 返回值

返回指定平台玩家的 UID。该接口不返回完整统计。

| 返回内容 | 含义 | 来源 |
|---|---|---|
| UID | 名称和平台对应的稳定玩家标识 | 官方文档确认用途 |
| 平台信息 | 若上游返回，用于说明匹配的平台 | 上游未公开 |
| 未找到结果 | 名称、平台不匹配或玩家不存在 | 结合官网错误码说明 |

脱敏语义示意：

~~~json
{
  "uid": "PLAYER_UID"
}
~~~

官网没有承诺固定包装字段；上例只说明业务结果，不构成完整 Schema。

## Go 解码

~~~go
response, err := client.NameToUID(ctx, apex.NameToUIDRequest{
    Player:   "PLAYER_NAME",
    Platform: apex.PlatformPS4,
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

SWITCH 不属于该接口支持的平台，SDK 会在发起请求前返回本地校验错误。通用鉴权和错误码见 [总览文档](../APEX_LEGENDS_API_ZH.md)。

上游可能用 HTTP 200 返回 `{"Error":"Player not found. Please try again (err origin lookup)"}`。SDK 会将非空 JSON `Error` 字符串识别为业务错误，而不只检查 HTTP 状态码；上述查找失败可用 `errors.Is(err, apex.ErrPlayerNotFound)` 判断。`APIError.StatusCode` 仍保留实际的 200，响应正文也原样保留。此错误也可能由上游 Origin 查找异常导致；需要 UID 时可参考 [OriginUID 的替代查询方式](origin.md#限制与错误)。
