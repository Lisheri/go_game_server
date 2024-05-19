package db

// 利用xrom作为orm框架
import (
	"fmt"
	"log"
	"ms_sg_back/config"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

// 测试数据库连接以及初始化
func TestAndInitDB() {
	mysqlConfig, errInfo := config.File.GetSection("mysql")
	if errInfo != nil {
		log.Println("数据库配置缺少", errInfo)
		panic(errInfo)
	}
	dbConnectInfo := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"], 
		mysqlConfig["charset"])
	var err error
	Engine, err = xorm.NewEngine("mysql", dbConnectInfo)
	if err != nil {
		log.Println("数据库连接失败", err)
		panic(err)
	}
	if err = Engine.Ping(); err != nil {
		log.Println("数据库ping不通", err)
		panic(err)
	}
	// 第二个是key, 第三个参数是默认值, 表示返回值一定是int
	max_idle := config.File.MustInt("mysql", "max_idle", 2)
	Engine.SetMaxIdleConns(max_idle)
	max_open := config.File.MustInt("mysql", "max_conn", 10)
	Engine.SetMaxOpenConns(max_open)
	Engine.ShowSQL(true) // 是否展示sql
	log.Println("数据库连接完成", dbConnectInfo)
}
