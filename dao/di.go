package dao

import (
	"github.com/wayne011872/goSterna/log"
	"github.com/wayne011872/ESGCheckSchedule/db"
)

type Di struct {
	*db.MySQLConf	`yaml:"mysql,omitempty"`
	*log.LoggerConf `yaml:"log,omitempty"`
}