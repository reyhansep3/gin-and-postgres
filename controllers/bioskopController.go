package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

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

// jawaban tugas hari 14
func GetBioskop(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var getBioskop []Bioskop

		sqlStatement := `SELECT * from bioskop`

		rows, err := db.Query(sqlStatement)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		for rows.Next() {
			var bioskop = Bioskop{}

			err = rows.Scan(&bioskop.ID, &bioskop.Nama, &bioskop.Lokasi, &bioskop.Rating)

			if err != nil {
				panic(err)
			}
			getBioskop = append(getBioskop, bioskop)
		}

		fmt.Println("data bioskop", getBioskop)
		ctx.JSON(http.StatusOK, gin.H{
			"bioskop": getBioskop,
		})
	}
}

func GetBioskopByID(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParams := ctx.Param("id")
		var bioskop Bioskop

		query := `SELECT * FROM bioskop WHERE id = $1`
		id, err := strconv.Atoi(idParams)
		if err != nil {
			panic(err)
		}

		err = db.QueryRow(query, id).Scan(&bioskop.ID, &bioskop.Nama, &bioskop.Lokasi, &bioskop.Rating)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "Bioskop tidak ditemukan"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"bioskop": bioskop})
	}
}

func UpdateBioskop(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParams := ctx.Param("id")
		id, err := strconv.Atoi(idParams)
		if err != nil {
			panic(err)
		}
		var bioskop Bioskop
		if err := ctx.ShouldBindJSON(&bioskop); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		sqlStatement := `
		UPDATE bioskop
		SET nama = $2, lokasi = $3, rating = $4
		WHERE id = $1;
		`
		hasil, err := db.Exec(sqlStatement, id, bioskop.Nama, bioskop.Lokasi, bioskop.Rating)
		if err != nil {
			panic(err)
		}

		rowsAffected, err := hasil.RowsAffected()
		if err != nil {
			panic(err)
		}

		if rowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"Error": "data not found"})
		}

		ctx.JSON(http.StatusOK, gin.H{"bioskop": bioskop})
	}
}

func DeleteBioskop(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		paramId := ctx.Param("id")
		var bioskop Bioskop
		if err := ctx.ShouldBindJSON(&bioskop); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		id, err := strconv.Atoi(paramId)
		if err != nil {
			panic(err)
		}
		sqlStatement := `DELETE FROM bioskop WHERE id = $1`
		hasil, err := db.Exec(sqlStatement, id)
		if err != nil {
			panic(err)
		}

		rowsAffected, err := hasil.RowsAffected()
		if err != nil {
			panic(err)
		}

		if rowsAffected == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Data not found",
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"data": "data has deleted"})
	}
}
