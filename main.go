package main

import (
	"agatra/api"
	"agatra/db"
	"agatra/model"
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

		conn.AutoMigrate(&model.City{}, &model.Center{}, &model.Session{}, &model.User{}, &model.Title{})
		conn.AutoMigrate(&model.Version{})
		conn.AutoMigrate(&model.Machine{}) 
		conn.AutoMigrate(&model.Location{})

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

	alpha := gin.Group("/agatra")
	{
		admin := alpha.Group("/admin", middleware.AuthAdmin(db))
		{
			users := admin.Group("/users")

			{
				users.POST("/add", apiHandler.UserAPIHandler.AddUser)
				users.PUT("/update/:id", apiHandler.UserAPIHandler.UpdateUser)
				users.DELETE("/delete/:id", apiHandler.UserAPIHandler.DeleteUser)
				users.GET("/privileged", apiHandler.UserAPIHandler.GetPrivileged)
				users.GET("/search", apiHandler.UserAPIHandler.SearchName)
			}

			city := admin.Group("/cities")
			{
				city.POST("/add", apiHandler.CityAPIHandler.AddCity)
				city.PUT("/update/:id", apiHandler.CityAPIHandler.UpdateCity)
				city.DELETE("/delete/:id", apiHandler.CityAPIHandler.DeleteCity)
			}
      
			center := admin.Group("/arcade_centers")
			{
				center.POST("/add", apiHandler.CenterAPIHandler.AddCenter)
				center.PUT("/update/:id", apiHandler.CenterAPIHandler.UpdateCenter)
				center.DELETE("/delete/:id", apiHandler.CenterAPIHandler.DeleteCenter)
			}

			location := admin.Group("/arcade_locations")
			{
				location.POST("/add", apiHandler.LocationAPIHandler.AddLocation)
				location.PUT("/update/:id", apiHandler.LocationAPIHandler.UpdateLocation)
				location.DELETE("/delete/:id", apiHandler.LocationAPIHandler.DeleteLocation)
			}

			machine := admin.Group("/arcade_machines")
			{
				machine.POST("/add", apiHandler.MachineAPIHandler.AddMachine)
				machine.PUT("/update/:id", apiHandler.MachineAPIHandler.UpdateMachine)
				machine.DELETE("/delete/:id", apiHandler.MachineAPIHandler.DeleteMachine)
			}

			version := admin.Group("/game_title_versions")
			{
				version.POST("/add", apiHandler.VersionAPIHandler.AddVersion)
				version.PUT("/update/:id", apiHandler.VersionAPIHandler.UpdateVersion)
				version.DELETE("/delete/:id", apiHandler.VersionAPIHandler.DeleteVersion)
			}

			title := admin.Group("/game_titles")
			{
				title.POST("/add", apiHandler.TitleAPIHandler.AddTitle)
				title.PUT("/update/:id", apiHandler.TitleAPIHandler.UpdateTitle)
				title.DELETE("/delete/:id", apiHandler.TitleAPIHandler.DeleteTitle)
			}
		}

		city := alpha.Group("/cities")
		{
			city.GET("/get/:id", apiHandler.CityAPIHandler.GetCityByID)
			city.GET("/list", apiHandler.CityAPIHandler.GetCityList)
		}

		center := alpha.Group("/arcade_centers")
		{
			center.GET("/get/:id", apiHandler.CenterAPIHandler.GetCenterByID)
			center.GET("/list", apiHandler.CenterAPIHandler.GetCenterList)
		}


		location := alpha.Group("/arcade_locations")
		{
			location.GET("/get/:id", apiHandler.LocationAPIHandler.GetLocationByID)
			location.GET("/list", apiHandler.LocationAPIHandler.GetLocationList)
			location.GET("/nearby/:lat/:long", apiHandler.LocationAPIHandler.GetLocationNearby)
		}

		machine := alpha.Group("/arcade_machines")
		{
			machine.GET("/get/:id", apiHandler.MachineAPIHandler.GetMachineByID)
			machine.GET("/list", apiHandler.MachineAPIHandler.GetMachineList)
		}

		version := alpha.Group("/game_title_versions")
		{
			version.GET("/get/:id", apiHandler.VersionAPIHandler.GetVersionByID)
			version.GET("/list", apiHandler.VersionAPIHandler.GetVersionList)
		}

		title := alpha.Group("/game_titles")
		{
			title.GET("/get/:id", apiHandler.TitleAPIHandler.GetTitleByID)
			title.GET("/list", apiHandler.TitleAPIHandler.GetTitleList)
		}
		alpha.POST("/login", apiHandler.UserAPIHandler.Login)
		alpha.POST("/register", apiHandler.UserAPIHandler.Register)
		alpha.Use(middleware.Auth())
		alpha.POST("/logout", apiHandler.UserAPIHandler.Logout)
	}

	return gin
}