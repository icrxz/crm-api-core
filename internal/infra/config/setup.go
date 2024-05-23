package config

import (
	"os"
	"time"

	"github.com/magiconair/properties"
)

const (
	defaultAppPropertyFilename = "resources/application.properties"
	appPropertyFilenameEnv     = "APP_PROPERTY_FILE"
)

type AppConfig struct {
	Database     Database `properties:"database"`
	SecretJWTKey string   `properties:"jwtKeyEnv"`
}

type Database struct {
	ConnStr         string        `properties:"connStr,default="`
	Driver          string        `properties:"driver,default=postgres"`
	Port            int           `properties:"port,default=5432"`
	Schema          string        `properties:"schema"`
	Username        string        `properties:"username"`
	PasswordEnv     string        `properties:"passwordEnv"`
	HostEnv         string        `properties:"hostEnv"`
	ServerTimezone  string        `properties:"serverTimezone,default=America/Sao_Paulo"`
	MaxOpenConns    int           `properties:"maxOpenConns,default=50"`
	MaxIdleConns    int           `properties:"maxIdleConns,default=50"`
	ConnMaxLifetime time.Duration `properties:"connMaxLifetime,default=10m"`
}

func (db Database) Host() string {
	return os.Getenv(db.HostEnv)
}

func (db Database) Password() string {
	return os.Getenv(db.PasswordEnv)
}

func (db AppConfig) Password() string {
	return os.Getenv(db.SecretJWTKey)
}

func AppPropertyFilename() string {
	propFile := defaultAppPropertyFilename
	if f := os.Getenv(appPropertyFilenameEnv); f != "" {
		propFile = f
	}
	return propFile
}

func Load() (*AppConfig, error) {
	p := properties.MustLoadFile(AppPropertyFilename(), properties.UTF8)

	config := &AppConfig{}
	err := p.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
