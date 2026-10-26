# OriginUID：Origin 名称查询 UID

- **SDK 方法：** Client.OriginUID
- **HTTP 接口：** GET /origin
- **官网文档：** <https://apexlegendsapi.com/#origin>
- **字段核对：** 2026-10-21
- **可运行示例：** [examples/origin/main.go](../examples/origin/main.go)

## 请求

必填参数是 player。该接口只查询 UID，不获取玩家统计；官网要求玩家曾经玩过 Apex Legends。

## 返回值

返回玩家的 Origin UID。官网没有说明顶层字段名或是否始终使用对象包装，因此 SDK 保留原始 JSON。

| 返回内容 | 含义 | 来源 |
|---|---|---|
| UID | 与 Origin 名称对应的稳定玩家标识 | 官方文档确认用途 |
| 未找到结果 | 玩家不存在、名称不匹配或没有 Apex 历史 | 结合官网错误码说明 |

脱敏语义示意：

~~~json
{
  "uid": "PLAYER_UID"
}
~~~

示意中的包装字段可能随上游变化，请以实际响应为准。

## Go 解码

~~~go
response, err := client.OriginUID(ctx, "ORIGIN_ACCOUNT_NAME")
if err != nil {
    return err
}

var data any
if err := response.Decode(&data); err != nil {
    return err
}
~~~

## 限制与错误

空名称会被 SDK 本地拒绝。404 需要结合响应正文判断，不能一律视为玩家不存在；其余错误码见 [总览文档](../APEX_LEGENDS_API_ZH.md)。

实测同一密钥查询同一玩家时，`/bridge` 能返回玩家 UID，而 `/origin` 返回 HTTP 404 和以下错误；另一个玩家的 `/origin` 查询也出现相同错误，说明上游 Origin 查找可能异常：

~~~json
{"Error": "Player not found. Please try again (err origin lookup)"}
~~~

SDK 将此响应识别为 `ErrPlayerNotFound`，可用 `errors.Is(err, apex.ErrPlayerNotFound)` 判断。该错误表示本次查找失败，不能据此断言账号不存在。无效密钥还可能返回 HTTP 404 和纯文本 `Error: API key doesn't exist !`，SDK 将其识别为 `ErrAccessDenied`。两者都保留 `APIError` 中的原始状态码和响应体。

上游 Origin 查找失败时，可显式使用玩家统计接口取得 UID：

~~~go
response, err := client.PlayerByName(ctx, apex.PlayerByNameRequest{
    Player:   "torchlee", // PC 使用绑定的 Origin 名称
    Platform: apex.PlatformPC,
})
if err != nil {
    return err
}
var player struct {
    Global struct {
        UID string `json:"uid"`
    } `json:"global"`
}
if err := response.Decode(&player); err != nil {
    return err
}
if player.Global.UID == "" {
    return fmt.Errorf("player response is missing global.uid")
}
fmt.Println(player.Global.UID)
~~~

上述替代方式会获取玩家统计；`OriginUID` 仍请求 `/origin` 并保留其原始响应，不会自动切换接口或伪造成功结果。
