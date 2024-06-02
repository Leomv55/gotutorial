package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"todo/models"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.Todo{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	route := gin.Default()
	route.LoadHTMLGlob("templates/*")

	route.GET("/todos/", func(c *gin.Context) {
		var todos []models.Todo
		db.Find(&todos)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Todo list",
			"todos": todos,
		})
	})

	route.POST("/todos/", func(c *gin.Context) {
		var todo models.Todo

		if err := c.ShouldBind(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db.Create(&todo)
		c.Redirect(http.StatusSeeOther, "/todos/")
	})

	route.DELETE("/todos/:id/", func(c *gin.Context) {
		var todo models.Todo
		id := c.Param("id")
		err := db.First(&todo, id)
		if err.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}
		db.Delete(&todo)
		c.Redirect(http.StatusSeeOther, "/todos/")
	})

	api_route := route.Group("/api")

	api_route.PATCH("/todos/:id/", func(c *gin.Context) {
		var todo models.Todo
		id := c.Param("id")
		err := db.First(&todo, id)
		if err.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		if err := c.ShouldBindBodyWithJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		db.Save(&todo)
		c.JSON(http.StatusOK, todo)
	})

	if err := route.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
