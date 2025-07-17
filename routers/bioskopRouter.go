package routers

import (
	"database/sql"
	controller "gin-and-postgres/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer(db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.POST("/bioskop", controller.CreateBioskop(db))

	// jawaban tugas hari 14
	router.GET("/bioskop", controller.GetBioskop(db))
	router.GET("/bioskop/:id", controller.GetBioskopByID(db))
	router.PUT("/bioskop/:id", controller.UpdateBioskop(db))
	router.DELETE("/bioskop/:id", controller.DeleteBioskop(db))

	return router
}
