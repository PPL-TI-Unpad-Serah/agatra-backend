package controller

import(
	api "agatra/API"
	service "agatra/Service"
	"agatra/middleware"

	"github.com/gin-gonic/gin"
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
		alpha.GET("/profile", apiHandler.UserAPIHandler.Profile)
		alpha.POST("/logout", apiHandler.UserAPIHandler.Logout)
	}

	return gin
}