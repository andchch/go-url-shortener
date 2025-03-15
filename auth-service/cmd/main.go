package main

import (
	"auth-service/internal/repository"
	"auth-service/internal/server"
	"auth-service/internal/utils"
	"log"
)

func main() {
	// Инициализация PostgreSQL
	//dsn := "postgres://user:password@localhost:5432/auth_db?sslmode=disable"
	//db, err := utils.InitDB(dsn)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()

	// Инициализация SQLite
	db, err := repository.NewSQLiteRepository("auth.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Запуск миграций
	err = utils.RunMigrations(db, "auth-service/migrations")
	if err != nil {
		log.Fatal(err)
	}

	// Запуск gRPC-сервера
	srv := server.NewAuthServer(db)
	// ... (остальной код сервера)
}
