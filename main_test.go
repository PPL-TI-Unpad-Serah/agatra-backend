package main

import (
	"agatra/db"
	"agatra/model"

	_ "embed"
	"fmt"
	"sync"
	"testing"
	"net/http"
	"net/http/httptest"
	"time"
	"strings"
	// "io"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func start() (*gin.Engine, *sync.WaitGroup){
	gin.SetMode(gin.ReleaseMode)
	config := &db.Config{
		Host:     "db.kkmkegheitaqvhygnixh.supabase.co",
		Port:     "5432",
		Password: "kalomaugenerikaja",
		User:     "postgres",
		SSLMode:  "disable",
		DBName:   "postgres",
	}
	router := gin.New()
	db := db.NewDB()
	conn, err := db.Connect(config)
	if err != nil {panic(err)}
	conn.AutoMigrate(&model.City{}, &model.Center{}, &model.Session{}, &model.User{}, &model.Title{})
	conn.AutoMigrate(&model.Version{})
	conn.AutoMigrate(&model.Machine{}) 
	conn.AutoMigrate(&model.Location{})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		router = RunServer(conn, router)
		fmt.Println("Server is running on port 8080")
		err := router.Run(":8080")
		if err != nil {
			panic(err)
		}
	}()

	return router, &wg
}

func stop(router *gin.Engine, wg *sync.WaitGroup){
	wg.Done()
}

func TestMain(t *testing.T){
	t.Run("Register", func(t *testing.T) {
		userBody := `{
			"username" : "Rommel22w",
			"email" : "rommela.malik@gmail.com",
			"password" : "abcd",
			"confirm_password" : "abcd"
		 }`
		req, _ := http.NewRequest("POST", "http://localhost:8080/agatra/register", strings.NewReader(userBody))
		w := httptest.NewRecorder()
		router, wg := start() 
		defer stop(router, wg)
		

		time.Sleep(10 * time.Second)
		fmt.Println("Registered Routes:", router.Routes())
		router.ServeHTTP(w, req)
		// Send the request to the server
		// res, err := http.DefaultClient.Do(req)
		// if err != nil {
		// 	t.Fatalf("failed to send request: %v", err)
		// }
		// defer res.Body.Close()

		// Print the response body for debugging
		// body, _ := io.ReadAll(res.Body)
		// fmt.Println("Response Body:", string(body))
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}