package main

import (
	"auth_service/config"
	pkgPostgres "auth_service/pkg/postgres"
	pkgRedis "auth_service/pkg/redis"
	"auth_service/storage/cache"
	postgres "auth_service/storage/postgres"
	"log"

	"auth_service/server"
	"auth_service/service"
)

func main() {
	cnf := config.NewConfig()
	cnf.Load()

	db, err := pkgPostgres.ConnectDB(cnf.Database)
	if err != nil {
		log.Fatal(err)
	}

	rClient := pkgRedis.ConnectDB(&cnf.Redis)
	authCache := cache.NewAuthCache(rClient)
	emailCacher := cache.NewEmailCache(rClient)
	tokenCacher := cache.NewTokenCache(rClient)

	user := postgres.NewUserManagementSQL(db, authCache)

	emailSenderService := service.NewEmailSender(cnf.EmailSender, emailCacher)

	authService := service.NewAuthService(user, emailSenderService, cnf, tokenCacher)

	if err := server.Run(authService, *cnf); err != nil {
		log.Fatal(err)
	}
}
