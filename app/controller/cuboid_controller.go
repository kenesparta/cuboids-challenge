package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"cuboid-challenge/app/service"
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
	var cuboidInput service.CuboidInput

	if err := ctx.BindJSON(&cuboidInput); err != nil {
		return
	}

	var (
		cuboidService = service.NewCuboidService(db.CONN, &cuboidInput)
		method        = ctx.Request.Method
		cuboidRes     *models.Cuboid
		err           error
		status        int
	)

	if http.MethodPost == method {
		cuboidRes, err = cuboidService.Create()
		status = http.StatusCreated
	}

	if http.MethodPut == method {
		cuboidRes, err = cuboidService.Update()
		status = http.StatusOK
	}

	if err == nil {
		ctx.JSON(status, &cuboidRes)

		return
	}

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, service.ErrNotFound):
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	case errors.Is(err, service.ErrBagDisabled):
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})
	case errors.Is(err, service.ErrInsufficientCapacity):
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
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
