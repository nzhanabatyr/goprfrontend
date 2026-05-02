package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"goprfrontend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

var jwtKey = []byte("my_secret_key")

func parseID(id string) uint {
	i, _ := strconv.Atoi(id)
	return uint(i)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := uint(claims["user_id"].(float64))

			c.Set("user_id", userID)
		} else {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
func main() {
	dsn := "host=localhost user=postgres password=nurlan050 dbname=task_management_db port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("✅ Connected to PostgreSQL!")

	_ = db

	db.AutoMigrate(&models.Task{}, &models.User{}, &models.FavoriteBook{}, &models.Book{})

	r := gin.Default()

	r.POST("/tasks", AuthMiddleware(), func(c *gin.Context) {
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

	r.POST("/register", func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		user.Password = string(hashedPassword)

		db.Create(&user)

		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
		})
	})

	r.POST("/login", func(c *gin.Context) {
		var input models.User
		var user models.User

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		db.Where("email = ?", input.Email).First(&user)

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, _ := token.SignedString(jwtKey)

		c.JSON(200, gin.H{
			"token": tokenString,
		})
	})

	r.PUT("/books/:id/favorites", AuthMiddleware(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		bookID := c.Param("id")

		fav := models.FavoriteBook{
			UserID: userID.(uint),
		}

		db.First(&models.Book{}, bookID) // просто проверка что книга есть

		fav.BookID = uint(parseID(bookID))

		db.Create(&fav)

		c.JSON(200, gin.H{"message": "Added to favorites"})
	})

	r.DELETE("/books/:id/favorites", AuthMiddleware(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		bookID := c.Param("id")

		db.Where("user_id = ? AND book_id = ?", userID, bookID).
			Delete(&models.FavoriteBook{})

		c.JSON(200, gin.H{"message": "Removed from favorites"})
	})

	r.GET("/books/favorites", AuthMiddleware(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		var books []models.Book

		db.Table("books").
			Joins("JOIN favorite_books ON favorite_books.book_id = books.id").
			Where("favorite_books.user_id = ?", userID).
			Find(&books)

		c.JSON(200, books)
	})

	r.Run(":8082")

}
