package main

import (
	"agatra/db"
	"agatra/model"
	"bytes"

	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

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

func getSuccess(w *bytes.Buffer){
	var temp model.SuccessResponse
	body, err := ioutil.ReadAll(w)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}
	resJson := []byte(body)
	err = json.Unmarshal(resJson, &temp)
	if err != nil {
		var temp2 model.ErrorResponse
		err = json.Unmarshal(resJson, &temp2)
		if err != nil{
			fmt.Println("Error parsing JSON:", err)
		}
		fmt.Println(temp2.Error)
	}else{
		fmt.Println(temp.Message)
	}
}

func stop(router *gin.Engine, wg *sync.WaitGroup){
	wg.Done()
}

func TestMain(t *testing.T){
	t.Run("Running Server", func(t *testing.T){
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
		t.Run("Cities", func(t *testing.T){
			t.Run("Get Cities", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/cities", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				fmt.Println("Get Cities Body:", w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Get Cities By ID", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/cities/1", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				fmt.Println("Get Cities Body By ID:", w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Add Cities", func(t *testing.T){
				cityBody := `{
					"Name" : "Bekasih"
				}`
				req, _ := http.NewRequest("POST", "/agatra/admin/cities", strings.NewReader(cityBody))
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Update Cities", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/cities", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var temp2 model.CityArrayResponse
				body, err := ioutil.ReadAll(w.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
				}
				resJson := []byte(body)
				err = json.Unmarshal(resJson, &temp2)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}

				cityBody := `{
					"Name" : "Bekasi"
				}`
				id := strconv.Itoa(temp2.Cities[len(temp2.Cities)-1].ID)
				req, _ = http.NewRequest("PUT", "/agatra/admin/cities/" + id, strings.NewReader(cityBody))
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Delete Cities", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/cities", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var temp2 model.CityArrayResponse
				body, err := ioutil.ReadAll(w.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
				}
				resJson := []byte(body)
				err = json.Unmarshal(resJson, &temp2)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}
				id := strconv.Itoa(temp2.Cities[len(temp2.Cities)-1].ID)
				req, _ = http.NewRequest("DELETE", "/agatra/admin/cities/" + id, nil)
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
		})
		t.Run("Arcade Locations", func(t *testing.T){
			t.Run("Get Arcade Locations", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/arcade_locations", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				fmt.Println("Get Arcade Locations Body:", w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Get Arcade Locations By ID", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/arcade_locations/1", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				fmt.Println("Get Arcade Locations Body By ID:", w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Add Arcade Locations", func(t *testing.T){
				locationBody := `{
					"Name" : "2",
					"Description" : "2",
					"Lat" : 2,
					"Long" : 2,
					"Center_id" : 1,
					"City_id" : 1
				}`
				req, _ := http.NewRequest("POST", "/agatra/admin/arcade_locations", strings.NewReader(locationBody))
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Update Arcade Locations", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/arcade_locations", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var temp2 model.LocationArrayResponse
				body, err := ioutil.ReadAll(w.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
				}
				resJson := []byte(body)
				err = json.Unmarshal(resJson, &temp2)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}

				locationBody := `{
					"Name" : "1",
					"Description" : "1",
					"Lat" : 1,
					"Long" : 1,
					"Center_id" : 1,
					"City_id" : 1
				}`
				id := strconv.Itoa(temp2.Locations[len(temp2.Locations)-1].ID)
				req, _ = http.NewRequest("PUT", "/agatra/admin/arcade_locations/" + id, strings.NewReader(locationBody))
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Delete Arcade Locations", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/arcade_locations", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var temp2 model.LocationArrayResponse
				body, err := ioutil.ReadAll(w.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
				}
				resJson := []byte(body)
				err = json.Unmarshal(resJson, &temp2)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}
				id := strconv.Itoa(temp2.Locations[len(temp2.Locations)-1].ID)
				req, _ = http.NewRequest("DELETE", "/agatra/admin/arcade_locations/" + id, nil)
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
		})
		t.Run("Arcade Centers", func(t *testing.T){
			t.Run("Get Arcade Centers", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/arcade_centers", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				fmt.Println("Get Arcade Centers Body:", w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Get Arcade Centers By ID", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/arcade_centers/1", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				fmt.Println("Get Arcade Centers Body By ID:", w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Add Arcade Centers", func(t *testing.T){
				centerBody := `{
					"Name" : "Tz",
					"Info" : "test"
				}`
				req, _ := http.NewRequest("POST", "/agatra/admin/arcade_centers", strings.NewReader(centerBody))
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Update Arcade Centers", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/arcade_centers", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var temp2 model.CenterArrayResponse
				body, err := ioutil.ReadAll(w.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
				}
				resJson := []byte(body)
				err = json.Unmarshal(resJson, &temp2)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}

				centerBody := `{
					"Name" : "Timez",
					"Info" : "test"
				}`
				id := strconv.Itoa(temp2.Centers[len(temp2.Centers)-1].ID)
				req, _ = http.NewRequest("PUT", "/agatra/admin/arcade_centers/" + id, strings.NewReader(centerBody))
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Delete Arcade Centers", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/arcade_centers", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var temp2 model.CenterArrayResponse
				body, err := ioutil.ReadAll(w.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
				}
				resJson := []byte(body)
				err = json.Unmarshal(resJson, &temp2)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}
				id := strconv.Itoa(temp2.Centers[len(temp2.Centers)-1].ID)
				req, _ = http.NewRequest("DELETE", "/agatra/admin/arcade_centers/" + id, nil)
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
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
		t.Run("Game Titles", func(t *testing.T){
			t.Run("Get Game Titles", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/game_titles", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				fmt.Println("Get Game Titles Body:", w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Get Game Titles By ID", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/game_titles/1", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				fmt.Println("Get Game Titles Body By ID:", w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Add Game Titles", func(t *testing.T){
				titleBody := `{
					"Name" : "maymay"
				}`
				req, _ := http.NewRequest("POST", "/agatra/admin/game_titles", strings.NewReader(titleBody))
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Update Game Titles", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/game_titles", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var temp2 model.TitleArrayResponse
				body, err := ioutil.ReadAll(w.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
				}
				resJson := []byte(body)
				err = json.Unmarshal(resJson, &temp2)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}

				titleBody := `{
					"Name" : "MaiMai"
				}`
				id := strconv.Itoa(temp2.Titles[len(temp2.Titles)-1].ID)
				req, _ = http.NewRequest("PUT", "/agatra/admin/game_titles/" + id, strings.NewReader(titleBody))
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
			t.Run("Delete Game Titles", func(t *testing.T){
				req, _ := http.NewRequest("GET", "/agatra/game_titles", nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				var temp2 model.TitleArrayResponse
				body, err := ioutil.ReadAll(w.Body)
				if err != nil {
					fmt.Println("Error reading response:", err)
				}
				resJson := []byte(body)
				err = json.Unmarshal(resJson, &temp2)
				if err != nil {
					fmt.Println("Error parsing JSON:", err)
				}
				id := strconv.Itoa(temp2.Titles[len(temp2.Titles)-1].ID)
				req, _ = http.NewRequest("DELETE", "/agatra/admin/game_titles/" + id, nil)
				req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
				router.ServeHTTP(w, req)
				getSuccess(w.Body)
				assert.Equal(t, http.StatusOK, w.Code)
			})
		})
		
		t.Run("Get Users", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/admin/users", nil)
			req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Println("Get Users Body:", w.Body)
			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Get Privileged Users", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/admin/users/privileged", nil)
			req.Header.Set("Authorization", "Bearer "+ temp.Data.ApiKey)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			fmt.Println("Get Privileged Body:", w.Body)
			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("Get Profile", func(t *testing.T){
			req, _ := http.NewRequest("GET", "/agatra/profile", nil)
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