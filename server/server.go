package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/canxium/supply-information/config"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server struct
type Server struct {
	echo *echo.Echo
	cfg  *config.Config
	db   *gorm.DB
}

// NewServer New Server constructor
func NewServer(cfg *config.Config) *Server {
	db, err := gorm.Open(postgres.Open(cfg.Server.Postgres), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	return &Server{echo: echo.New(), cfg: cfg, db: db}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           s.cfg.Server.Addr,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		fmt.Printf("Server is listening on: %s\n", s.cfg.Server.Addr)
		if err := s.echo.StartServer(server); err != nil {
			fmt.Printf("Error starting Server: %s\n", err.Error())
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	fmt.Println("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
