package main

import (
	"github.com/xamust/couponApp/pkg/config"
	"github.com/xamust/couponApp/pkg/db/postgres"
	"github.com/xamust/couponApp/pkg/migration/goose"
	"log"
	"os"
)

/*
go run migrator status
go run migrator create filename sql
go run migrator up
go run migrator down
*/

func init() {
	if err := config.ViperInit(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("error parse config file: %v", err)
		return
	}
	db, err := postgres.Database(cfg).DB()
	if err != nil {
		panic(err)
	}
	migr := migrator.NewGooseMigrator(db)
	if err := migr.Commands(os.Args[1], os.Args[2:]...); err != nil {
		panic(err)
	}
}
