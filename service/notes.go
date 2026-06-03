package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/models"
	"gorm.io/gorm"
)

func GetNoteByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var note models.Note

	result := db.First(&note, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Заметка не найдена",
		})
		return
	}

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusOK, note)
}

func GetNotesByStudentID(ctx *gin.Context) {
	id := ctx.Param("id")

	var student models.Student

	result := db.First(&student, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Студент не найден",
		})
		return
	}

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	var notes []models.Note

	result = db.Where("student_id = ?", id).Find(&notes)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusOK, notes)
}

func CreateNote(ctx *gin.Context) {
	var request models.CreateNoteRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректное тело запроса",
		})
		return
	}

	if strings.TrimSpace(request.Text) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Текст заметки обязателен",
		})
		return
	}

	var student models.Student

	result := db.First(&student, request.StudentID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Студент не существует",
		})
		return
	}

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	note := models.Note{
		StudentID: request.StudentID,
		Author:    request.Author,
		Text:      request.Text,
	}

	result = db.Create(&note)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusCreated, note)
}

func UpdateNote(ctx *gin.Context) {
	id := ctx.Param("id")

	var note models.Note

	result := db.First(&note, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Заметка не найдена",
		})
		return
	}

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	var request models.UpdateNoteRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректное тело запроса",
		})
		return
	}

	if request.StudentID != nil {
		var student models.Student

		result = db.First(&student, *request.StudentID)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Студент не существует",
			})
			return
		}

		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка базы данных",
			})
			return
		}

		note.StudentID = *request.StudentID
	}

	if request.Author != nil {
		note.Author = *request.Author
	}

	if request.Text != nil {
		if strings.TrimSpace(*request.Text) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Текст заметки не может быть пустым",
			})
			return
		}

		note.Text = *request.Text
	}

	result = db.Save(&note)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusOK, note)
}

func DeleteNote(ctx *gin.Context) {
	id := ctx.Param("id")

	result := db.Delete(&models.Note{}, id)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Заметка не найдена",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Заметка удалена",
	})
}
