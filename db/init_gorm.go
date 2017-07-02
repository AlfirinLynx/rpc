package db

import (
	"github.com/jinzhu/gorm"
	"sync"
	"github.com/antipin1987@gmail.com/rpcj/config"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	loadOnce sync.Once
	dbConn *gorm.DB
)

func DB() *gorm.DB {
	loadOnce.Do(newConnection)
	return dbConn
}

func Close() {
	if dbConn != nil {
		dbConn.Close()
	}
}

func newConnection() {
	var err error
	conf := config.Get().Sub("db")
	dbConn, err = gorm.Open(conf.GetString("dialect"), conf.GetString("dsn"))
	if err != nil {
		panic(err)
	}
}