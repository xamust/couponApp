package postgres

import (
	"fmt"
	"github.com/xamust/couponApp/pkg/config"
	"gorm.io/gorm"
	"log"
	"time"

	_ "github.com/lib/pq"
	postgresdriver "gorm.io/driver/postgres"
	gormlogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func Database(cfg *config.Config) *gorm.DB {
	if db != nil {
		return db
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Database, cfg.DB.SSL)

	loglevel := gormlogger.Silent
	if cfg.DB.Debug {
		loglevel = gormlogger.Info
	}

	newLogger := gormlogger.New(
		log.Default(),
		gormlogger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      loglevel,    // Log level
			Colorful:      false,       // Disable color
		},
	)

	conn, err := gorm.Open(postgresdriver.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to create a new postgres (v2) connection: %s", err.Error())
		return nil
	}

	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatalf("failed to get sql db: %s", err.Error())
		return nil
	}

	maxIdleConn := 1
	if cfg.DB.MaxIdle > 0 {
		maxIdleConn = cfg.DB.MaxIdle
	}
	sqlDB.SetMaxIdleConns(maxIdleConn)

	maxOpenConn := 1
	if cfg.DB.MaxOpen > 0 {
		maxOpenConn = cfg.DB.MaxOpen
	}
	sqlDB.SetMaxOpenConns(maxOpenConn)

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping to database: %s", err.Error())
		return nil
	}

	db = conn

	return db
}
