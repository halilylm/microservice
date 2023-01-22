package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"time"
)

type MysqlConn struct {
	DB                    *sql.DB
	host                  string
	port                  int
	user                  string
	password              string
	name                  string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
	log                   *zap.Logger
}

type MysqlConnOptions struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
	Log                   *zap.Logger
}

func NewMysqlConn(opts MysqlConnOptions) (*MysqlConn, error) {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	db := &MysqlConn{
		host:                  opts.Host,
		port:                  opts.Port,
		user:                  opts.User,
		password:              opts.Password,
		name:                  opts.Name,
		maxOpenConnections:    opts.MaxOpenConnections,
		maxIdleConnections:    opts.MaxIdleConnections,
		connectionMaxLifetime: opts.ConnectionMaxLifetime,
		connectionMaxIdleTime: opts.ConnectionMaxIdleTime,
		log:                   opts.Log,
	}
	if err := db.Connect(); err != nil {
		return nil, err
	}
	return db, nil
}

func (sd *MysqlConn) Connect() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sd.user, sd.password, sd.host, sd.port, sd.name))
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(sd.connectionMaxLifetime)
	db.SetConnMaxIdleTime(sd.connectionMaxIdleTime)
	db.SetMaxIdleConns(sd.maxIdleConnections)
	db.SetConnMaxLifetime(sd.connectionMaxLifetime)
	sd.DB = db
	return nil
}

func (sd *MysqlConn) Ping() error {
	return sd.DB.Ping()
}
