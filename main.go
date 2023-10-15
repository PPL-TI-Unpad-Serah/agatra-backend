package main

import (
	"agatra/db"
	// "agatra/routes"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func setRoutes(app *fiber.App){
	// app.Post("/login", routes.Login)
	// app.Post("/register", routes.Register)
}

func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Missing .env file. Probably okay on dockerized environment")
	}
	config := &db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db.Connect(config)
	app := fiber.New()

	setRoutes(app)

	app.Listen(":8080")
}