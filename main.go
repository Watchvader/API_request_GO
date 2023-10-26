package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/users", getUsers)
	r.GET("/users/:id", getUser)
	r.POST("/users", createUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	return r
}

func getUsers(c *gin.Context) {
	var users []User
	db.Find(&users)

	c.JSON(200, gin.H{
		"data": users,
	})
}

func getUser(c *gin.Context) {
	var user User

	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "not found!"})
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})
}

func createUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Create(&user)

	c.JSON(200, gin.H{
		"data": user,
	})
}

func updateUser(c *gin.Context) {
	var user User

	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Save(&user)

	c.JSON(200, gin.H{
		"data": user,
	})
}

func deleteUser(c *gin.Context) {
	var user User

	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "not found"})
		return
	}

	db.Delete(&user)

	c.JSON(200, gin.H{
		"message": "deleted",
	})
}

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/prueba?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect")
	}

	db.AutoMigrate(&User{})

	r := setupRouter()
	r.Run(":8080")
}
