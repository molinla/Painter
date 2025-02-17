package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	conf "painter-server-new/conf"
	"painter-server-new/tolog"
	"time"
)

var DbEngine *gorm.DB

func InitDB() error {
	dataSourceName := conf.MysqlConf.User + ":" + conf.MysqlConf.Password + "@tcp(" + conf.MysqlConf.Host + ":" + conf.MysqlConf.Port + ")/" + conf.MysqlConf.Database + "?charset=utf8&parseTime=true"

	configLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
		},
	)

	config := gorm.Config{}

	if conf.Server.Model == "debug" {
		config.Logger = configLog
	}

	db, err := gorm.Open(mysql.Open(dataSourceName), &config)
	if err != nil {
		tolog.Log().Errorf("Mysql init error %e", err).PrintAndWriteSafe()
		return err
	}
	DbEngine = db
	err = Migrate()
	if err != nil {
		tolog.Log().Errorf("Migrate user table error %e:", err).PrintAndWriteSafe()
		return err
	}
	tolog.Log().Infof("Connect to mysql: Success").PrintAndWriteSafe()
	InitSettings()
	InitRules()
	return nil
}

func GetDB() *gorm.DB {
	return DbEngine
}
