package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goprfrontend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=nurlan050 dbname=task_management_db port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("✅ Connected to PostgreSQL!")

	_ = db

	db.AutoMigrate(&models.Task{})

	r := gin.Default()

	r.POST("/tasks", func(c *gin.Context) {
		var task models.Task

		if err := c.BindJSON(&task); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		db.Create(&task)

		c.JSON(200, task)

	})

	r.GET("/tasks", func(c *gin.Context) {
		var tasks []models.Task

		db.Find(&tasks)

		c.JSON(200, tasks)
	})

	r.GET("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")

		var task models.Task

		if err := db.First(&task, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Task not found"})
			return
		}

		c.JSON(200, task)
	})

	r.PUT("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")

		var task models.Task

		if err := db.First(&task, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Task not found"})
			return
		}

		if err := c.BindJSON(&task); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		db.Save(&task)

		c.JSON(200, task)
	})

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")

		var task models.Task

		if err := db.First(&task, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Task not found"})
			return
		}

		db.Delete(&task)

		c.JSON(200, gin.H{"message": "Task deleted"})
	})

	r.GET("/tasks/status/:status", func(c *gin.Context) {
		status := c.Param("status")

		var tasks []models.Task

		db.Where("status = ?", status).Find(&tasks)

		c.JSON(200, tasks)
	})

	r.GET("/tasks/search/:title", func(c *gin.Context) {
		title := c.Param("title")

		var tasks []models.Task

		db.Where("title ILIKE ?", "%"+title+"%").Find(&tasks)

		c.JSON(200, tasks)
	})

	r.GET("/tasks/count", func(c *gin.Context) {
		var count int64

		db.Model(&models.Task{}).Count(&count)

		c.JSON(200, gin.H{
			"count": count,
		})
	})

	r.GET("/tasks/done", func(c *gin.Context) {
		var tasks []models.Task

		db.Where("status = ?", "done").Find(&tasks)

		c.JSON(200, tasks)
	})

	r.DELETE("/tasks/done", func(c *gin.Context) {
		db.Where("status = ?", "done").Delete(&models.Task{})

		c.JSON(200, gin.H{
			"message": "All done tasks deleted",
		})
	})

	r.Run(":8082")

}
