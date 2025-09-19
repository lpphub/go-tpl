package global

import (
	"fmt"
	"go-tpl/ext/logext"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDb() {
	var (
		err    error
		dbConf = Conf.Mysql

		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=True&loc=Asia%%2FShanghai",
			dbConf.User, dbConf.Password, dbConf.Addr, dbConf.Database)
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logext.NewGormLogger(dbConf.Database, dbConf.Addr, 1024), // 重写日志
	})

	if err != nil {
		panic("init db error: " + err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic("get db error: " + err.Error())
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
}
