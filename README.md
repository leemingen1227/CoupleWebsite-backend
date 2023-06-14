# couple website - backend

此專案為將 https://github.com/leemingen1227/CoupleWebsite-SideProject 之後端以 golang 開發之實作。專注於使用者之間配對的關係建立，和實作如何寄信邀請其他用戶成為配對關係。

## 專案講解

### 使用環境及工具

此專案採用Golang語言的Gin框架開發RESTful API，使用PostgreSQL作為資料庫，以Redis實作 task queue.

### 如何運行該專案(使用docker-compose)

可利用本專案的docker-compose.yaml會一次啟動Backend、PostgreSQL、Redis，方便直接運行測試。

總共起了以下四個Docker Container

+ postgres

  作為本專案的主要資料庫。

+ postgres-client

  這個是採用pgAdmin作為client端，開啟**localhost:5432**，並需要輸入帳號密碼: **couple@example.com*、**couple**，登入後連接postgres即可。

+ redis

    利用 asynq 套件實現異步任務，在專案中用於寄送認證郵件和邀請郵件。

+ backend

  此為Golang的Backend Service。

  可打開URL：http://localhost:8080/swagger/index.html ，利用Swagger框架打造RESTful API，可透過該框架直接測試API。


### 檢查環境變數是否一致

在本專案下有一個app.env.example檔案，裡面定義了本專案需要的環境變數。實際使用時請將檔名更改為 app.env


