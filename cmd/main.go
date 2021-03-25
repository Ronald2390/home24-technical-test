package main

import (
	"log"
	"os"

	"home24-technical-test/config"
	"home24-technical-test/database"
	"home24-technical-test/database/seeder"
	internalhttp "home24-technical-test/internal/http"
	"home24-technical-test/internal/user"
	"home24-technical-test/internal/user/adapter"
	"home24-technical-test/internal/user/service"
	userStoragePostgres "home24-technical-test/internal/user/storage/postgres"
	userStorageRedis "home24-technical-test/internal/user/storage/redis"
	"home24-technical-test/pkg/data"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

func main() {
	// remove or can replaced with actual env var
	os.Setenv("ENVIRONMENT", "development")

	cfg, err := config.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}

	db, err := sqlx.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	defer db.Close()

	// Migrate the db
	database.MigrateUp(cfg)
	seeder.SeedUp(cfg.DBConnectionString, cfg.IsDevelopment)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	_, err = redisClient.Ping().Result()
	if err != nil {
		log.Fatalln(err)
	}

	userService := service.NewService(
		user.NewService(
			userStoragePostgres.NewPostgresStorage(db),
		),
		user.NewSessionService(
			userStorageRedis.NewSessionStorage(redisClient),
		),
	)

	getUserAdapter := adapter.NewGetUserAdapter(userService)
	getLoginSessionAdapter := adapter.NewGetLoginSessionAdapter(userService)
	loginAdapter := adapter.NewLoginAdapter(userService)
	logoutAdapter := adapter.NewLogoutAdapter(userService)
	changePasswordAdapter := adapter.NewChangePasswordAdapter(userService)

	dataManager := data.NewManager(db)

	s := internalhttp.NewServer(
		getUserAdapter,
		getLoginSessionAdapter,
		loginAdapter,
		logoutAdapter,
		changePasswordAdapter,
		dataManager,
	)
	s.Serve(cfg.ExposingPort)
}
