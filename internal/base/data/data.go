package data

import (
	"fmt"

	"github.com/fntsky/ddl_guard/internal/base/conf"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

type Data struct {
	DB *xorm.Engine
}

func NewData(db *xorm.Engine) (*Data, error) {
	return &Data{
		DB: db,
	}, nil
}

func NewDB(debug bool) (*xorm.Engine, error) {
	cfg := conf.Global()
	if cfg == nil || cfg.Data.Database.Connection == "" {
		return nil, fmt.Errorf("config is not loaded or database config is empty")
	}
	dataBase := cfg.Data.Database
	if dataBase.Driver == "" {
		dataBase.Driver = "postgres"
	}
	engine, err := xorm.NewEngine(dataBase.Driver, dataBase.Connection)
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(debug)
	return engine, nil
}
