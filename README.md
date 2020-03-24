# go-gin-blog
A go example project

# day1
#### 一、初始化项目目录

go-gin-blog/

├── conf

├── middleware

├── models

├── pkg

├── routers

└── runtime

1、conf：用于存储配置文件

2、middleware：应用中间件

3、models：应用数据库模型

4、pkg：第三方包

5、routers 路由逻辑处理

6、runtime：应用运行时数据


#### 二、添加 Go Modules Replace
打开 go.mod 文件，新增 replace 配置项

？？？ 为什么要添加replace配置项？

答：在go.mod中我们可以看到，我们用的完整的外部模块引用路径，github.com/EDDYCJY/go-gin-example/xxx），而这个模块还没推送到远程，是没有办法下载下来的，因此需要用 replace 将其指定读取本地的模块路径，这样子就可以解决本地模块读取的问题
后续每新增一个本地应用目录，你都需要主动去 go.mod 文件里新增一条 replace



#### 三、编写配置包

1、拉去go-ini/ini 的依赖包：go get -u github.com/go-ini/ini

2、在conf目录下创建app.ini配置文件

3、建立调用配置的setting模块，在pkg目录下新建setting目录，新建setting.go文件，该文件的作用是引入app.ini配置文件，并进行一些初始化处理

#### 四、编写API错误包
1、建立错误码的e模块，在pkg目录下新建e目录，新建code.go和msg.go文件
，该模块的作用主要是接口的错误码和错误信息，以及相关函数的定义

#### 五、编写分页工具包
1、在pkg目录下新建util目录，并拉取com的依赖：go get -u github.com/unknwon/com

2、在util目录下新建pagination.go

#### 六、初始化mysql的models
1、拉取gorm的依赖包：go get -u github.com/jinzhu/gorm

2、拉取mysql驱动的依赖包：go get -u github.com/go-sql-driver/mysql

3、在models目录下新建models.go

#### 七、编写Demo

#### 八、优化路由
在routers目录下新疆router.go文件


# day2
#### 1、完成标签类的增删该接口的定义和编写
    a、获取标签列表：GET("/tags")

    b、新建标签：POST("/tags")

    c、更新指定标签：PUT("/tags/:id")

    d、删除指定标签：DELETE("/tags/:id")

#### 2、引用包

    a、beego-validation：本节采用的beego的表单验证库，中文文档：https://beego.me/docs/mvc/controller/validation.md。
    
    b、gorm，对开发人员友好的ORM框架，英文文档：http://gorm.io/docs/
    
    c、com，一个小而美的工具包
    
#### 3、完成文章类接口定义和编写
    a、获取文章列表：GET("/articles")
    
    b、获取指定文章：POST("/articles/:id")
    
    c、新建文章：POST("/articles")
    
    c、更新指定文章：PUT("/articles/:id")
    
    d、删除指定文章：DELETE("/articles/:id")

# day3
    1、加入权限用户校验

    2、加载jwt-go依赖包：go get -u github.com/dgrijalva/jwt-go
    加入token中间件验证（使用的JWT 文档[https://godoc.org/github.com/dgrijalva/jwt-go#SigningMethodHMAC]）
    
    3、测试通过
    
# day4
    1、自定义log，并接入接口
    2、使用endless包实现优雅重启
    
# day5
    1、golang应用部署docker
    2、修改配置文件app.ini，修改mysql的host配置项，改为mysql:3306，因为接下来我们启动应用容器时，会直接关联mysql容器
    其中--link 容器名:别名（如：--link mysql_db:mysql），这里用别名（mysql）即可访问link的容器
    2、编写Dockerfile文件，创建镜像：docker build -t go-gin-docker .
    3、启动容器，同时关联mysql容器
    docker run --link mysql_db:mysql --net backend_default  -p 8000:8000 gin-blog-docker
    
# day6
    1、重写GORM 的 Callbacks
    
# day7
    1、增加定时任务调度，使用依赖：go get -u github.com/robfig/cron




