package models

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetEnforcer() (*casbin.Enforcer, error) {
	a, err := adapter.NewAdapterByDB(DB)
	if err != nil {
		log.Error().Err(err).Msg("NewAdapter")
		return nil, err
	}

	m, err := model.NewModelFromString(`
    [request_definition]
    r = sub, obj, act
    
    [policy_definition]
    p = sub, obj, act
    
    [role_definition]
    g = _, _
    
    [policy_effect]
    e = some(where (p.eft == allow))
    
    [matchers]
    m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
  `)
	if err != nil {
		log.Error().Err(err).Msg("NewModelFromString")
		return nil, err
	}

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Error().Err(err).Msg("NewEnforcer")
		return nil, err
	}

	_, _ = e.AddPolicy("admin", "*", "*")
	_, _ = e.AddPolicy("user", "/*/logout", "POST")

	if err = e.LoadPolicy(); err != nil {
		log.Error().Err(err).Msg("LoadPolicy")
		return nil, err
	}

	return e, nil
}

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("gorm.Open")
	}

	if err = DB.AutoMigrate(&User{}); err != nil {
		log.Error().Err(err).Msg("AutoMigrate")
	}

	DB.Model(&User{}).FirstOrCreate(&User{
		Username: "admin",
		Password: "admin",
		Role:     "admin",
		Email:    "yzg963@gmail.com",
	}, User{Username: "admin"})
}
