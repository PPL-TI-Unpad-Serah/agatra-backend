package main_test

import (
	"agatra/controller"
	"agatra/db"
	"agatra/model"
	"testing"
	"net/http"
	"net/http/httptest"

	_ "embed"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func start() *gin.Engine{
	gin.SetMode(gin.ReleaseMode)
	config := &db.Config{
		Host:     "https://kkmkegheitaqvhygnixh.supabase.co",
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
	router = controller.RunServer(conn, router)
	fmt.Println("Server is running on port 8080")
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}

	return router
}

func TestMain(t *testing.T){
	t.Run("ServerRunning", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/arcade_locations", nil)
		w := httptest.NewRecorder()
		router := start() 
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}