package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/wayne011872/ESGCheckSchedule/mail"
)

type MySQLConf struct {
	User      string `yaml:"user"`
	Pass      string `yaml:"pass"`
	Network   string `yaml:"network"`
	Server    string `yaml:"server"`
	Port      string `yaml:"port"`
	DefaultDB string `yaml:"defaul"`

	Dsn       string
}

func (mc *MySQLConf) GetDsn() string {
	mc.Dsn = fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",mc.User,mc.Pass,mc.Network,mc.Server,mc.Port,mc.DefaultDB)
	return mc.Dsn
}

func (mc *MySQLConf) NewMySQLDB()(*gorm.DB){
	mysqlDB,err := gorm.Open(mysql.Open(mc.GetDsn()),&gorm.Config{})
	if err != nil {
		mailContent := fmt.Sprintf("<h3><strong>--------傳承發票排程錯誤通知--------</strong></h3></br><p>以下為錯誤訊息 :%s</p></br>", err.Error())
		mail.SendMail("傳承發票排程錯誤通知", mailContent)
		panic("使用 gorm 連線 DB 發生錯誤，原因為 " + err.Error())
	}
	return mysqlDB
}