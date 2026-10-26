# Predator：猎杀者门槛

- **SDK 方法：** Client.Predator
- **HTTP 接口：** GET /predator
- **官网文档：** <https://apexlegendsapi.com/#predator>
- **可运行示例：** [examples/predator/main.go](../examples/predator/main.go)

## 请求

此接口没有官网列出的查询参数，只需要 API Key。

## 返回值

返回各平台的 Apex Predator（顶尖猎杀者）分数门槛、用于确定门槛的参考账号信息，以及大师与猎杀者的合计人数。

根据 [API 官网](https://apexlegendsapi.com/#predator)，该接口用于获取 PC、PlayStation、Xbox、Switch 的猎杀者 RP/AP 门槛和大师人数。官网没有逐项公布响应字段的稳定 Schema；以下字段名和 JSON 类型依据用户提供的实际响应整理，未公开的内部取值逻辑会单独说明。

### 排位规则与门槛

[EA 排位说明](https://help.ea.com/en/articles/apex-legends/ranked/) 将 Apex Predator 定义为各平台大师段位中排名前 750 的玩家。因此猎杀者门槛会随榜单竞争变化，并非达到一个永久固定分数即可获得猎杀者身份；各平台的门槛应分别读取。

这份响应只有 `RP`，对应大逃杀排位积分（Ranked Points）。官网还提及 AP（竞技场积分），但不能据此假定当前响应一定包含 `AP`，也不应把 `RP` 下的 `val` 当作 AP 或其他评分。

### 顶层与平台字段

| 路径 | JSON 类型 | 含义 |
|---|---|---|
| `RP` | object | 大逃杀排位的各平台猎杀者门槛数据 |
| `RP.PC` | object | PC 平台的门槛与统计 |
| `RP.PS4` | object | PlayStation 平台的门槛与统计；`PS4` 是 API 的平台键，并非只指 PS4 主机 |
| `RP.X1` | object | Xbox 平台的门槛与统计；保留 API 的平台键 `X1` |
| `RP.SWITCH` | object | Switch 平台的门槛与统计 |

EA 说明同一主机平台的不同世代（如 PS4、PS5）按同一排位平台处理；平台键应使用 API 的原始拼写，不能自行替换成 `PS5`、`Xbox` 等展示名称。[EA 排位说明](https://help.ea.com/en/articles/apex-legends/ranked/)

### 各平台对象内的字段

下表中的字段分别位于 `RP.PC`、`RP.PS4`、`RP.X1`、`RP.SWITCH` 内。

| 字段 | JSON 类型 | 含义与示例 |
|---|---|---|
| `foundRank` | number（整数） | 结合字段名和样例，可理解为上游查到的门槛参考账号的榜单名次；PC、PS4 为 `750`，X1 为 `749`。它不是猎杀者人数，也不是调用者的排名；精确的参考账号选取算法未公开 |
| `val` | number（整数） | 该平台返回的猎杀者分数门槛参考值；在 `RP` 对象中单位为 RP，如 PC 的 `33795`。应结合 `foundRank`、`uid` 判断是否有有效参考账号 |
| `uid` | string | 与门槛参考记录对应的玩家 UID，如 `"1015157856349"`；它不是 API key 所有者的账号，也不是榜单名次。应保留为字符串；`"-1"` 不是有效参考玩家 UID |
| `updateTimestamp` | number（整数） | 该平台数据的更新时间，Unix 秒级时间戳，如 `1792569930`；不是毫秒，也不是每次客户端请求的时间 |
| `totalMastersAndPreds` | number（整数） | 该平台大师与猎杀者的合计统计人数，如 PC 的 `13581`；包括猎杀者，不是只统计未进入猎杀者的大师玩家 |

字段名区分大小写。Go 中 `foundRank` 要使用可表示负数的类型；`uid` 应使用 `string`，避免大整数在前端或其他语言中经过浮点转换丢失精度，并保留 `"-1"` 这样的特殊值。

### 样例数值的解读

| 平台 | `foundRank` | `val`（RP） | 大师与猎杀者合计 | 更新时间（北京时间） |
|---|---|---|---|---|
| `PC` | 750 | 33795 | 13581 | 2026-10-21 16:05:30 |
| `PS4` | 750 | 24404 | 6297 | 2026-10-21 16:05:56 |
| `X1` | 749 | 19785 | 2344 | 2026-10-21 16:06:15 |
| `SWITCH` | -1 | 16000 | 444 | 2026-10-21 16:06:15 |

PC 的记录可作为“查到第 750 名参考玩家时，其分数为 33795 RP”的门槛快照。X1 实际返回的是 `749`，应按原值展示；不能擅自改成 `750`，也不能据此推断 Xbox 的猎杀者名额只有 749 个。上游没有公开为什么该次选择了第 749 名。

Switch 同时返回 `foundRank: -1` 和 `uid: "-1"`，表示这条记录没有提供有效的参考名次和玩家 UID。`val: 16000` 仍是上游返回的门槛参考值，但不能当成某个第 750 名玩家的实测分数，也不能写成所有赛季通用的大师门槛。`444` 少于 750 与名额尚未满的情况相符，但这仅是可能的解释；人数统计可能低估，官网也没有定义这些特殊值的生成规则，不能据此断言该平台没有猎杀者或上游请求失败。

`totalMastersAndPreds` 应作为估计人数使用。[Apex Legends Status 的门槛页面](https://apexlegendsstatus.com/points-for-predator) 说明，人数基于其数据库内查到的最高大师榜单名次估算；未收录所有玩家时很可能低估。不能把该字段当作 EA 公布的精确实时总人数，也不应直接减去 750 推算纯大师人数。

各平台数据分别更新，样例中的时间戳并不完全相同。例如 `1792569930` 对应 `2026-10-21 08:05:30 UTC`，即北京时间 `2026-10-21 16:05:30`。建议展示更新时间，避免将旧门槛当作实时晋级保证；上游网站标注门槛页面每小时更新，但 API 官网未单独承诺此接口的固定刷新周期。[Apex Legends Status 门槛页面](https://apexlegendsstatus.com/points-for-predator)

### JSON 示例

以下为提供的完整响应；分数和人数仅表示该次快照。

~~~json
{
  "RP": {
    "PC": {
      "foundRank": 750,
      "val": 33795,
      "uid": "1015157856349",
      "updateTimestamp": 1792569930,
      "totalMastersAndPreds": 13581
    },
    "PS4": {
      "foundRank": 750,
      "val": 24404,
      "uid": "7199572041294906769",
      "updateTimestamp": 1792569956,
      "totalMastersAndPreds": 6297
    },
    "X1": {
      "foundRank": 749,
      "val": 19785,
      "uid": "1011866096780",
      "updateTimestamp": 1792569975,
      "totalMastersAndPreds": 2344
    },
    "SWITCH": {
      "foundRank": -1,
      "val": 16000,
      "uid": "-1",
      "updateTimestamp": 1792569975,
      "totalMastersAndPreds": 444
    }
  }
}
~~~

## Go 解码

~~~go
response, err := client.Predator(ctx)
if err != nil {
    return err
}

type threshold struct {
    FoundRank            int    `json:"foundRank"`
    Value                int64  `json:"val"`
    UID                  string `json:"uid"`
    UpdateTimestamp      int64  `json:"updateTimestamp"`
    TotalMastersAndPreds int64  `json:"totalMastersAndPreds"`
}
var data struct {
    RP map[string]threshold `json:"RP"`
}
if err := response.Decode(&data); err != nil {
    return err
}

if pc, ok := data.RP["PC"]; ok {
    fmt.Printf("PC 猎杀者门槛参考：%d RP，大师与猎杀者合计：%d\n",
        pc.Value, pc.TotalMastersAndPreds)
    fmt.Println("更新时间（UTC）：", time.Unix(pc.UpdateTimestamp, 0).UTC().Format(time.RFC3339))
    if pc.FoundRank <= 0 || pc.UID == "" || pc.UID == "-1" {
        fmt.Println("未提供有效的门槛参考账号")
    }
}
~~~

示例的展示代码需要导入 `fmt`、`time`。用平台映射读取数据可以兼容平台增减，读取单个平台前先检查键是否存在。仍可解码到 `map[string]any` 或 `any`，SDK 的 `Decode` 会将其中的数字保留为 `json.Number`。

## 限制与错误

平台字段名和返回包装可能随上游变化；请按实际 JSON 解码。通用错误码见 [总览文档](../APEX_LEGENDS_API_ZH.md)。
