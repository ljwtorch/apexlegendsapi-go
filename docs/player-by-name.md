# PlayerByName：按名称查询玩家统计

- **SDK 方法：** Client.PlayerByName
- **HTTP 接口：** GET /bridge
- **官网文档：** <https://apexlegendsapi.com/#query-by-name>
- **字段核对：** 2026-10-21
- **可运行示例：** [examples/player-by-name/main.go](../examples/player-by-name/main.go)

## 请求

必填参数是 player 和 platform。PC 玩家必须传绑定的 Origin 名称，即使实际通过 Steam 游玩；支持 PC、PS4、X1。可选的 version、enableClubsBeta、skipRank、merge 和 removeMerged 的含义见 [总览文档](../APEX_LEGENDS_API_ZH.md)。

## 返回值

返回一个 JSON 对象。官网没有发布完整 Schema；下面的字段表根据用户提供的 /bridge 样例整理，属于“样例观察”，实际字段会随版本、玩家数据和上游更新变化。

| 路径 | 类型 | 含义 |
|---|---|---|
| global | object | 账号概况、段位、战斗通行证和徽章等全局信息 |
| global.name | string | 玩家展示名称 |
| global.tag | string | 玩家标签；缺失时可能没有该字段 |
| global.uid | string | 玩家 UID |
| global.avatar | string | 头像 URL |
| global.platform | string | 统计所属平台 |
| global.level | number | 账号等级 |
| global.toNextLevelPercent | number | 距离下一级的进度百分比 |
| global.bans | object | 封禁状态 |
| global.bans.isActive | boolean | 当前是否处于封禁状态 |
| global.bans.remainingSeconds | number | 剩余封禁秒数 |
| global.bans.last_banReason | string | 最近一次封禁原因 |
| global.rank | object | 大逃杀排位信息 |
| global.rank.rankScore | number | 排位分数（RP） |
| global.rank.rankName | string | 段位名称 |
| global.rank.rankDiv | number | 段位分级 |
| global.rank.ladderPosPlatform | number | 平台排行榜位置；无排名时可能为负数 |
| global.rank.rankedSeason | string | 排位赛季标识 |
| global.rank.ALStopPercent / ALStopPercentGlobal | number | 相对平台或全球玩家的百分位 |
| global.rank.ALStopInt / ALStopIntGlobal | number | 相对平台或全球的比较值 |
| global.rank.ALSFlag | boolean | Apex Legends Status 的内部标记 |
| global.arena | object | 竞技场排位信息；没有该模式数据时可能是 null 或特殊字符串 |
| global.battlepass | object | 战斗通行证等级和历史；没有数据时可能为 null |
| global.badges | array | 账号徽章，元素通常含 name 和 value |
| global.levelPrestige | number | 传承等级 |
| realtime | object | 在线、组队和当前游戏状态 |
| realtime.isOnline / isInGame / canJoin / partyFull | number | 实时状态标志，样例中以 0 或 1 表示 |
| realtime.selectedLegend | string | 当前选择的传奇名称 |
| realtime.currentState / currentStateAsText | string | 机器可读和可读的当前状态 |
| realtime.currentStateSinceTimestamp | number | 状态开始时间的 Unix 时间戳 |
| legends.selected | object | 当前选择传奇的追踪器、游戏信息和图片资源 |
| legends.all | object | 按传奇名称分组的全部追踪器 |
| legends.*.data[] | array | 单项追踪器；常见字段为 name、value、key |
| legends.*.data[].rank | object | 全局排名和百分位 |
| legends.*.data[].rankPlatformSpecific | object | 平台内排名和百分位 |
| legends.*.gameInfo | object | 皮肤、姿势、语音和传奇徽章等信息 |
| legends.*.ImgAssets | object | 传奇图标和横幅 URL |
| mozambiquehere_internal | object | 上游内部处理信息，不建议作为业务契约 |
| total | object | 按 tracker key 汇总的统计；每项通常含 name 和 value |
| processingTime | number | 上游处理耗时，单位通常为秒 |

skipRank 会省略排名数据；merge 和 removeMerged 会改变 legends、total 中 tracker 的合并方式。字段缺失、null、负数和 No game this split 等值都可能是合法的无数据表示。

## 脱敏示例

下面只展示结构，真实玩家名称、UID、头像和数值已替换：

~~~json
{
  "global": {
    "name": "PLAYER_NAME",
    "uid": "PLAYER_UID",
    "platform": "PC",
    "level": 0,
    "rank": {
      "rankName": "Master",
      "rankScore": 0
    },
    "badges": []
  },
  "realtime": {
    "isOnline": 0,
    "currentState": "offline"
  },
  "legends": {
    "selected": {
      "LegendName": "LegendName",
      "data": []
    },
    "all": {}
  },
  "total": {},
  "processingTime": 0
}
~~~

## Go 解码

SDK 保留原始正文，不预定义上述字段：

~~~go
response, err := client.PlayerByName(ctx, apex.PlayerByNameRequest{
    Player:   "PLAYER_NAME",
    Platform: apex.PlatformPC,
})
if err != nil {
    return err
}

var data map[string]any
if err := response.Decode(&data); err != nil {
    return err
}
~~~

## 限制与错误

官网说明这是非官方服务，默认所有接口合计每秒 5 次请求。常见 HTTP 状态码及鉴权方式见 [总览文档](../APEX_LEGENDS_API_ZH.md)。收到 404 时通常表示玩家未找到；429 表示达到速率限制。
