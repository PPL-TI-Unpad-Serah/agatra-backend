package main

import (
	"agatra/api"
	"agatra/db"
	// "agatra/model"
	"agatra/service"
	"agatra/middleware"
	"log"
	"os"

	"fmt"
	"sync"
	_ "embed"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

type APIHandler struct {
	MachineAPIHandler	api.MachineAPI
	UserAPIHandler     	api.UserAPI
	VersionAPIHandler	api.VersionAPI
	TitleAPIHandler		api.TitleAPI
	LocationAPIHandler	api.LocationAPI
	CityAPIHandler		api.CityAPI
	CenterAPIHandler	api.CenterAPI
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

		// conn.AutoMigrate(&model.Machine{}) //,&model.User{}, &model.Session{}, &model.Category{}, &model.Task{})

		router = RunServer(conn, router)

		fmt.Println("Server is running on port 8080")
		err = router.Run(":8080")
		if err != nil {
			panic(err)
		}

	}()

	wg.Wait()
}

func RunServer(db *gorm.DB, gin *gin.Engine) *gin.Engine {
	machineService := service.NewMachineService(db)
	userService := service.NewUserService(db)
	sessionService := service.NewSessionService(db)
	versionService := service.NewVersionService(db)
	titleService := service.NewTitleService(db)
	locationService := service.NewLocationService(db)
	cityService := service.NewCityService(db)
	centerService := service.NewCenterService(db)

	machineAPIHandler := api.NewMachineAPI(machineService)
	userAPIHandler := api.NewUserAPI(userService, sessionService)
	versionAPIHandler := api.NewVersionAPI(versionService)
	titleAPIHandler := api.NewTitleAPI(titleService)
	locationAPIHandler := api.NewLocationAPI(locationService)
	cityAPIHandler := api.NewCityAPI(cityService)
	centerAPIHandler := api.NewCenterAPI(centerService)

	apiHandler := APIHandler{
		MachineAPIHandler: 	machineAPIHandler,
		UserAPIHandler:     userAPIHandler,
		VersionAPIHandler: 	versionAPIHandler,
		TitleAPIHandler: 	titleAPIHandler,
		LocationAPIHandler: locationAPIHandler,
		CityAPIHandler: 	cityAPIHandler,	
		CenterAPIHandler: 	centerAPIHandler,
	}

	alpha := gin.Group("/v0")
	{
		admin := alpha.Group("/admin")
		{
			users := admin.Group("/user")
			{
				users.Use(middleware.Auth())
				users.POST("/add", apiHandler.UserAPIHandler.AddUser)
				users.GET("/get/:id", apiHandler.UserAPIHandler.GetUserByID)
				users.PUT("/update/:id", apiHandler.UserAPIHandler.UpdateUser)
				users.DELETE("/delete/:id", apiHandler.UserAPIHandler.DeleteUser)
				users.GET("/list", apiHandler.UserAPIHandler.GetUserList)
			}
		}

		city := alpha.Group("/city")
		{
			city.Use(middleware.Auth())
			city.POST("/add", apiHandler.CityAPIHandler.AddCity)
			city.GET("/get/:id", apiHandler.CityAPIHandler.GetCityByID)
			city.PUT("/update/:id", apiHandler.CityAPIHandler.UpdateCity)
			city.DELETE("/delete/:id", apiHandler.CityAPIHandler.DeleteCity)
			city.GET("/list", apiHandler.CityAPIHandler.GetCityList)
		}

		center := alpha.Group("/center")
		{
			center.Use(middleware.Auth())
			center.POST("/add", apiHandler.CenterAPIHandler.AddCenter)
			center.GET("/get/:id", apiHandler.CenterAPIHandler.GetCenterByID)
			center.PUT("/update/:id", apiHandler.CenterAPIHandler.UpdateCenter)
			center.DELETE("/delete/:id", apiHandler.CenterAPIHandler.DeleteCenter)
			center.GET("/list", apiHandler.CenterAPIHandler.GetCenterList)
		}


		location := alpha.Group("/location")
		{
			location.Use(middleware.Auth())
			location.POST("/add", apiHandler.LocationAPIHandler.AddLocation)
			location.GET("/get/:id", apiHandler.LocationAPIHandler.GetLocationByID)
			location.PUT("/update/:id", apiHandler.LocationAPIHandler.UpdateLocation)
			location.DELETE("/delete/:id", apiHandler.LocationAPIHandler.DeleteLocation)
			location.GET("/list", apiHandler.LocationAPIHandler.GetLocationList)
		}

		machine := alpha.Group("/machine")
		{
			machine.Use(middleware.Auth())
			machine.POST("/add", apiHandler.MachineAPIHandler.AddMachine)
			machine.GET("/get/:id", apiHandler.MachineAPIHandler.GetMachineByID)
			machine.PUT("/update/:id", apiHandler.MachineAPIHandler.UpdateMachine)
			machine.DELETE("/delete/:id", apiHandler.MachineAPIHandler.DeleteMachine)
			machine.GET("/list", apiHandler.MachineAPIHandler.GetMachineList)
		}

		version := alpha.Group("/version")
		{
			version.Use(middleware.Auth())
			version.POST("/add", apiHandler.VersionAPIHandler.AddVersion)
			version.GET("/get/:id", apiHandler.VersionAPIHandler.GetVersionByID)
			version.PUT("/update/:id", apiHandler.VersionAPIHandler.UpdateVersion)
			version.DELETE("/delete/:id", apiHandler.VersionAPIHandler.DeleteVersion)
			version.GET("/list", apiHandler.VersionAPIHandler.GetVersionList)
		}

		title := alpha.Group("/title")
		{
			title.Use(middleware.Auth())
			title.POST("/add", apiHandler.TitleAPIHandler.AddTitle)
			title.GET("/get/:id", apiHandler.TitleAPIHandler.GetTitleByID)
			title.PUT("/update/:id", apiHandler.TitleAPIHandler.UpdateTitle)
			title.DELETE("/delete/:id", apiHandler.TitleAPIHandler.DeleteTitle)
			title.GET("/list", apiHandler.TitleAPIHandler.GetTitleList)
		}

		user := alpha.Group("/user")
		{
			user.POST("/login", apiHandler.UserAPIHandler.Login)
			user.POST("/register", apiHandler.UserAPIHandler.Register)

			user.Use(middleware.Auth())
			user.POST("/logout", apiHandler.UserAPIHandler.Logout)
// 			user.GET("/tasks", apiHandler.UserAPIHandler.GetUserTaskCategory)
		}

// 		task := alpha.Group("/task")
// 		{
// 			task.Use(middleware.Auth())
// 			task.POST("/add", apiHandler.TaskAPIHandler.AddTask)
// 			task.GET("/get/:id", apiHandler.TaskAPIHandler.GetTaskByID)
// 			task.PUT("/update/:id", apiHandler.TaskAPIHandler.UpdateTask)
// 			task.DELETE("/delete/:id", apiHandler.TaskAPIHandler.DeleteTask)
// 			task.GET("/list", apiHandler.TaskAPIHandler.GetTaskList)
// 			task.GET("/category/:id", apiHandler.TaskAPIHandler.GetTaskListByCategory)
// 		}

// 		category := alpha.Group("/category")
// 		{
// 			category.Use(middleware.Auth())
// 			category.POST("/add", apiHandler.CategoryAPIHandler.AddCategory)
// 			category.GET("/get/:id", apiHandler.CategoryAPIHandler.GetCategoryByID)
// 			category.PUT("/update/:id", apiHandler.CategoryAPIHandler.UpdateCategory)
// 			category.DELETE("/delete/:id", apiHandler.CategoryAPIHandler.DeleteCategory)
// 			category.GET("/list", apiHandler.CategoryAPIHandler.GetCategoryList)
// 		}
	}

	return gin
}