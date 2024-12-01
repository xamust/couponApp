package main

import (
	"context"
	"fmt"
	_ "github.com/xamust/couponApp/docs"
	v1 "github.com/xamust/couponApp/internal/adapter/api/v1"
	"github.com/xamust/couponApp/internal/adapter/db/postgre"
	"github.com/xamust/couponApp/internal/service/coupon_applier"
	"github.com/xamust/couponApp/internal/usecase"
	"github.com/xamust/couponApp/pkg/config"
	"github.com/xamust/couponApp/pkg/db/postgres"
	"github.com/xamust/couponApp/pkg/logger"
	_ "golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	if err := config.ViperInit(); err != nil {
		log.Fatal(err)
	}
}

//go:generate swag init --parseDependency --parseInternal --propertyStrategy pascalcase --parseDepth 3
func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("panic: %v", err)
		}
	}()
	ctx := context.Background()
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("error parse config file: %v", err)
		return
	}
	level, err := logger.ParseLevel(cfg.App.Log)
	if err != nil {
		log.Fatalf("error parse log level: %v", err)
		return
	}

	logger.NewLogger(logger.WithLevel(level))
	slog.InfoContext(ctx, fmt.Sprintf("config: %v", cfg))
	db := postgres.Database(cfg)

	userRepo := postgre.NewUserRepository(db)
	user := usecase.NewUserUsecase(userRepo)

	coupRepo := postgre.NewCouponRepository(db)
	coupRelation := postgre.NewCouponRelationRepository(db)
	coup := usecase.NewCouponUsecase(coupRepo, coupRelation)

	coup_app := coupon_applier.NewCouponApplier(
		coupRepo,
		userRepo,
		coupRelation,
		postgre.NewCouponApplierRepository(db))

	myRouter, err := v1.NewRouter(&user, &coup, &coup_app).Build(ctx)
	if err != nil {
		log.Fatalf("error build router: %v", err)
		return
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return myRouter.Start(cfg.App.Addr)
	})
	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return myRouter.Shutdown(ctx)
		case sig := <-quit:
			log.Println(ctx, fmt.Sprintf("recieved signal: [%+v]", sig))
			return myRouter.Shutdown(ctx)
		}
	})
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
