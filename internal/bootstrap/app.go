package bootstrap

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-admin/server/internal/config"
	"go-admin/server/internal/handler"
	"go-admin/server/internal/repository"
	"go-admin/server/internal/router"
	"go-admin/server/internal/service"
	"go-admin/server/pkg/logger"
)

type App struct {
	server *http.Server
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	log := logger.New()
	log.Printf("mysql addr: %s", cfg.MySQLDSN)
	log.Printf("redis addr: %s", cfg.RedisAddr)

	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("init mysql db handle: %w", err)
	}
	if err := sqlDB.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf(
			"ping redis addr=%s db=%d password_set=%t: %w; check redis protected-mode/bind/requirepass and server firewall",
			cfg.RedisAddr,
			cfg.RedisDB,
			strings.TrimSpace(cfg.RedisPassword) != "",
			err,
		)
	}

	repo := repository.New(db)
	authService := service.NewAuthService(repo, redisClient, cfg)
	rbacService := service.NewRBACService(repo)
	systemService := service.NewSystemService(repo)
	publicService := service.NewPublicService(repo)

	engine := router.New(
		cfg,
		repo,
		redisClient,
		handler.NewAuthHandler(authService, rbacService),
		handler.NewSystemHandler(systemService, cfg),
		handler.NewPublicHandler(publicService),
		rbacService,
	)
	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("starting %s on :%s", cfg.AppName, cfg.AppPort)
	log.Printf("swagger ui: http://127.0.0.1:%s/swagger/index.html", cfg.AppPort)
	log.Printf("swagger yaml(generated): http://127.0.0.1:%s/api-docs/swagger.yaml", cfg.AppPort)
	if hostIP := detectHostIP(); hostIP != "" {
		log.Printf("swagger public: http://%s:%s/swagger/index.html", hostIP, cfg.AppPort)
		log.Printf("swagger yaml public(generated): http://%s:%s/api-docs/swagger.yaml", hostIP, cfg.AppPort)
	} else {
		log.Printf("swagger public: http://<your-server-ip>:%s/swagger/index.html", cfg.AppPort)
		log.Printf("swagger yaml public(generated): http://<your-server-ip>:%s/api-docs/swagger.yaml", cfg.AppPort)
	}
	return &App{server: server}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func detectHostIP() string {
	conn, err := net.DialTimeout("udp", "8.8.8.8:80", 2*time.Second)
	if err != nil {
		return ""
	}
	defer conn.Close()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok || addr.IP == nil {
		return ""
	}
	ip := addr.IP.String()
	if ip == "" || ip == "<nil>" {
		return ""
	}
	return fmt.Sprintf("%s", ip)
}
