package models

import (
	"api-server/config"
	"net/url"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetEnforcer() (*casbin.Enforcer, error) {
	var e *casbin.Enforcer

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
    m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*") || r.sub == "admin"
  `)
	if err != nil {
		log.Error().Err(err).Msg("NewModelFromString")
		return nil, err
	}

	e, err = casbin.NewEnforcer(m, a)
	if err != nil {
		log.Error().Err(err).Msg("NewEnforcer")
		return nil, err
	}

	if err = e.GetAdapter().(*adapter.Adapter).Transaction(e, func(e casbin.IEnforcer) error {
		if _, err = e.AddPolicy("user", "/*/logout", "POST"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Error().Err(err).Msg("Transaction")
		return nil, err
	}

	if err = e.LoadPolicy(); err != nil {
		log.Error().Err(err).Msg("LoadPolicy")
		return nil, err
	}

	return e, nil
}

func init() {
	u, err := url.Parse(config.G.DB.URL)
	if err != nil {
		log.Fatal().Err(err).Str("db", config.G.DB.URL).Msg("url.Parse")
	}

	var dialector gorm.Dialector
	switch u.Scheme {
	case "sqlite":
		dialector = sqlite.Open(u.Host)
	case "postgres":
		dialector = postgres.Open(config.G.DB.URL)
	default:
		dialector = mysql.Open(config.G.DB.URL)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("gorm.Open")
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
