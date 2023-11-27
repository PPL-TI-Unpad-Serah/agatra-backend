package main

import (
	api "agatra/API"
	service "agatra/Service"
	"agatra/db"
	"agatra/middleware"
	"agatra/model"
	"log"
	"os"

	_ "embed"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
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
				users.POST("", apiHandler.UserAPIHandler.AddUser)
				users.PUT("/:id", apiHandler.UserAPIHandler.UpdateUser)
				users.DELETE("/:id", apiHandler.UserAPIHandler.DeleteUser)
				users.GET("/privileged", apiHandler.UserAPIHandler.GetPrivileged)
				users.GET("", apiHandler.UserAPIHandler.GetUserList)
				users.GET("/:id", apiHandler.UserAPIHandler.GetUserByID)
			}

			city := admin.Group("/cities")
			{
				city.POST("", apiHandler.CityAPIHandler.AddCity)
				city.PUT("/:id", apiHandler.CityAPIHandler.UpdateCity)
				city.DELETE("/:id", apiHandler.CityAPIHandler.DeleteCity)
			}
      
			center := admin.Group("/arcade_centers")
			{
				center.POST("", apiHandler.CenterAPIHandler.AddCenter)
				center.PUT("/:id", apiHandler.CenterAPIHandler.UpdateCenter)
				center.DELETE("/:id", apiHandler.CenterAPIHandler.DeleteCenter)
			}

			location := admin.Group("/arcade_locations")
			{
				location.POST("", apiHandler.LocationAPIHandler.AddLocation)
				location.PUT("/:id", apiHandler.LocationAPIHandler.UpdateLocation)
				location.DELETE("/:id", apiHandler.LocationAPIHandler.DeleteLocation)
			}

			version := admin.Group("/game_title_versions")
			{
				version.POST("", apiHandler.VersionAPIHandler.AddVersion)
				version.PUT("/:id", apiHandler.VersionAPIHandler.UpdateVersion)
				version.DELETE("/:id", apiHandler.VersionAPIHandler.DeleteVersion)
			}

			title := admin.Group("/game_titles")
			{
				title.POST("", apiHandler.TitleAPIHandler.AddTitle)
				title.PUT("/:id", apiHandler.TitleAPIHandler.UpdateTitle)
				title.DELETE("/:id", apiHandler.TitleAPIHandler.DeleteTitle)
			}
		}

		maintainer := alpha.Group("/maintainer", middleware.AuthMaintainer(db))
		{
			machine := maintainer.Group("/arcade_machines")
			{
				machine.POST("", apiHandler.MachineAPIHandler.AddMachine)
				machine.PUT("/:id", apiHandler.MachineAPIHandler.UpdateMachine)
				machine.DELETE("/:id", apiHandler.MachineAPIHandler.DeleteMachine)
			}
			
			location := maintainer.Group("/arcade_locations")
			{
				location.POST("", apiHandler.LocationAPIHandler.AddLocation)
				location.PUT("/:id", apiHandler.LocationAPIHandler.UpdateLocation)
				location.DELETE("/:id", apiHandler.LocationAPIHandler.DeleteLocation)
			}
		}

		city := alpha.Group("/cities")
		{
			city.GET("/:id", apiHandler.CityAPIHandler.GetCityByID)
			city.GET("", apiHandler.CityAPIHandler.GetCityList)
		}

		center := alpha.Group("/arcade_centers")
		{
			center.GET("/:id", apiHandler.CenterAPIHandler.GetCenterByID)
			center.GET("", apiHandler.CenterAPIHandler.GetCenterList)
		}


		location := alpha.Group("/arcade_locations")
		{
			location.GET("/:id", apiHandler.LocationAPIHandler.GetLocationByID)
			location.GET("", apiHandler.LocationAPIHandler.GetLocationList)
			location.GET("/nearby", apiHandler.LocationAPIHandler.GetLocationNearby)
			location.GET("/search", apiHandler.LocationAPIHandler.SearchLocation)
		}

		machine := alpha.Group("/arcade_machines")
		{
			machine.GET("/:id", apiHandler.MachineAPIHandler.GetMachineByID)
			machine.GET("", apiHandler.MachineAPIHandler.GetMachineList)
		}

		version := alpha.Group("/game_title_versions")
		{
			version.GET("/:id", apiHandler.VersionAPIHandler.GetVersionByID)
			version.GET("", apiHandler.VersionAPIHandler.GetVersionList)
		}

		title := alpha.Group("/game_titles")
		{
			title.GET("/:id", apiHandler.TitleAPIHandler.GetTitleByID)
			title.GET("", apiHandler.TitleAPIHandler.GetTitleList)
		}
		alpha.POST("/login", apiHandler.UserAPIHandler.Login)
		alpha.POST("/register", apiHandler.UserAPIHandler.Register)
		alpha.Use(middleware.Auth())
		alpha.POST("/profile", apiHandler.UserAPIHandler.Profile)
		alpha.POST("/logout", apiHandler.UserAPIHandler.Logout)
	}

	return gin
}