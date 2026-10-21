# MapRotation：地图轮换

- **SDK 方法：** Client.MapRotation
- **HTTP 接口：** GET /maprotation
- **官网文档：** <https://apexlegendsapi.com/#map-rotation>
- **可运行示例：** [examples/map-rotation/main.go](../examples/map-rotation/main.go)

## 请求

`version=1` 或省略版本时只返回大逃杀普通模式；`version=2` 请求所有模式的地图轮换。这是 [API 官网](https://apexlegendsapi.com/#map-rotation) 明确说明的版本区别。

官网介绍仍提及竞技场和控制模式，但没有逐项公布当前响应字段。本文以下字段以用户提供的 2026-10-21、`Version: 2` 实际响应为依据，结合 EA 官方游戏说明解释模式含义，不将历史模式列表视为固定响应结构。

## 返回值

返回一个 JSON 对象，顶层按游戏模式分类，每个模式下包含 `current`（当前轮换）和 `next`（下一轮换）。SDK 保留完整原始响应，不修改地图名、时间或模式字段。

### 顶层模式

| 字段 | JSON 类型 | 游戏含义 | 本次响应的当前 → 下一轮换 |
|---|---|---|---|
| `battle_royale` | object | 大逃杀普通模式（非排位） | Broken Moon → World's Edge |
| `ranked` | object | 大逃杀排位模式 | World's Edge → Olympus |
| `ltm` | object | 本次对应 Mixtape 轮换，包含具体玩法和对应场地；不能仅凭键名认定其一定是临时活动 | Control / Hammond Labs → Gun Run / Skulltown |
| `wildcard` | object | Wildcard，规则更精简、节奏更快的独立模式 | Olympus → E-District |

上述 JSON 键及其映射来自本次响应。游戏背景可参见 [EA 游戏模式说明](https://help.ea.com/en/articles/apex-legends/game-modes/)：Ranked 与 Unranked 分别对应排位和非排位，Wildcard 有独立规则；[EA 对 Mixtape 的介绍](https://www.ea.com/ea-play/news/revelry-season-game-updates) 确认其轮换包含 Control、Gun Run 和 Team Deathmatch。样例里的 `ltm` 同时轮换玩法和地图，应同时读取 `eventName`、`map`。

这四个键是本次返回的模式，不代表所有版本、赛季都必定返回相同模式。读取时应允许模式增减；样例中没有独立的 `arenas` 或 `control` 顶层节点，Control 位于 `ltm` 内。

### current 和 next

| 字段 | JSON 类型 | 含义 |
|---|---|---|
| `<模式>.current` | object | 此次响应时刻的当前地图及其轮换时间段；包含剩余时间字段 |
| `<模式>.next` | object | 该模式紧接着的下一轮换；本次响应不包含剩余时间字段 |

`<模式>` 表示上述任一顶层模式。各模式有自己的时间表，不能用大逃杀的切图时间推算 Mixtape 的切换时间。

### 轮换条目的字段

下表列出 `current`、`next` 对象内的全部样例字段。类型来自实际响应，字段用途依据样例值及其相互关系解释；API 官网未提供这些字段的稳定 Schema。

| 字段 | JSON 类型 | 含义与示例 |
|---|---|---|
| `start` | number（整数） | 本轮换开始的 Unix 时间戳，单位为秒；如 `1792564200` |
| `end` | number（整数） | 本轮换结束的 Unix 时间戳，单位为秒；如 `1792569600` |
| `readableDate_start` | string | `start` 对应的可读时间，格式为 `YYYY-MM-DD HH:mm:ss`；样例为 UTC，如 `2026-10-21 06:30:00` |
| `readableDate_end` | string | `end` 对应的可读时间；样例为 UTC，如 `2026-10-21 08:00:00` |
| `map` | string | 地图或玩法场地的英文名称，如 `Broken Moon`、`Hammond Labs`；不一定是一张完整大逃杀地图 |
| `code` | string | 上游提供的地图/玩法轮换标识，如 `broken_moon_rotation`、`control_olympus_labs`；没有公开完整枚举，不建议自行拆分字符串推断规则 |
| `DurationInSecs` | number（整数） | 此轮换时间段的总秒数；样例中等于 `end - start`，如 `5400` |
| `DurationInMinutes` | number（整数） | 此轮换时间段的总分钟数，如 `90`；样例中等于 `DurationInSecs / 60` |
| `asset` | string | 地图图片 URL；可能是场地所属大地图的图片，如 Hammond Labs 使用 `Olympus.png` |
| `remainingSecs` | number（整数） | 此次响应时刻距当前轮换结束的剩余秒数，如 `1161` |
| `remainingMins` | number（整数） | 剩余分钟的整数部分；样例中 `1161` 秒对应 `19` 分钟，`261` 秒对应 `4` 分钟 |
| `remainingTimer` | string | 剩余时间的 `HH:mm:ss` 展示字符串，如 `00:19:21` |
| `isActive` | boolean | 活动/玩法的启用标记，样例中都为 `true`；更完整的取值语义未在官网公布，不能用来区分当前与下一轮换 |
| `eventName` | string | 具体玩法名称，如 `Control`、`Gun Run`；`map` 则是此玩法使用的场地 |

字段名大小写需与 JSON 一致，尤其是首字母大写的 `DurationInSecs`、`DurationInMinutes`。`next` 没有 `remaining*`，非 `ltm` 模式没有 `isActive`、`eventName`，不应把缺失字段直接解释为零、未启用或未知玩法。

同一 `code` 也可能用于不同模式：本次 `ranked.next` 和 `wildcard.current` 都是 `olympus_rotation`。若需要标识某个具体轮换，应结合模式键及时间段，不能只用 `code`。

### 时间与轮换时长

本次样例中，`1792569600` 对应 `2026-10-21 08:00:00 UTC`，即北京时间（`Asia/Shanghai`）`2026-10-21 16:00:00`。可读时间字符串没有携带时区，上述 UTC 解释由样例时间戳核对得出；展示和时间计算推荐使用 `start`、`end` 并显式转换时区。

| 模式 | 当前轮换区间（北京时间，2026-10-21） | 下一轮换区间（北京时间，同日） | 本次每个轮换的时长 |
|---|---|---|---|
| `battle_royale` | 14:30–16:00 | 16:00–17:30 | 5400 秒 / 90 分钟 |
| `ranked` | 11:30–16:00 | 16:00–20:30 | 16200 秒 / 270 分钟 |
| `ltm` | 15:30–15:45 | 15:45–16:00 | 900 秒 / 15 分钟 |
| `wildcard` | 14:30–16:00 | 16:00–17:30 | 5400 秒 / 90 分钟 |

这些是地图/玩法在匹配列表中开放的时间段，`DurationInMinutes` 不表示一场比赛要打多久。表中的 90、270、15 分钟仅描述这次响应，不应写成固定轮换规则。

本次每个模式都有 `current.end == next.start`。倒计时是服务端生成响应时的快照：例如大逃杀剩余 `1161` 秒，而 `ltm` 剩余 `261` 秒，相差 15 分钟，与二者的结束时间差一致。界面持续展示时可用 `end - 当前 Unix 时间` 更新倒计时；到期后重新获取轮换，不能只把旧倒计时重置成一个固定周期。

### JSON 示例

以下保留本次响应的 `ltm` 部分，完整展示共有字段及玩法专有字段；其他三个模式的地图和时间见上表。

~~~json
{
  "ltm": {
    "current": {
      "start": 1792567800,
      "end": 1792568700,
      "readableDate_start": "2026-10-21 07:30:00",
      "readableDate_end": "2026-10-21 07:45:00",
      "map": "Hammond Labs",
      "code": "control_olympus_labs",
      "DurationInSecs": 900,
      "DurationInMinutes": 15,
      "isActive": true,
      "eventName": "Control",
      "asset": "https://apexlegendsstatus.com/assets/maps/Olympus.png",
      "remainingSecs": 261,
      "remainingMins": 4,
      "remainingTimer": "00:04:21"
    },
    "next": {
      "start": 1792568700,
      "end": 1792569600,
      "readableDate_start": "2026-10-21 07:45:00",
      "readableDate_end": "2026-10-21 08:00:00",
      "map": "Skulltown",
      "code": "freedm_gungame_skulltown",
      "DurationInSecs": 900,
      "DurationInMinutes": 15,
      "isActive": true,
      "eventName": "Gun Run",
      "asset": "https://apexlegendsstatus.com/assets/maps/Arena_Skulltown.png"
    }
  }
}
~~~

`asset` 中的 `https:\/\/` 是 JSON 对 `/` 的合法转义，解码后与 `https://` 相同。图片路径中的 `Arena_` 也不能用来判定当前玩法，示例应以 `eventName: "Gun Run"` 为准。

## Go 解码

~~~go
response, err := client.MapRotation(ctx, &apex.MapRotationOptions{Version: 2})
if err != nil {
    return err
}

// 按模式键解码，便于兼容上游新增模式；仅声明业务需要的字段。
type entry struct {
    Map            string `json:"map"`
    Start          int64  `json:"start"`
    End            int64  `json:"end"`
    DurationInSecs int64  `json:"DurationInSecs"`
    EventName      string `json:"eventName"`
    IsActive       *bool  `json:"isActive"`
    RemainingSecs  *int64 `json:"remainingSecs"`
}
var rotation map[string]struct {
    Current entry `json:"current"`
    Next    entry `json:"next"`
}
if err := response.Decode(&rotation); err != nil {
    return err
}

if br, ok := rotation["battle_royale"]; ok {
    location, err := time.LoadLocation("Asia/Shanghai")
    if err != nil {
        return err
    }
    fmt.Printf("当前：%s，下一张：%s，切换时间：%s\n",
        br.Current.Map,
        br.Next.Map,
        time.Unix(br.Current.End, 0).In(location).Format(time.RFC3339))
}
~~~

时间展示代码需要导入 `fmt` 和 `time`。用 `*bool`、`*int64` 可区分字段缺失与实际返回的 `false`、`0`。仍可按原方式解码到 `map[string]any`；SDK 的 `Decode` 会将其中的数字保留为 `json.Number`。

## 限制与错误

传入版本之外的值会被 SDK 本地拒绝；上游错误码和鉴权方式见 [总览文档](../APEX_LEGENDS_API_ZH.md)。
