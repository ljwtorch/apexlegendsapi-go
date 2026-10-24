# PlayerByName：按名称查询玩家统计

- **SDK 方法：** <code>Client.PlayerByName</code>
- **HTTP 接口：** <code>GET /bridge</code>
- **官方接口说明：** <https://apexlegendsapi.com/documentation#query-by-name>
- **可运行示例：** [examples/player-by-name/main.go](../examples/player-by-name/main.go)

本文按接口返回 JSON 的结构解释字段。上游没有公布完整、固定的响应 Schema；字段和数据会随 API 版本、玩家资料、请求选项和上游更新变化。以下追踪器键名及个别嵌套字段来自提供的 <code>/bridge</code> 响应解析，适合作为读取指南，不应视为所有响应都必定包含的字段清单。

## 请求要点

<code>player</code> 和 <code>platform</code> 必填。PC 玩家应使用关联的 Origin 账号名称，即使玩家通过 Steam 启动游戏。当前文档列出的名称查询平台为 <code>PC</code>、<code>PS4</code>、<code>X1</code>。

官方支持的可选参数：

| 参数 | 作用 |
|---|---|
| <code>version</code> | 选择接口版本；当前默认版本为 5。 |
| <code>enableClubsBeta</code> | 尝试附带玩家俱乐部数据；该功能仍为 Beta，可能无数据。 |
| <code>skipRank</code> | 省略追踪器排名数据。 |
| <code>merge</code> | 将相同类型的追踪器合并，例如不同来源的击杀追踪器。 |
| <code>removeMerged</code> | 移除已合并的原始追踪器，只保留合并值。 |

<code>merge</code> 会影响 <code>legends</code> 与 <code>total</code> 中追踪器的键和值；比较不同请求的结果时需使用相同选项。<code>skipRank</code> 会影响下文的排名字段。

## 顶层结构

| 路径 | 含义 | 附件样例 |
|---|---|---|
| <code>global</code> | 账号资料、排位、通行证和账号徽章。 | 对象 |
| <code>realtime</code> | 查询时可取得的在线、对局和队伍状态。 | <code>currentState=offline</code> |
| <code>legends.selected</code> | 当前所选英雄及其展示追踪器、外观和图片。 | <code>LegendName=Pathfinder</code> |
| <code>legends.all</code> | 按英雄名分组的已记录追踪器及展示数据；<code>Global</code> 是全局追踪器组。 | <code>Global</code>、<code>Pathfinder</code> 等 |
| <code>total</code> | 按追踪器键汇总的统计项。每项通常包含 <code>name</code> 和 <code>value</code>。 | <code>career_kills=15281</code> |
| <code>mozambiquehere_internal</code> | 上游内部元数据；字段不属于稳定的业务数据契约。 | <code>isNewToDB=false</code>、<code>clusterSrv=fr-1</code> |
| <code>processingTime</code> | 上游处理本次请求所用的时间，单位为秒。 | <code>0.3147578239440918</code> |

## <code>global</code>：账号资料与排位

### 账号资料

| 路径 | 含义 | 附件样例 |
|---|---|---|
| <code>global.name</code> | 玩家展示名称。 | <code>PLAYER_NAME</code> |
| <code>global.tag</code> | 玩家标签；可能缺失。 | <code>PLAYER_TAG</code> |
| <code>global.uid</code> | 玩家唯一标识。 | <code>PLAYER_UID</code> |
| <code>global.avatar</code> | 头像 URL。 | <code>https://…/avatar.jpg</code> |
| <code>global.platform</code> | 该统计资料对应的平台。 | <code>PC</code> |
| <code>global.level</code> | 当前显示的账号等级。 | <code>309</code> |
| <code>global.toNextLevelPercent</code> | 当前等级进度，百分比数值。 | <code>44</code> |
| <code>global.levelPrestige</code> | 账号 Prestige 阶段。 | <code>2</code> |
| <code>global.internalUpdateCount</code> | 上游记录的内部更新计数。 | <code>27330</code> |
| <code>global.internalParsingVersion</code> | 上游数据解析版本号。 | <code>2</code> |

### 处罚信息：<code>global.bans</code>

| 路径 | 解释 | 附件样例 |
|---|---|---|
| <code>global.bans.isActive</code> | 处罚当前是否生效。 | <code>false</code> |
| <code>global.bans.remainingSeconds</code> | 剩余秒数。 | <code>0</code> |
| <code>global.bans.last_banReason</code> | 最近一次处罚的原因代码。 | <code>COMPETITIVE_DODGE_COOLDOWN</code> |

### 大逃杀排位：<code>global.rank</code>

| 路径 | 解释 | 附件样例 |
|---|---|---|
| <code>global.rank.rankScore</code> | 排位分数（RP）。 | <code>16448</code> |
| <code>global.rank.rankName</code> | 段位名称。 | <code>Master</code> |
| <code>global.rank.rankDiv</code> | 段位分级。 | <code>1</code> |
| <code>global.rank.ladderPosPlatform</code> | 玩家在当前平台排位榜上的位置。 | <code>4361</code> |
| <code>global.rank.rankImg</code> | 段位图标资源标识或 URL。 | <code>master1.png</code> |
| <code>global.rank.rankedSeason</code> | 排位模式与赛季阶段标识。 | <code>br_ranked_s30_s2</code> |
| <code>global.rank.rankedSeasonMeta.start</code> | 排位阶段开始时间，Unix 秒级时间戳。 | <code>1789491600</code> |
| <code>global.rank.rankedSeasonMeta.end</code> | 排位阶段结束时间，Unix 秒级时间戳。 | <code>1793728800</code> |
| <code>global.rank.ALStopPercent</code> | ALS 平台相关百分比字段。 | <code>0.24</code> |
| <code>global.rank.ALStopInt</code> | ALS 平台相关整数值。 | <code>13319</code> |
| <code>global.rank.ALStopPercentGlobal</code> | ALS 全局相关百分比字段。 | <code>0.31</code> |
| <code>global.rank.ALStopIntGlobal</code> | ALS 全局相关整数值。 | <code>18463</code> |
| <code>global.rank.ALSFlag</code> | ALS 标志位。 | <code>false</code> |

<code>ALStop*</code> 字段没有在公开 Schema 中给出足以可靠解码的完整定义。样例中的值如 <code>0.24</code> 和 <code>13319</code> 仅记录原始返回值；不要直接将其解释成“全球前 X%”或猎杀门槛。排名字段也可能是数字、字符串占位值或缺失。

### 其他全局对象

- <code>global.arena</code>：竞技场排位资料（历史模式）。附件样例中 <code>rankScore=0</code>、<code>rankName=Unranked</code>、<code>rankDiv=4</code>、<code>ladderPosPlatform=-1</code>、<code>rankImg=unranked4.png</code>、<code>rankedSeason=arenas_ranked</code>；ALS 字段返回 <code>No game this split</code>，<code>ALSFlag=true</code>。这些值说明应同时兼容数字和字符串。
- <code>global.battlepass.level</code>、<code>global.battlepass.history</code>：通行证当前等级和历史记录；附件样例两者均为 <code>null</code>。
- <code>global.badges[]</code>：账号徽章数组，字段见[徽章对象](#徽章对象)。附件中包含 <code>Ranked Season 29</code>（<code>value=112</code>）、<code>R-99 Mastery</code>（<code>value=564</code>）等记录。

## <code>realtime</code>：查询时状态

| 路径 | 解释 | 附件样例 |
|---|---|---|
| <code>realtime.lobbyState</code> | 大厅加入状态；单独看它不能证明玩家此刻在线或可加入。 | <code>open</code> |
| <code>realtime.isOnline</code> | 玩家是否在线。 | <code>0</code> |
| <code>realtime.isInGame</code> | 玩家是否正在对局。 | <code>0</code> |
| <code>realtime.canJoin</code> | 是否允许加入玩家队伍。 | <code>0</code> |
| <code>realtime.partyFull</code> | 玩家队伍是否已满。 | <code>0</code> |
| <code>realtime.selectedLegend</code> | 当前选择的英雄名，不等于玩家最常用或最擅长的英雄。 | <code>Pathfinder</code> |
| <code>realtime.currentState</code> | 机器可读状态码。 | <code>offline</code> |
| <code>realtime.currentStateAsText</code> | 上述状态的可读文本。 | <code>Offline</code> |
| <code>realtime.currentStateSinceTimestamp</code> | 状态开始时间戳；无有效时间时可能为 -1。 | <code>-1</code> |

状态标志应结合解释；离线时队伍加入能力等字段可能没有实际意义。由于它是查询时的上游状态，不能将其当作持续实时连接或游戏内遥测。

## <code>legends</code>：英雄追踪数据

### 英雄分组

- <code>legends.selected</code> 是本次响应中所选英雄的资料，可能包含 <code>LegendName</code>、<code>data</code>、<code>gameInfo</code>、<code>ImgAssets</code>。
- <code>legends.all.Global</code> 保存全局追踪器；其他键通常是英文英雄名。
- <code>legends.all.&lt;Legend&gt;.data[]</code> 是该英雄已记录的追踪器。不同英雄的数据项可能不同；缺少 <code>data</code> 不代表数值为零或从未使用该英雄。

附件样例中的 <code>legends.selected.LegendName</code> 是 <code>Pathfinder</code>。其展示追踪器包括 <code>BR Kills=2624</code>、<code>Grapple: Travel distance=289318</code> 和 <code>BR Damage=1101447</code>；其他英雄追踪器仍以各自的 <code>key</code> 区分。

不要把 <code>selectedLegend</code> 当作历史使用率，也不要假设所有英雄都有 <code>gameInfo</code>、徽章或统计数据。

### 追踪器对象：<code>data[]</code>

| 字段 | 解释 | 附件样例 |
|---|---|---|
| <code>name</code> | 展示名称，用于阅读，文字可能变化。 | <code>BR Kills</code> |
| <code>key</code> | 追踪器机器键名；优先用于程序识别统计类型。 | <code>specialEvent_kills</code> |
| <code>value</code> | 该追踪器值；保留原始 JSON 类型，勿假定所有统计都是整数。 | <code>2624</code> |
| <code>global</code> | 标记该追踪器是否为全局统计。 | <code>false</code> |
| <code>rank</code> | 综合排名对象。 | <code>{"rankPos":280791,"topPercent":10.43}</code> |
| <code>rankPlatformSpecific</code> | 平台范围排名对象。 | <code>{"rankPos":158859,"topPercent":9.16}</code> |
| <code>rank.rankPos</code> | 综合排名位置。 | <code>280791</code> |
| <code>rank.topPercent</code> | 综合排名百分比，基于上游统计样本；不是全体 Apex 玩家总体的精确百分位。 | <code>10.43</code> |
| <code>rankPlatformSpecific.rankPos</code> | 平台内排名位置。 | <code>158859</code> |
| <code>rankPlatformSpecific.topPercent</code> | 平台内排名百分比，同样以第三方样本为基础。 | <code>9.16</code> |

排名值可能返回数字，也可能是 <code>NOT_CALCULATED_YET</code> 等字符串（附件中的竞技场字段示例），或因 <code>skipRank</code> 而缺失。将这些情况与数值 0 区分处理。

### <code>gameInfo</code>：游戏内展示资料

样例中的 <code>gameInfo</code> 可包含已装备外观及其品质：

| 字段 | 含义 | 附件样例 |
|---|---|---|
| <code>skin</code>、<code>skinRarity</code> | 英雄皮肤及品质。 | <code>War Machine</code>、<code>Legendary</code> |
| <code>frame</code>、<code>frameRarity</code> | 旗帜边框及品质。 | <code>Disco Tech</code>、<code>Legendary</code> |
| <code>pose</code>、<code>poseRarity</code> | 旗帜姿势及品质。 | <code>Proud Machine</code>、<code>Rare</code> |
| <code>intro</code>、<code>introRarity</code> | 开场语音及品质。 | <code>Losing isn't fun</code>、<code>Rare</code> |
| <code>badges[]</code> | 英雄旗帜徽章。 | 含 <code>Apex Pathfinder</code>、<code>Assassin</code> 等徽章对象 |

<code>ImgAssets.icon</code> 和 <code>ImgAssets.banner</code> 是上游提供的英雄图标/横幅资源；<code>ImgAssets.banner</code> 不是玩家装备的 <code>gameInfo.frame</code>。

### 徽章对象

<code>global.badges[]</code> 和 <code>gameInfo.badges[]</code> 样例元素可包含：

| 字段 | 解释 | 附件样例 |
|---|---|---|
| <code>name</code> | 徽章名，通常为英文。 | <code>Apex Pathfinder</code> |
| <code>value</code> | API 返回的徽章内部值；可能代表样式、等级或编码，不能一概当作解锁数量或直接按数值判断徽章等级。 | <code>6</code> |
| <code>category</code> | 徽章所属分类，可能是英雄名称。 | <code>Pathfinder</code> |

附件样例中还出现 <code>Assassin</code>（<code>value=5</code>）、<code>Fool Proof</code>（<code>value=2</code>）以及字符串名称 <code>"null"</code>。徽章数组应按原始记录处理；不要假设名称唯一或把 <code>value</code> 直接当作完成场数。字符串 <code>"null"</code> 也不等同 JSON <code>null</code>。

## <code>total</code>：汇总统计

<code>total</code> 通常是以追踪器 <code>key</code> 为属性名的对象，例如：

    {
      "career_kills": {"name": "Career Kills", "value": 15281},
      "kd": {"name": "KD", "value": "-1"}
    }

| 结构 | 解释 |
|---|---|
| <code>total.&lt;key&gt;.name</code> | 统计项显示名称。 |
| <code>total.&lt;key&gt;.value</code> | 该统计项汇总值；可为数字、字符串占位值或其他上游类型。 |

附件响应中出现的追踪器键说明如下。该列表是样例中出现的键，不意味着每次都返回这些键：

| <code>key</code> | 含义 | 附件样例 <code>value</code> |
|---|---|---:|
| <code>career_kills</code> | 生涯击杀追踪器。 | <code>15281</code> |
| <code>career_revives</code> | 生涯扶起队友次数追踪器。 | <code>3572</code> |
| <code>career_wins</code> | 生涯胜场追踪器。 | <code>484</code> |
| <code>specialEvent_kills</code> | 大逃杀击杀追踪器；键名中的 <code>specialEvent</code> 不代表活动模式击杀。 | <code>8049</code> |
| <code>specialEvent_damage</code> | 大逃杀伤害追踪器。 | <code>4066168</code> |
| <code>specialEvent_wins</code> | 大逃杀胜场追踪器。 | <code>212</code> |
| <code>kills</code>、<code>damage</code> | 击杀追踪器和伤害追踪器。 | <code>kills=166</code>；<code>damage=229486</code> |
| <code>wins_season_21</code> | 第 21 赛季胜场追踪器。 | <code>4</code> |
| <code>mastery_car_damage_done</code> | C.A.R. 冲锋枪伤害追踪器。 | <code>956020</code> |
| <code>mastery_car_kills</code> | C.A.R. 冲锋枪击杀追踪器。 | <code>2295</code> |
| <code>mastery_r99_damage_done</code> | R-99 冲锋枪伤害追踪器。 | <code>772266</code> |
| <code>rev_ultimate_knocks_kills</code> | 亡灵终极技能相关击倒/击杀追踪器。 | <code>463</code> |
| <code>shields_charged</code> | 友方护盾充能相关追踪器。 | <code>7747</code> |
| <code>distance_on_stim</code> | 动力小子兴奋剂移动距离。 | <code>901388</code> |
| <code>enemies_scanned</code> | 寻血猎犬扫描敌人追踪器。 | <code>6734</code> |
| <code>grapple_travel_distance</code> | 探路者抓钩移动距离。 | <code>289318</code> |
| <code>tactical_tp_distance</code> | 罗芭战术技能传送距离。 | <code>79903</code> |
| <code>ultimate_loot_taken</code> | 罗芭黑市相关物品获取追踪器。 | <code>7971</code> |
| <code>tactical_focus_of_attention_hits</code> | 希尔战术技能命中追踪器。 | <code>632</code> |
| <code>ultimate_tempest_damage</code> | 弹道终极技能伤害追踪器。 | <code>62846</code> |
| <code>ultimate_energy_barricade_damage</code> | 导线管终极技能伤害追踪器。 | <code>4278</code> |
| <code>alter_tactical</code> | 变幻战术技能穿越表面次数追踪器。 | <code>762</code> |
| <code>kd</code> | 击杀死亡比（K/D）。 | <code>"-1"</code> |

相同 <code>name</code> 可能对应不同 <code>key</code>（例如 <code>BR Kills</code> 对应 <code>specialEvent_kills</code> 或 <code>kills</code>）。按 key 归一化业务统计时要处理重复追踪器和 <code>merge</code> 选项，避免重复相加。<code>total</code> 与单个英雄的 <code>data[]</code> 用途不同，也不能不加区分地把两者累加。

## <code>mozambiquehere_internal</code>

附件样例为 <code>{"isNewToDB":false,"clusterSrv":"fr-1"}</code>。这两个值是上游数据库/服务集群元数据；<code>clusterSrv</code> 不是玩家正在游玩的游戏服务器。通常仅用于诊断，不应暴露为玩家属性。

## 缺失值和解码建议

- 字段可能不存在、为 JSON <code>null</code>、是字符串占位值，或返回 <code>-1</code>；这些情况不等于数值 0。
- 对不稳定的动态结构，可先用 <code>map[string]any</code> 解码并按原始值分支处理；若使用强类型结构，为排名、徽章值等字段预留多种 JSON 类型。
- 保留未知 tracker key、徽章名和额外字段，以便上游新增内容时不丢数据。
- 该接口提供的是第三方追踪器/状态数据，并不保证覆盖完整生涯历史；不能仅凭本响应推导每场比赛记录、真实 K/D、游戏时长或英雄使用率。

## Go 解码

SDK 保留上游响应正文，调用方决定解码结构：

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

<code>Response.Decode</code> 使用 <code>json.Decoder.UseNumber</code>，因此解码到 <code>any</code> 的 JSON 数字为 <code>json.Number</code>，可避免大整数精度损失。

## 相关限制

这是非官方 API，无服务可用性保证。认证、速率限制及错误码请参阅[官方接口文档](https://apexlegendsapi.com/documentation)。
