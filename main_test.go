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
	"encoding/json"
	"io/ioutil"

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
	t.Run("Get", func(t *testing.T){
		userBody := `{
			"Username" : "Rommel22w",
			"Password" : "abcd"
		 }`
		req, _ := http.NewRequest("POST", "/agatra/login", strings.NewReader(userBody))
		w := httptest.NewRecorder()
		router, wg := start() 
		defer stop(router, wg)

		time.Sleep(1 * time.Second)
		router.ServeHTTP(w, req)
		var temp model.LoginResponse
		body, err := ioutil.ReadAll(w.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
		}
		resJson := []byte(body)
		err = json.Unmarshal(resJson, &temp)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}
		assert.Equal(t, http.StatusOK, w.Code)
		t.Run("Get Cities", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/cities", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Println("Get Cities Body:", w.Body)
			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Get Arcade Locations", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/arcade_locations", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Println("Get Locations Body:", w.Body)
			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Get Arcade Centers", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/arcade_centers", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Println("Get Centers Body:", w.Body)
			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Get Arcade Machines", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/arcade_machines", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Println("Get Machines Body:", w.Body)
			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Get Game Title Versions", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/game_title_versions", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Println("Get Versions Body:", w.Body)
			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Get Game Titles", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/game_titles", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Println("Get Titles Body:", w.Body)
			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Get Profile", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/profile", nil)
			fmt.Println("APIKEY:", temp.Data.ApiKey)
			req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Println("Get Profiles Body:", w.Body)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	
	// t.Run("Register", func(t *testing.T) {
	// 	userBody := `{
	// 		"username" : "Rommel22w",
	// 		"email" : "rommela.malik@gmail.com",
	// 		"password" : "abcd",
	// 		"confirm_password" : "abcd"
	// 	 }`
	// 	req, _ := http.NewRequest("POST", "/agatra/register", strings.NewReader(userBody))
	// 	w := httptest.NewRecorder()
	// 	router, wg := start() 
	// 	defer stop(router, wg)
		

	// 	time.Sleep(10 * time.Second)
	// 	fmt.Println("Registered Routes:", router.Routes())
	// 	router.ServeHTTP(w, req)
	// 	assert.Equal(t, http.StatusCreated, w.Code)
	// })
}