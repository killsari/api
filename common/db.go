package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/url"
	"github.com/killsari/api/config"
)

var DB *sqlx.DB

func init() {
	mysqlConfig := config.Conf.Mysql
	db, err := sqlx.Connect(
		"mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=%s",
			mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database,
			url.QueryEscape("Asia/Shanghai")),
	) //with ping
	if err != nil {
		logrus.Error(err)
	}
	DB = db
}
