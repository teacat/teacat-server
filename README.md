<p align="center">
  <img src="https://cloud.githubusercontent.com/assets/7308718/21562106/97c9ad20-ceb0-11e6-960a-664fa507bd68.png" alt="kitsvc" width="60">
  <br><br><strong>KitSvc</strong> 是一個 Go 的單個微服務初始包 <br>提供了 Go kit、Consul、Prometheus 相關模塊和 Gorm 與 NSQ。（<a href="https://github.com/TeaMeow/KitGate">依賴 KitGate</a>）
</p>
# WIP 尚未完成
## 特色

- Go kit
- 具獨立性的微服務架構
- 透過 Gorm 與資料庫連線
- Consul、Prometheus、NSQ
- 內建預設範例

## 主要結構

在 KitSvc 中所有系統核心檔案都以 `core_*.go` 作為命名格式，這意味著你並不會動到他們。下面是一個你能夠編輯的資料結構。

```js
KitSvc
├── service
│  ├── *handlers.go        // 路由處理
│  ├── *instrumenting.go   // 效能測量層
│  ├── *logging.go         // 紀錄層
│  ├── *model.go           // 資料、邏輯處理
│  ├── *service.go         // 主服務功能
│  └── *transport.go       // 轉繼層和進入點
└── *config.yml            // 設定檔案
```

## 依賴性

KitSvc 依賴下列服務，請確保你有安裝。

* **[Consul](https://www.consul.io/)**
* **[Prometheus](https://prometheus.io/)**
* **[NSQ](http://nsq.io/)**
* **[MySQL](https://www.mysql.com/downloads/)（或 [MSSQL](https://www.microsoft.com/zh-tw/server-cloud/products/sql-server/overview.aspx)、[SQLite](https://sqlite.org/)、[PostgreSQL](https://www.postgresql.org/)）**

## 從這開始

請注意，KitSvc 並不能透過 `go run` 直接執行，你必須先透過 `go build` 編譯後才能執行。執行前請先確認 NSQ 與 Consul 皆有啟動。

```bash
# clone（複製）或 fork（分歧）這個倉庫
$ cd KitSvc
$ go build

# 執行微服務
$ ./KitSvc
```

## 環境變數設置

為了方便在 Docker 上部署，KitSvc 沒有設定檔，而是透過環境變數來設置。

```bash
# 服務名稱
KITSVC_NAME="StringService"
# 服務暴露網址
KITSVC_URL="http://127.0.0.1:8080"
# 服務位置
KITSVC_ADDR="127.0.0.1:8080"
# 服務埠口
KITSVC_PORT=8080
# 服務註釋
KITSVC_USAGE="Operations about the string."
# 服務版本
KITSVC_VERSION="0.0.1"

# 資料庫名稱
KITSVC_DATABASE_NAME="service"
# 資料庫主機位置與埠口
KITSVC_DATABASE_HOST="127.0.0.1:3306"
# 資料庫帳號
KITSVC_DATABASE_USER="root"
# 資料庫密碼
KITSVC_DATABASE_PASSWORD="root"
# 資料庫字符集
KITSVC_DATABASE_CHARSET="utf8"
# 資料庫時間地區
KITSVC_DATABASE_LOC="Local"
# 是否解析時間
KITSVC_DATABASE_PARSE_TIME=true

# 訊息產生者位置
KITSVC_NSQ_PRODUCER="127.0.0.1:4150"
# 訊息中心位置（以無空白 `,` 逗號新增多個位置）
KITSVC_NSQ_LOOKUPS="127.0.0.1:4161"

#
KITSVC_ES_SERVER_URL="http://127.0.0.1:2113"
KITSVC_ES_USERNAME="admin"
KITSVC_ES_PASSWORD="changeit"

# 紀錄的命名空間
KITSVC_PROMETHEUS_NAMESPACE="my_group"
# 紀錄的服務名稱
KITSVC_PROMETHEUS_SUBSYSTEM="string_service"

# 服務中心的健康檢查時間
KITSVC_CONSUL_CHECK_INTERVAL="10s"
# 服務中心的健康檢查逾時時間
KITSVC_CONSUL_CHECK_TIMEOUT="1s"
# 服務中心的服務標籤（以無空白 `,` 逗號新增多個位置）
KITSVC_CONSUL_TAGS="string,micro"
```

## 模塊介紹

* **資料庫（Database）**：一個微服務擁有一個資料庫，同屬性的微服務可共享同一個資料庫，在這裡我們採用 Gorm 與資料庫連線。
* **服務探索（Discovery）**：向 Consul 服務中心註冊，表示自己可供使用，此舉是利於負載平衡做相關處理。
* **效能測量（Instrumenting）**：每個函式的執行時間、呼叫次數都可以被測量，並傳送給 Prometheus 伺服器彙整成視覺化資料。
* **紀錄層（Logging）**：傳入、輸出請求都可以被記錄，最終可以儲存程記錄檔供未來除錯。
* **訊息傳遞（Messaging）**：微服務之間並不會直接溝通（很少數），但他們可以透過 NSQ 訊息中心廣播讓另一個微服務處理相關事情，且無需等待該微服務處理完畢（即異步處理，不會阻擋）。

### 缺少和未加入

這裡是目前尚未加入的功能，而有些功能則會在 [KitGate](https://github.com/TeaMeow/KitGate/) 中實作。

* **同步溝通（Synchronous Communication）**：與另一個微服務有所溝通，並取得該服務結束後回傳的資料接著繼續執行作業。
* **版本控制（Versioning）**：能夠將舊的微服務替換掉，並且無須停機即可升級至新版本服務。
* **速率限制（Rate Limiting）**：避免客戶端在短時間內發送大量請求而導致癱瘓。

## 從這開始

### 請求與回應

**服務處理**

`service/handlers.go`

`service/transport.go`

**回應內容**

```json
{
  "status" : "success",
  "code"   : "success",
  "message": "",
  "payload": {
    "username": "YamiOdymel"
  }
}
```

```json
{
  "status" : "error",
  "code"   : "str_empty",
  "message": "The string is empty.",
  "payload": null
}
```

### 資料與邏輯

`service/model.go`

### 中央控制器

`service/controller.go`

### 訊息與廣播

`service/handlers.go`

### 紀錄

`logging/service.go`

### 效能測量

`instrumenting/instrumenting.go`

## 單元測試

## 還請參閱

這裡整理了一些也許能夠協助你理解微服務如何運作的文件。

**正體中文**

[一個基於 Golang 的基本 Go kit 微服務範例](https://yami.io/go-kit-example/)

[1. 什麼是微服務？——Golang 微服務實作教學與範例](https://yami.io/golang-microservice-1/)

[2. 微服務概念與溝通——Golang 微服務實作教學與範例](https://yami.io/golang-microservice-2/)

[3. 相關工具介紹與安裝——Golang 微服務實作教學與範例](https://yami.io/golang-microservice-3/)

**官方範例**

[go-kit/stringsvc3](https://github.com/go-kit/kit/tree/master/examples/stringsvc3)

[go-examples/nsq](https://github.com/ibmendoza/go-examples/tree/master/nsq)

**原文參考**

[An Introduction to Microservices, Part 1](https://auth0.com/blog/an-introduction-to-microservices-part-1/)

[API Gateway. An Introduction to Microservices, Part 2](https://auth0.com/blog/an-introduction-to-microservices-part-2-API-gateway/)

[An Introduction to Microservices, Part 3: The Service Registry](https://auth0.com/blog/an-introduction-to-microservices-part-3-the-service-registry/)

[Intro to Microservices, Part 4: Dependencies and Data Sharing](https://auth0.com/blog/introduction-to-microservices-part-4-dependencies/)

# License

MIT &copy; [Yami Odymel](https://github.com/YamiOdymel)