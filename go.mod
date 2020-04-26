module github.com/Songkun007/go-gin-blog

go 1.13

require (
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/EDDYCJY/go-gin-example v0.0.0-20200322073714-2b22b57dfce9 // indirect
	github.com/astaxie/beego v1.12.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.5.0
	github.com/go-ini/ini v1.52.0
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/protobuf v1.3.4 // indirect
	github.com/gomodule/redigo v2.0.1-0.20180401191855-9352ab68be13+incompatible
	github.com/jinzhu/gorm v1.9.12
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/tealeg/xlsx v1.0.5
	github.com/unknwon/com v1.0.1
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace (
	github.com/Songkun007/go-gin-blog/conf => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/conf
	github.com/Songkun007/go-gin-blog/middleware => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/middleware
	github.com/Songkun007/go-gin-blog/models => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/models
	github.com/Songkun007/go-gin-blog/pkg/app => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/pkg/app
	github.com/Songkun007/go-gin-blog/pkg/e => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/pkg/e
	github.com/Songkun007/go-gin-blog/pkg/file => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/pkg/file
	github.com/Songkun007/go-gin-blog/pkg/gredis => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/pkg/gredis
	github.com/Songkun007/go-gin-blog/pkg/setting => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/pkg/setting
	github.com/Songkun007/go-gin-blog/pkg/upload => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/pkg/upload
	github.com/Songkun007/go-gin-blog/pkg/util => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/pkg/util
	github.com/Songkun007/go-gin-blog/routers => /data1/htdocs/go_project/src/github.com/Songkun007/go-gin-blog/routers
)
