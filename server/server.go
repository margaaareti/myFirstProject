package server

import (
	"Test_derictory/auth"
	"Test_derictory/auth/delivery/authhttp"
	"Test_derictory/auth/service"
	"Test_derictory/auth/storage/postgres"
	"Test_derictory/auth/storage/redst"
	"Test_derictory/mainpage"
	"Test_derictory/mainpage/delivery/mainhttp"
	"Test_derictory/server/repository"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server
	authUC     auth.UseCase
	homeUC     mainpage.HomePage
	redisDB    *redis.Client
}

func NewApp() *App {

	//redst
	mediator := make(chan *redis.Client)

	go func() {
		conn, err := repository.NewRedisCon()
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("Redis is up!")
		mediator <- conn
	}()

	rdb := <-mediator

	//psql
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	userRepo := postgres.NewAuthPostgres(db)
	tokenStg := redst.NewAuthRedis(rdb)
	//homeUseCase := mainservice.NewHomeUseCase(tokenStg)

	return &App{
		authUC: service.NewAuthUseCase(userRepo, tokenStg),
		//homeUC:  homeUseCase,
		redisDB: rdb,
	}

}

func (a *App) Run(port string) error {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	router.LoadHTMLGlob("../ui/html/**")

	// Set up http handlers
	// SignIn/SignUp endpoints
	authhttp.RegisterHTTPEndPoints(router, a.authUC)

	//API endpoints
	authMiddleware := authhttp.NewAuthMiddleware(a.authUC, a.redisDB)
	api := router.Group("/api", authMiddleware)

	mainhttp.RegisterHTTPEndPoints(api, a.homeUC, a.authUC)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		MaxHeaderBytes: 1 << 20, //1mb
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			logrus.Fatal("Failed to listen and serve: %+v", err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
