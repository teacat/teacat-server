# 已棄用 DEPRECATED

此專案僅供研究或結構參考用途。

---

<p align="center">
  <img src="https://cloud.githubusercontent.com/assets/7308718/21562106/97c9ad20-ceb0-11e6-960a-664fa507bd68.png" alt="kitsvc" width="60">
  <br><br><strong>KitSvc</strong> 是一個 Go 的單個微服務初始包 <br>提供了 Gin、Consul、Prometheus、EventStore 相關模塊和 Gorm 與 NSQ。（<a href="https://github.com/TeaMeow/KitGate">建議搭配 KitGate</a>）
</p>

# 特色

- 採用 [Gin](https://github.com/gin-gonic/gin) 框架，以一般網站應用程式的方式處理請求。
- 具獨立性的微服務架構
- 透過 Gorm 與資料庫連線
- Consul、Prometheus、NSQ、EventStore
- 以 [Melody](https://github.com/olahol/melody) 作為支援 WebSocket 的框架
- 內建預設範例（亦有 WebSocket、NSQ 與 EventStore 和 [JSON Web Token](jwt.io) 的使用方式）

# 主要結構

KitSvc 是極具獨立性的微服務初始包，能作為單個微服務。基本使用方式與一般網站應用程式相同，但多了數個專門配合微服務的工具。

```js
KitSvc
├── client               // 客戶端
├── errno                // 錯誤代碼
├── model                // 資料結構
├── module
│  ├── event
│  │   ├── eventstore    // 事件邏輯
│  │   └── *event.go     // 事件函式介面
│  ├── metrics           // 效能測量
│  ├── mq
│  │   ├── mqstore       // 訊息佇列邏輯
│  │   └── *mq.go        // 訊息佇列函式介面
│  └── sd                // 服務探索
├── router
│  ├── middleware        // 中介軟體
│  └── *router.go        // 路由
├── server               // 主要邏輯
│  └── *main_test.go     // 單元測試
├── service              // 主要邏輯
├── shared               // 資源庫、共享函式、工具
├── store
│  ├── datastore         // 資料庫邏輯
│  └── *store.go         // 資料庫函式介面
├── vendor               // 依賴性套件
└── version              // 版本
```

# 依賴性

KitSvc 依賴下列服務，請確保你有安裝。

* **[Consul](https://www.consul.io/)**
* **[Prometheus](https://prometheus.io/)**
* **[NSQ](http://nsq.io/)**
* **[EventStore](https://geteventstore.com/)**
* **[MySQL](https://www.mysql.com/downloads/)（或 [MSSQL](https://www.microsoft.com/zh-tw/server-cloud/products/sql-server/overview.aspx)、[SQLite](https://sqlite.org/)、[PostgreSQL](https://www.postgresql.org/)）**

# 環境變數設置

為了方便在 Docker 上部署，KitSvc 沒有設定檔，而是透過環境變數來設置。下面這些設置都是預設值，如果你的環境符合這些預設值那麼你就不需要手動宣告，可以直接啟動程式進行測試。

```bash
# 服務名稱
KITSVC_NAME="Service"
# 服務暴露網址
KITSVC_URL="http://127.0.0.1:8080"
# 服務位置
KITSVC_ADDR="127.0.0.1:8080"
# 服務埠口
KITSVC_PORT=8080
# 服務註釋
KITSVC_USAGE="Operations about the users."
# JSON Web Token 的加密密碼
KITSVC_JWT_SECRET="4Rtg8BPKwixXy2ktDPxoMMAhRzmo9mmuZjvKONGPZZQSaJWNLijxR42qRgq0iBb5"
# Ping 伺服器的最大嘗試次數
KITSVC_MAX_PING_COUNT=20
# 除錯模式
KITSVC_DEBUG=false

# 資料庫驅動
KITSVC_DATABASE_DRIVER="mysql"
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
# 訊息產生者的 HTTP 位置
KITSVC_NSQ_PRODUCER_HTTP="127.0.0.1:4151"
# 訊息中心位置（以無空白 `,` 逗號新增多個位置）
KITSVC_NSQ_LOOKUPS="127.0.0.1:4161"

# 事件存儲中心的 HTTP 位置
KITSVC_ES_SERVER_URL="http://127.0.0.1:2113"
# 事件存儲中心帳號
KITSVC_ES_USERNAME="admin"
# 事件存儲中心密碼
KITSVC_ES_PASSWORD="changeit"

# 紀錄的命名空間
KITSVC_PROMETHEUS_NAMESPACE="service"
# 紀錄的服務名稱
KITSVC_PROMETHEUS_SUBSYSTEM="user"

# 服務中心的健康檢查時間
KITSVC_CONSUL_CHECK_INTERVAL="30s"
# 服務中心的健康檢查逾時時間
KITSVC_CONSUL_CHECK_TIMEOUT="1s"
# 服務中心的服務標籤（以無空白 `,` 逗號新增多個位置）
KITSVC_CONSUL_TAGS="user,micro"
```

# 模塊介紹

* **資料庫（Database）**：一個微服務擁有一個資料庫，同屬性的微服務可共享同一個資料庫，在這裡我們採用 Gorm 與資料庫連線。
* **服務探索（Discovery）**：向 Consul 服務中心註冊，表示自己可供使用，此舉是利於負載平衡做相關處理。
* **效能測量（Instrumenting）**：每個函式的執行時間、呼叫次數都可以被測量，並傳送給 Prometheus 伺服器彙整成視覺化資料。
* **紀錄層（Logging）**：傳入、輸出請求都可以被記錄，最終可以儲存程記錄檔供未來除錯。
* **訊息傳遞（Messaging）**：微服務之間並不會直接溝通，但可以透過 NSQ 訊息中心廣播讓另一個微服務處理相關事情，且無需等待該微服務處理完畢（即異步處理，不會阻擋）。
* **事件中心（EventStore）**：微服務之間可以發送、監聽事件，並在重新上線時重播事件找回資料庫的內容。

**非此套件**

這些功能會在 [KitGate](https://github.com/TeaMeow/KitGate/) 中實作，而不是 KitSvc。

* **版本控制（Versioning）**：能夠將舊的微服務替換掉，並且無須停機即可升級至新版本服務。
* **速率限制（Rate Limiting）**：避免客戶端在短時間內發送大量請求而導致癱瘓。

**未加入**

在 KitSvc 中這些功能沒有被加入，可能是計畫中。

* **斷路器（Circuit Breaker）**：斷路器能在服務錯誤時果斷拒絕外來請求，並給予一定的時間回復作業。

# 從這開始

請注意，KitSvc 並不能直接透過 `go get` 取得，因此你需要手動 `git clone` 一份回家。順帶記得的是：KitSvc 是一個開發用的模板，而不是套件。

```bash
# 設置環境變數。
$ export PATH=$PATH:$GOPATH/bin

# 從 Git 上複製一份此倉庫回家。
$ git clone git@github.com:TeaMeow/KitSvc.git $GOPATH/src/github.com/TeaMeow/KitSvc
$ cd $GOPATH/src/github.com/TeaMeow/KitSvc
```

## 可用指令

KitSvc 有下列指令可供開發環境時使用。

```bash
make deps                   # 將依賴性套件複製進 GOPATH
make test                   # 進行單元測試
make test_mysql test_sqlite # 測試資料庫功能
make run                    # 建置並執行程式
make build                  # 建置程式
```

如果有其他問題可以參考 `.drone.yml` 檔案，該檔案為自動測試環境的設置檔，你可以以此作為 KitSvc 所需的環境、設定依據。

## 索引

* [Client - 客戶端](#client)
* [Errno - 錯誤代碼](#errno)
* [Model - 資料結構](#model)
* [Event - 事件代碼](#event)
* [MQ - 訊息代碼](#mq)
* [Router - 路由處理](#router)
* [Server - 啟動核心](#server)
* [Service - 服務邏輯](#service)
* [Datastore - 資料庫邏輯](#datastore)
* [Version - 版本](#version)

## Client

```go
// GetUser gets an user by the user identifier.
func (c *client) GetUser(username string) (out *model.User, err error) {
	resp, err := grequests.Get(uri(pathSpecifiedUser, c.base, username), &grequests.RequestOptions{})
	resp.JSON(&out)
	return
}
```

`/client` 是客戶端，在這裡我們會撰寫透過 HTTP 和自己的程式溝通的函式。當需要進行單元測試時，我們就可以呼叫客戶端來測試自己的伺服器。如此一來就可以免去透過 Postman 逐一測試 API 的困擾。

## Errno

```go
var (
	ErrBind = &Err{
		Code:       "BIND_ERR",
		Message:    "Error occurred while binding the request body to the struct.",
		StatusCode: http.StatusBadRequest}
)
```

`/errno` 裡面存放的是錯誤回傳代碼，我們會在這裡定義程式的錯誤代碼、訊息、HTTP 狀態碼。之後發生錯誤時就可以回傳這個代碼。

## Model

```go
// User represents a registered user.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" gorm:"not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"not null" binding:"required" validate:"min=8,max=128"`
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *User) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)
	return
}
```

`/model` 就像是 MVC 結構中的 Model，但這裡存放的是物件定義，例如使用者建構體、文章建構體，我們亦會在這個地方替這些建構體新增一些輔助函式。

## Event

```go
var (
	EvtUserCreated = "user_created"
	EvtUserDeleted = "user_deleted"
)
```

`/module/event` 存放著我們會在此服務所用到的所有事件代碼，詳細的事件廣播範例請參考下方的 `/service` 段落。

## MQ

```go
var (
	MsgSendMail = "send_mail"
)
```

`/module/mq` 存放著我們在此服務會用到的訊息主題代碼，詳細的發布訊息範例請參考下方的 `/service` 段落。

## Router

```go
user := g.Group("/user")
{
  user.POST("", service.CreateUser)
  user.GET("/:username", service.GetUser)
  user.DELETE("/:id", service.DeleteUser)
  user.PUT("/:id", service.UpdateUser)
  user.POST("/token", service.PostToken)
}

// WebSockets.
w.Handle("/websocket", service.WatchUser)
// Message handlers.
m.Capture("user", mq.MsgSendMail, service.SendMail)
// Event handlers.
e.Capture(event.EvtUserCreated, service.UserCreated)
```

`/router` 用來定義路徑，我們亦會在此定義 WebSocket、訊息佇列、事件監聽器。

## Server

```go
func TestPostUser(t *testing.T) {
	assert := assert.New(t)

	u, err := c.PostUser(&model.User{
		Username: "admin",
		Password: "testtest",
	})
	assert.True(err == nil)
	err = u.Compare("testtest")
	assert.True(err == nil)
}
```

`/server` 是服務的主要核心，用來啟動服務的地方。我們亦會在此撰寫單元測式，而這個測試方式式呼叫客戶端與我們的程式溝通，然後判斷回傳的內容是否符合預期。

## Service

```go
// GetUser gets an user by the user identifier.
func GetUser(c *gin.Context) {
	// Get the `username` from the url parameter.
	username := c.Param("username")
	// Get the user by the `username` from the database.
	u, err := store.GetUser(c, username)
	if err != nil {
		errno.Abort(errno.ErrUserNotFound, err, c)
		return
	}

	c.JSON(http.StatusOK, u)
}
```

`/service` 是服務中最重要的邏輯地帶，這裡是用來處理服務邏輯的，我們亦會在此處理接收到的事件、訊息。

```go
// Publish the `send_mail` message to the message queue.
mq.Publish(c, mq.M{
  Topic: mq.MsgSendMail,
  Data:  &u,
})
// Send the `user_created` event to Event Store.
event.Send(c, event.E{
  Stream: event.EvtUserCreated,
  Data:   &u,
})
```

除了在此處理事件、訊息，我們也能夠透過上述方式廣播訊息、傳遞事件給其他服務。

## Datastore

```go
// GetUser gets an user by the user identifier.
func (db *datastore) GetUser(username string) (*model.User, error) {
	u := &model.User{}
	d := db.Where(&model.User{Username: username}).First(&u)
	return u, d.Error
}
```

`/store` 是用來存放專門存取資料庫的函式，我們亦會在此撰寫單元測試，這個測試只會測試和資料庫的存取，而不是整個服務的測試。

## Version

```go
var (
	// VersionMajor is for an API incompatible changes
	VersionMajor = 0
	// VersionMinor is for functionality in a backwards-compatible manner
	VersionMinor = 1
	// VersionPatch is for backwards-compatible bug fixes
	VersionPatch = 0
	// VersionDev indicates development branch. Releases will be empty string.
	VersionDev = "dev"
)
```

`/version` 裡面定義了這個服務的版本，這個資訊會在註冊至服務中心一併轉換成服務中心內的標籤。

# 還請參閱

這裡整理了一些也許能夠協助你理解微服務如何運作的文件。

**正體中文**

[一個基於 Golang 的基本 Go kit 微服務範例](https://yami.io/go-kit-example/)

[1. 什麼是微服務？——Golang 微服務實作教學與範例](https://yami.io/golang-microservice-1/)

[2. 微服務概念與溝通——Golang 微服務實作教學與範例](https://yami.io/golang-microservice-2/)

[3. 相關工具介紹與安裝——Golang 微服務實作教學與範例](https://yami.io/golang-microservice-3/)

**Go Kit 範例**

[go-kit/stringsvc3](https://github.com/go-kit/kit/tree/master/examples/stringsvc3)

[go-examples/nsq](https://github.com/ibmendoza/go-examples/tree/master/nsq)

**原文參考**

[An Introduction to Microservices, Part 1](https://auth0.com/blog/an-introduction-to-microservices-part-1/)

[API Gateway. An Introduction to Microservices, Part 2](https://auth0.com/blog/an-introduction-to-microservices-part-2-API-gateway/)

[An Introduction to Microservices, Part 3: The Service Registry](https://auth0.com/blog/an-introduction-to-microservices-part-3-the-service-registry/)

[Intro to Microservices, Part 4: Dependencies and Data Sharing](https://auth0.com/blog/introduction-to-microservices-part-4-dependencies/)

**啟發**

[drone/drone](https://github.com/drone/drone)

# License

MIT &copy; [Yami Odymel](https://github.com/YamiOdymel)
