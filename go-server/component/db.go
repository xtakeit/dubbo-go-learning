package component

import (
	"fmt"

	"go-server/library/clean"
	"go-server/library/mysql"
)

var DBContainer *mysql.DBContainer

type DBConfig struct {
	Name        string `env:"DB_NAME"`
	Host        string `env:"DB_HOST"`
	Port        string `env:"DB_PORT"`
	UserName    string `env:"DB_USERNAME"`
	Password    string `env:"DB_PASSWORD"`
	MaxLifeTime int    `env:"DB_MAX_LIFE_TIME"`
	MaxOpenConn int    `env:"DB_MAX_OPEN_CONN"`
	MaxIdleConn int    `env:"DB_MAX_IDLE_CONN"`
}

func SetupDB() (err error) {
	DBContainer, err = mysql.NewDBContainer(getDBConf)
	if err != nil {
		err = fmt.Errorf("mysql.NewDBContainer: %w", err)
		return
	}

	clean.Push(DBContainer)
	Conf.PushUpdater(DBContainer)

	return
}

func getDBConf() (cf *mysql.DBConf, err error) {
	cfg := &DBConfig{}

	if err = Conf.Scan(cfg, "env"); err != nil {
		err = fmt.Errorf("Conf.Scan: %w", err)
		return
	}

	cf = &mysql.DBConf{
		Name:        cfg.Name,
		Host:        cfg.Host,
		Port:        cfg.Port,
		UserName:    cfg.UserName,
		Password:    cfg.Password,
		MaxLifeTime: cfg.MaxLifeTime,
		MaxOpenConn: cfg.MaxOpenConn,
		MaxIdleConn: cfg.MaxIdleConn,
	}

	return
}
