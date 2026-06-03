package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/models"
	"gorm.io/gorm"
)

func GetAllGroups(ctx *gin.Context) {
	var groups []models.Group

	query := db.Model(&models.Group{})

	week := ctx.Query("week")
	if week != "" {
		weekInt, err := strconv.Atoi(week)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Некорректный week",
			})
			return
		}

		query = query.Where("current_week = ?", weekInt)
	}

	finished := ctx.Query("finished")
	if finished == "true" {
		query = query.Where("is_finished = ?", true)
	}

	result := query.Find(&groups)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusOK, groups)
}

func GetGroupByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var group models.Group

	result := db.First(&group, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Группа не найдена",
		})
		return
	}

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusOK, group)
}

func CreateGroup(ctx *gin.Context) {
	var request models.CreateGroupRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректное тело запроса",
		})
		return
	}

	group := models.Group{
		Title:       request.Title,
		CurrentWeek: request.CurrentWeek,
		TotalWeeks:  request.TotalWeeks,
		IsFinished:  request.IsFinished,
	}

	result := db.Create(&group)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusCreated, group)
}

func UpdateGroup(ctx *gin.Context) {
	id := ctx.Param("id")

	var group models.Group

	result := db.First(&group, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Группа не найдена",
		})
		return
	}

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	var request models.UpdateGroupRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректное тело запроса",
		})
		return
	}

	if request.Title != nil {
		group.Title = *request.Title
	}

	if request.CurrentWeek != nil {
		group.CurrentWeek = *request.CurrentWeek
	}

	if request.TotalWeeks != nil {
		group.TotalWeeks = *request.TotalWeeks
	}

	if request.IsFinished != nil {
		group.IsFinished = *request.IsFinished
	}

	result = db.Save(&group)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusOK, group)
}

func DeleteGroup(ctx *gin.Context) {
	id := ctx.Param("id")

	result := db.Delete(&models.Group{}, id)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Группа не найдена",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Группа удалена",
	})
}

func GetOfferStats(ctx *gin.Context) {
	id := ctx.Param("id")

	var group models.Group

	result := db.First(&group, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Группа не найдена",
		})
		return
	}

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	var totalStudents int64

	result = db.Model(&models.Student{}).
		Where("group_id = ?", id).
		Count(&totalStudents)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	if totalStudents == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"group_id":      group.ID,
			"offer_percent": 0,
		})
		return
	}

	var offerStudents int64

	result = db.Model(&models.Student{}).
		Where("group_id = ? AND study_status = ?", id, models.StudyOffer).
		Count(&offerStudents)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	offerPercent := float64(offerStudents) / float64(totalStudents) * 100

	ctx.JSON(http.StatusOK, gin.H{
		"group_id":      group.ID,
		"offer_percent": offerPercent,
	})
}
