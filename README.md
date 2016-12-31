## 從這開始

請注意，KitSvc 並不能透過 `go run` 直接執行，你必須先透過 `go build` 編譯後才能執行。執行前請先確認 NSQ 與 Consul 皆有啟動。

```bash
# clone（複製）或 fork（分歧）這個倉庫
$ cd KitSvc
$ go build

# 執行微服務
$ ./KitSvc
```

## 模塊介紹

* **資料庫（Database）**：一個微服務擁有一個資料庫，同屬性的微服務可共享同一個資料庫，在這裡我們採用 Gorm 與資料庫連線。
* **服務探索（Discovery）**：向 Consul 服務中心註冊，表示自己可供使用。
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
