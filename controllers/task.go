package controllers

import (
	"TESTOVOE/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateTaskInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func CreateTask(c *gin.Context) {
	var input CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{Title: input.Title, Content: input.Content}
	models.DB.Create(&task)

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func FindTasks(c *gin.Context) {
	var tasks []models.Task
	models.DB.Find(&tasks)

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func FindTask(c *gin.Context) {
	var task models.Task

	if err := models.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

type UpdateTaskInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func UpdateTask(c *gin.Context) {
	var Task models.Task
	if err := models.DB.Where("id = ?", c.Param("id")).First(&Task).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	var input UpdateTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTask := models.Task{Title: input.Title, Content: input.Content}

	models.DB.Model(&Task).Updates(&updatedTask)
	c.JSON(http.StatusOK, gin.H{"data": Task})
}

func DeleteTask(c *gin.Context) {
	var task models.Task
	if err := models.DB.Where("id = ?", c.Param("id")).First(&task).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "record not found"})
		return
	}

	models.DB.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"data": "success"})
}

