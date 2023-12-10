package main

import (
	"agatra/db"
	"agatra/controller"
	"agatra/model"
	"log"
	"os"

	_ "embed"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main(){
	// gin.SetMode(gin.ReleaseMode)
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

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		router := gin.New()
		db := db.NewDB()
		router.Use(gin.Recovery())
		router.ForwardedByClientIP = true
		router.SetTrustedProxies([]string{"127.0.0.1"})

		conn, err := db.Connect(config)
		if err != nil {
			panic(err)
		}

		conn.AutoMigrate(&model.City{}, &model.Center{}, &model.Session{}, &model.User{}, &model.Title{})
		conn.AutoMigrate(&model.Version{})
		conn.AutoMigrate(&model.Machine{}) 
		conn.AutoMigrate(&model.Location{})

		router = controller.RunServer(conn, router)

		fmt.Println("Server is running on port 8080")
		err = router.Run(":8080")
		if err != nil {
			panic(err)
		}

	}()

	wg.Wait()
}