package main

import (
	"agatra/api"
	"agatra/db"
	"agatra/model"
	"agatra/service"
	"agatra/middleware"
	"log"
	"os"

	// "embed"
	"fmt"
	// "net/http"
	"sync"
	// "time"

	_ "embed"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	// "github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// func setRoutes(app *fiber.App){
// 		app.Post("/login", routes.Login)
// 		app.Post("/register", routes.Register)
// }

type APIHandler struct {
	MachineAPIHandler	api.MachineAPI
	UserAPIHandler     	api.UserAPI
	VersionAPIHandler	api.VersionAPI
// 	CategoryAPIHandler api.CategoryAPI
// 	TaskAPIHandler     api.TaskAPI
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

	machineAPIHandler := api.NewMachineAPI(machineService)
	userAPIHandler := api.NewUserAPI(userService, sessionService)
	versionAPIHandler := api.NewVersionAPI(versionService)

	apiHandler := APIHandler{
		MachineAPIHandler: 	machineAPIHandler,
		UserAPIHandler:     userAPIHandler,
		VersionAPIHandler: versionAPIHandler,
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