package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Bioskop struct {
	ID     string  `json:"id"`
	Nama   string  `json:"nama"`
	Lokasi string  `json:"lokasi"`
	Rating float64 `json:"rating"`
}

func CreateBioskop(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newBioskop Bioskop

		if err := ctx.ShouldBindJSON(&newBioskop); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if newBioskop.Nama == "" || newBioskop.Lokasi == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan Lokasi tidak boleh kosong"})
			return
		}

		sqlStatement := `
		INSERT INTO bioskop (nama, lokasi, rating)
		VALUES ($1, $2, $3)
		Returning id
		`

		err := db.QueryRow(sqlStatement, newBioskop.Nama, newBioskop.Lokasi, newBioskop.Rating).
			Scan(&newBioskop.ID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"bioskop": newBioskop,
		})
	}
}
