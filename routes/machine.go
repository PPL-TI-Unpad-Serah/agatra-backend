package routes

import (
	"encoding/csv"
	"fmt"
	"imazine/models"
	"imazine/storage"
	"imazine/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
	"github.com/gin-gonic/gin"
)

func CreateMachine(context *fiber.Ctx) error{
	
}