# Servers：服务器状态

- **SDK 方法：** Client.Servers
- **HTTP 接口：** GET /servers
- **官网文档：** <https://apexlegendsapi.com/#server-status>
- **可运行示例：** [examples/servers/main.go](../examples/servers/main.go)

## 请求

此接口没有官网列出的查询参数，只需要 API Key。

## 返回值

返回 Apex Legends Status 的服务状态监测快照。[API 官网](https://apexlegendsapi.com/#server-status) 说明，此接口提供其网站展示的服务器状态；具体字段以下方提供的实际响应为依据，官网未发布完整的稳定 Schema。

### 响应结构与服务分类

响应采用两层对象：`服务分类 → 地区或组件 → 状态记录`。第二层不一定是地区，也可能是上游网站的自检组件或外部平台服务。

| 顶层字段 | JSON 类型 | 含义 |
|---|---|---|
| `Origin_login` | object | Origin/EA 体系的 PC 登录相关服务，各地区分别监测；`Origin` 是 API 保留的键名，不能据此认定只涉及旧版 Origin 客户端 |
| `EA_novafusion` | object | EA Novafusion 服务的地区监测结果；API 官网未公开该服务的具体内部职责，不应直接解释为全部游戏对局服务器状态 |
| `EA_accounts` | object | EA 玩家账号相关服务的地区监测结果 |
| `ApexOauth_Crossplay` | object | Apex 跨平台认证相关服务的地区监测结果；不是跨平台匹配成功率或玩家的跨平台设置 |
| `selfCoreTest` | object | Apex Legends Status（ALS）自身网站、统计查询及相关接口组件的自检结果 |
| `otherPlatforms` | object | PlayStation Network、Xbox Live 等外部平台的状态汇总 |

这些分类结合实际键名与 [ALS 状态网站](https://apexlegendsstatus.com/) 的 PC/Desktop logins、Player accounts、Crossplay auth、ALS website 和 PSN/Xbox Live 栏目解释。ALS 是独立监测服务；需要核对游戏或平台的官方状态时，可查看 [EA 官方服务器状态页](https://help.ea.com/en/server-status/)。

### 地区键

前四类服务在样例中都包含以下地区。它们是上游的地区监测分组，不是玩家当前所在地区、游戏内具体数据中心名称或可直接连接的服务器地址。

| 地区键 | 含义 |
|---|---|
| `EU-West` | 欧洲西部 |
| `EU-East` | 欧洲东部 |
| `US-West` | 美国西部 |
| `US-Central` | 美国中部 |
| `US-East` | 美国东部 |
| `SouthAmerica` | 南美洲 |
| `Asia` | 亚洲 |

### selfCoreTest 的组件

| 组件键 | 含义 | 样例状态 |
|---|---|---|
| `Status-website` | ALS 状态网站自身的可用性检查 | `UP`，HTTP 200，92 毫秒 |
| `Stats-API` | ALS 玩家统计查询接口的检查 | `SLOW`，HTTP 200，5145 毫秒 |
| `Overflow-#1` | ALS 名为 Overflow #1 的组件；具体部署及流量调度方式未公开 | `DOWN`，HTTP 503，耗时 `-1` |
| `Overflow-#2` | ALS 名为 Overflow #2 的组件；不能据名称推断其与 #1 的切换关系 | `UP`，HTTP 200，1 毫秒 |
| `Origin-API` | ALS 自检中的 Origin 相关查询接口 | `SLOW`，HTTP 200，6272 毫秒 |
| `Playstation-API` | ALS 自检中的 PlayStation 相关查询接口 | `DOWN`，HTTP 404，耗时 `-1` |
| `Xbox-API` | ALS 自检中的 Xbox 相关查询接口 | `SLOW`，HTTP 200，991 毫秒 |

`selfCoreTest.Origin-API` 与 `Origin_login` 属于不同检测项目，不能视作同一个端点。`Playstation-API` 故障也不等于整个 PlayStation Network 故障；样例中 `otherPlatforms.Playstation-Network` 仍为 `UP`。

### otherPlatforms 的平台

| 平台键 | 含义 |
|---|---|
| `Playstation-Network` | PlayStation Network 的状态汇总 |
| `Xbox-Live` | Xbox Live 的状态汇总 |

样例中这两个平台只有 `Status`、`QueryTimestamp`，没有 `HTTPCode`、`ResponseTime`。这是字段缺失，不能补成 HTTP 0 或 0 毫秒。

### 状态记录字段

| 字段 | JSON 类型 | 含义与示例 |
|---|---|---|
| `Status` | string | 上游对该节点的健康状态判断；样例出现 `UP`、`SLOW`、`DOWN` |
| `HTTPCode` | number（整数，可缺失） | 上游监测该节点时记录的 HTTP 状态码，如 `200`、`503`、`404`；不是本次调用 `/servers` 的 HTTP 状态码 |
| `ResponseTime` | number（整数，可缺失） | 上游记录的节点响应耗时，按 ALS 网站的毫秒单位解释，如 `5145` 约为 5.145 秒；不是本地程序的请求总耗时，也不是玩家的游戏内 Ping |
| `QueryTimestamp` | number（整数） | 该节点状态的采集/查询时间，Unix 秒级时间戳，如 `1792575009`；各节点可以不同，不是游戏开服或故障开始时间 |

字段名区分大小写，必须保留 `Status`、`HTTPCode`、`ResponseTime`、`QueryTimestamp` 的原始拼写。

### 状态与特殊值

| `Status` | 解释 | 样例 |
|---|---|---|
| `UP` | 上游本次监测认为该节点正常 | `Origin_login.EU-West`，HTTP 200，28 毫秒 |
| `SLOW` | 节点能够响应，但上游判定响应偏慢 | `selfCoreTest.Stats-API`，HTTP 200，5145 毫秒 |
| `DOWN` | 上游判定该节点异常或不可用 | `selfCoreTest.Overflow-#1`，HTTP 503 |

以上是样例中出现的状态，官网没有公布完整枚举或统一的慢响应阈值。`Xbox-API` 的 991 毫秒也被标为 `SLOW`，不应自行规定“超过 1000 毫秒才算慢”，也不应仅根据 `HTTPCode == 200` 将 `SLOW` 改为 `UP`。遇到未识别的状态值应保留原值。

`ResponseTime: -1` 在样例中与 `DOWN` 一起出现，应按无有效耗时数据的特殊值处理，不能理解为负延迟。它也不能单独证明是连接超时，因为样例同时提供了 HTTP 503 或 404。`ResponseTime: 0` 是上游实际返回的数值，不等于字段缺失，也不代表网络传输完全不耗时；其计时精度、取整或缓存细节未公开。

调用 `/servers` 成功，只表示取得了状态快照；JSON 内某个节点的 `HTTPCode: 404` 或 `503` 是被监测节点的结果，不会让 `client.Servers` 因该内层字段直接返回 `APIError`。反过来，节点 `UP` 也不能保证某个玩家一定能登录或进入对局，仍可能存在账号、具体数据中心或本地网络问题。

### 时间解读

`1792575001` 对应 `2026-10-21 09:30:01 UTC`，即北京时间（`Asia/Shanghai`）`2026-10-21 17:30:01`。样例中各节点时间从 `1792575001` 到 `1792575024`，相差 23 秒，不能把整份响应当成完全同一时刻的测量结果，也不能凭这些时间差断言具体探测执行顺序。

展示状态时应同时展示节点的更新时间，避免把缓存快照当作实时结果；API 官网没有给出此接口的固定刷新间隔或故障恢复时间承诺。

### JSON 示例

以下节选保留六类服务，并在地区服务中各保留一个地区；`selfCoreTest` 和 `otherPlatforms` 保留全部节点。数值只描述提供的那次快照。

~~~json
{
  "Origin_login": {
    "EU-West": {
      "Status": "UP",
      "HTTPCode": 200,
      "ResponseTime": 28,
      "QueryTimestamp": 1792575002
    }
  },
  "EA_novafusion": {
    "EU-West": {
      "Status": "UP",
      "HTTPCode": 200,
      "ResponseTime": 63,
      "QueryTimestamp": 1792575003
    }
  },
  "EA_accounts": {
    "EU-West": {
      "Status": "UP",
      "HTTPCode": 200,
      "ResponseTime": 30,
      "QueryTimestamp": 1792575003
    }
  },
  "ApexOauth_Crossplay": {
    "EU-West": {
      "Status": "UP",
      "HTTPCode": 200,
      "ResponseTime": 75,
      "QueryTimestamp": 1792575004
    }
  },
  "selfCoreTest": {
    "Status-website": {
      "Status": "UP",
      "HTTPCode": 200,
      "ResponseTime": 92,
      "QueryTimestamp": 1792575004
    },
    "Stats-API": {
      "Status": "SLOW",
      "HTTPCode": 200,
      "ResponseTime": 5145,
      "QueryTimestamp": 1792575009
    },
    "Overflow-#1": {
      "Status": "DOWN",
      "HTTPCode": 503,
      "ResponseTime": -1,
      "QueryTimestamp": 1792575009
    },
    "Overflow-#2": {
      "Status": "UP",
      "HTTPCode": 200,
      "ResponseTime": 1,
      "QueryTimestamp": 1792575009
    },
    "Origin-API": {
      "Status": "SLOW",
      "HTTPCode": 200,
      "ResponseTime": 6272,
      "QueryTimestamp": 1792575016
    },
    "Playstation-API": {
      "Status": "DOWN",
      "HTTPCode": 404,
      "ResponseTime": -1,
      "QueryTimestamp": 1792575020
    },
    "Xbox-API": {
      "Status": "SLOW",
      "HTTPCode": 200,
      "ResponseTime": 991,
      "QueryTimestamp": 1792575021
    }
  },
  "otherPlatforms": {
    "Playstation-Network": {
      "Status": "UP",
      "QueryTimestamp": 1792575021
    },
    "Xbox-Live": {
      "Status": "UP",
      "QueryTimestamp": 1792575024
    }
  }
}
~~~

## Go 解码

~~~go
response, err := client.Servers(ctx)
if err != nil {
    return err
}

type nodeStatus struct {
    Status         string `json:"Status"`
    HTTPCode       *int   `json:"HTTPCode"`
    ResponseTime   *int64 `json:"ResponseTime"`
    QueryTimestamp int64  `json:"QueryTimestamp"`
}
// 第一层是服务分类，第二层可能是地区、组件或外部平台。
var status map[string]map[string]nodeStatus
if err := response.Decode(&status); err != nil {
    return err
}

for service, nodes := range status {
    for name, node := range nodes {
        fmt.Printf("%s / %s：%s", service, name, node.Status)
        if node.HTTPCode != nil {
            fmt.Printf("，HTTP %d", *node.HTTPCode)
        }
        switch {
        case node.ResponseTime == nil:
            fmt.Print("，未提供耗时")
        case *node.ResponseTime < 0:
            fmt.Print("，无有效耗时")
        default:
            fmt.Printf("，%d 毫秒", *node.ResponseTime)
        }
        fmt.Println("，采集时间（UTC）：",
            time.Unix(node.QueryTimestamp, 0).UTC().Format(time.RFC3339))
    }
}
~~~

展示代码需要导入 `fmt`、`time`。指针字段可以区分“未提供 HTTPCode/ResponseTime”与实际返回 `0` 的情况，带符号整数可保留 `-1`。使用两层映射可以适应地区、组件和平台增减；上游新增状态字符串时也可以直接展示。仍可解码到 `any` 或 `map[string]any`，其中数字由 SDK 的 `Decode` 保留为 `json.Number`。

## 限制与错误

展示来自该接口的数据时，必须链接 <https://apexlegendsstatus.com>，或注明 Data from apexlegendsstatus.com。字段和状态枚举以实际响应为准；通用错误码见 [总览文档](../APEX_LEGENDS_API_ZH.md)。
