package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListBags(ctx *gin.Context) {
	var bags []models.Bag
	if r := db.CONN.Preload("Cuboids").Find(&bags); r.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	ctx.JSON(http.StatusOK, bags)
}

func GetBag(ctx *gin.Context) {
	bagID := ctx.Param("bagID")

	var bag models.Bag
	if r := db.CONN.Preload("Cuboids").First(&bag, bagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	ctx.JSON(http.StatusOK, &bag)
}

func CreateBag(ctx *gin.Context) {
	var bagInput struct {
		Title  string
		Volume uint
	}

	if err := ctx.BindJSON(&bagInput); err != nil {
		return
	}

	bag := models.Bag{
		Title:   bagInput.Title,
		Volume:  bagInput.Volume,
		Cuboids: []models.Cuboid{},
	}
	if r := db.CONN.Create(&bag); r.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(r.Error, &err); ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	ctx.JSON(http.StatusCreated, &bag)
}

func DeleteBag(ctx *gin.Context) {
	bagID := ctx.Param("bagID")

	var bag models.Bag
	if r := db.CONN.First(&bag, bagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	if r := db.CONN.Delete(&bag); r.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}
