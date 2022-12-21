package database

import (
	"fmt"
	"strings"
)

// MySQLConfig schema
type MySQLConfig struct {
	Host     string `json:"host" mapstructure:"host" yaml:"host"`
	Database string `json:"database" mapstructure:"database" yaml:"database"`
	Port     int    `json:"port" mapstructure:"port" yaml:"port"`
	Username string `json:"username" mapstructure:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" yaml:"password"`
	Options  string `json:"options" mapstructure:"options" yaml:"options"`
	Debug    bool   `json:"debug" mapstructure:"debug" yaml:"debug"`
}

// String return MySQL connection url
func (m MySQLConfig) String() string {
	return fmt.Sprintf("mysql://%s", m.DSN())
}

// DSN return Data Source Name
func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", m.Username, m.Password, m.Host, m.Port, m.Database, m.Options)
}

// PostgreSQLConfig hold config for postgresql
type PostgreSQLConfig struct {
	Host     string `json:"host" mapstructure:"host" yaml:"host"`
	Database string `json:"database" mapstructure:"database" yaml:"database"`
	Port     int    `json:"port" mapstructure:"port" yaml:"port"`
	Username string `json:"username" mapstructure:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" yaml:"password"`
	Options  string `json:"options" mapstrucuture:"options" yaml:"options"`
}

// String return Postgres connection string
// TODO: support ssl
func (m PostgreSQLConfig) String() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s%s", m.Username, m.Password, m.Host, m.Port, m.Database, m.Options)
}

// MySQLDefaultConfig return default config for mysql, usually use on development
func MySQLDefaultConfig() MySQLConfig {
	return MySQLConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "sample",
		Username: "root",
		Password: "khongbiet",
		Options:  "",
	}
}

//ElasticSearchConfig hold config
type ElasticsearchConfig struct {
	AddressesString     string `json:"addresses" mapstructure:"addresses"`
	MaxIdleConnsPerHost int    `json:"max_idle_conns_per_host" mapstructure:"max_idle_conns_per_host"`
	Timeout             int    `json:"timeout" mapstructure:"timeout"`
}

func (e ElasticsearchConfig) Addresses() []string {
	return strings.Split(e.AddressesString, ",")
}

func ElasticsearchDefaultConfig() ElasticsearchConfig {
	return ElasticsearchConfig{
		AddressesString:     "http://localhost:9200",
		MaxIdleConnsPerHost: 10,
		Timeout:             5000,
	}
}
