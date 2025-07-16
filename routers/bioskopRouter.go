package routers

import (
	"database/sql"
	controller "gin-and-postgres/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer(db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.POST("/bioskop", controller.CreateBioskop(db))
	return router
}
