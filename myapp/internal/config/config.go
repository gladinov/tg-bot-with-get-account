package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	//Env                       string `yaml:"env"`
	ClientsHosts           Clients         `yaml:"clients"`
	DbType                 string          `yaml:"dbType"`
	StorageSQLLitePath     string          `yaml:"storageSQLLitePath"`
	MigrationsSqllitePath  string          `yaml:"migrationsSqllitePath"`
	MigrationsPostgresPath string          `yaml:"migrationsPostgresPath"`
	PostgresHost           PostgresHost    `yaml:"postgresHost"`
	RedisHTTPServer        RedisHTTPServer `yaml:"redis"`
}

type PostgresHost struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Port     int    `yaml:"port"`
	SslMode  string `yaml:"sslmode"`
}

type Clients struct {
	BondReportServiceHost string `yaml:"bondReportServiceHost"`
	TelegramHost          string `yaml:"telegramHost"`
	TinkoffApiHost        string `yaml:"tinkoffApiHost"`
}

type RedisHTTPServer struct {
	Address     string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

func (p *PostgresHost) GetStringHost() (string, error) {
	if p.Host == "" {
		return "", errors.New("empty host in config")
	}
	if p.User == "" {
		return "", errors.New("empty user in config")
	}
	if p.Password == "" {
		return "", errors.New("empty password in config")
	}
	if p.Dbname == "" {
		return "", errors.New("empty dbname in config")
	}
	if p.Port == 0 {
		return "", errors.New("empty port in config")
	}
	host := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s",
		p.Host,
		p.User,
		p.Password,
		p.Dbname,
		p.Port,
		p.SslMode)
	return host, nil
}

func (p *PostgresHost) GetHostToGoMigrate() (string, error) {
	if p.Host == "" {
		return "", errors.New("empty host in config")
	}
	if p.User == "" {
		return "", errors.New("empty user in config")
	}
	if p.Password == "" {
		return "", errors.New("empty password in config")
	}
	if p.Dbname == "" {
		return "", errors.New("empty dbname in config")
	}
	if p.Port == 0 {
		return "", errors.New("empty port in config")
	}
	host := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=%s", p.User, p.Password, p.Host, p.Port, p.Dbname, p.SslMode)
	return host, nil
}

func MustInitConfig(rootPath string) Config {
	path := filepath.Join(rootPath, "configs", "config.yaml")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", path)
	}

	var config Config
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return config
}
