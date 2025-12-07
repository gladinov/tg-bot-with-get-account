package lib

import (
	"errors"
	"fmt"
)

type PostgresHost struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Port     int    `yaml:"port"`
	SslMode  string `yaml:"sslmode"`
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
