package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListCuboids(ctx *gin.Context) {
	var cuboids []models.Cuboid
	if r := db.CONN.Find(&cuboids); r.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	ctx.JSON(http.StatusOK, cuboids)
}

func GetCuboid(ctx *gin.Context) {
	var (
		cuboidID = ctx.Param("cuboidID")
		cuboid   models.Cuboid
	)

	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	ctx.JSON(http.StatusOK, &cuboid)
}

func CreateCuboid(ctx *gin.Context) {
	var cuboidInput struct {
		ID     uint `json:"id"`
		Width  uint
		Height uint
		Depth  uint
		BagID  uint `json:"bagId"`
	}

	if err := ctx.BindJSON(&cuboidInput); err != nil {
		return
	}

	cuboid := models.Cuboid{
		Width:  cuboidInput.Width,
		Height: cuboidInput.Height,
		Depth:  cuboidInput.Depth,
		BagID:  cuboidInput.BagID,
	}

	var bag models.Bag
	if r := db.CONN.Preload("Cuboids").First(&bag, cuboid.BagID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	if !bag.HasCapacity(cuboid.PayloadVolume()) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})

		return
	}

	if bag.Disabled {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})

		return
	}

	var (
		queryRes *gorm.DB
		status   int
		method   = ctx.Request.Method
	)

	switch method {
	case http.MethodPost:
		queryRes = db.CONN.Create(&cuboid)
		status = http.StatusCreated
	case http.MethodPut:
		cuboid.ID = cuboidInput.ID
		queryRes = db.CONN.Updates(&cuboid)

		if queryRes.RowsAffected == 0 {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})

			return
		}

		status = http.StatusOK
	default:
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Not Allowed"})

		return
	}

	if queryRes.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(queryRes.Error, &err); ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": queryRes.Error.Error()})
		}

		return
	}

	ctx.JSON(status, &cuboid)
}

func DeleteCuboid(ctx *gin.Context) {
	var (
		cuboidID = ctx.Param("cuboidID")
		cuboid   models.Cuboid
	)

	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	if r := db.CONN.Delete(&cuboid); r.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": cuboid.ID})
}
