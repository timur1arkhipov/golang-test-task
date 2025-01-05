package app

import (
	"context"
	"fmt"
	"golangTestTask/internal/bootstrap"
	"golangTestTask/internal/config"
	"golangTestTask/internal/repositories/socksrepository/socksgorm"
	"golangTestTask/internal/services/socksservice"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

func Run(cfg config.Config) error {
	db, err := bootstrap.InitGormDB(cfg)
	if err != nil {
		return err
	}

	socksService := socksservice.New(socksgorm.New(db))
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Error while unwrapping gorm db: %s", err)
	}
	migrationsDir := filepath.Join(os.Getenv("PWD"), "migrations")
	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}
	n, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		fmt.Printf("Error while migrate: %s", err)
	}
	fmt.Printf("Applied %d migrations!\n", n)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler: socksService.GetHandler(),
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	gracefullyShutdown(ctx, cancel, server)
	return nil
}

func gracefullyShutdown(ctx context.Context, cancel context.CancelFunc, server *http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
	if err := server.Shutdown(ctx); err != nil {
		log.Warning(err)
	}
	cancel()
}
