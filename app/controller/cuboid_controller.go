package controller

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListCuboids(c *gin.Context) {
	var cuboids []models.Cuboid
	if r := db.CONN.Find(&cuboids); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	c.JSON(http.StatusOK, cuboids)
}

func GetCuboid(c *gin.Context) {
	var (
		cuboidID = c.Param("cuboidID")
		cuboid   models.Cuboid
	)
	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}
	c.JSON(http.StatusOK, &cuboid)
}

func CreateCuboid(c *gin.Context) {
	var cuboidInput struct {
		ID     uint `json:"id"`
		Width  uint
		Height uint
		Depth  uint
		BagID  uint `json:"bagId"`
	}

	if err := c.BindJSON(&cuboidInput); err != nil {
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
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}
		return
	}

	if !bag.HasCapacity(cuboid.PayloadVolume()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Insufficient capacity in bag"})
		return
	}

	if bag.Disabled {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bag is disabled"})
		return
	}

	var (
		queryRes *gorm.DB
		status   int
		method   = c.Request.Method
	)
	switch method {
	case http.MethodPost:
		queryRes = db.CONN.Create(&cuboid)
		status = http.StatusCreated
	case http.MethodPut:
		cuboid.ID = cuboidInput.ID
		queryRes = db.CONN.Updates(&cuboid)
		if queryRes.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		status = http.StatusOK
	default:
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Not Allowed"})
		return
	}

	if queryRes.Error != nil {
		var err models.ValidationErrors
		if ok := errors.As(queryRes.Error, &err); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": queryRes.Error.Error()})
		}

		return
	}

	c.JSON(status, &cuboid)
}

func DeleteCuboid(c *gin.Context) {
	var (
		cuboidID = c.Param("cuboidID")
		cuboid   models.Cuboid
	)
	if r := db.CONN.First(&cuboid, cuboidID); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		}

		return
	}

	if r := db.CONN.Delete(&cuboid); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"id": cuboid.ID})
}
