package service

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/models"
	"gorm.io/gorm"
)

func GetAllStudents(ctx *gin.Context) {
	var students []models.Student

	query := db.Model(&models.Student{})

	groupID := ctx.Query("group_id")
	if groupID != "" {
		if _, err := strconv.Atoi(groupID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Некорректный group_id",
			})
			return
		}

		query = query.Where("group_id = ?", groupID)
	}

	paymentStatus := ctx.Query("payment_status")
	if paymentStatus != "" {
		validPaymentStatuses := map[string]bool{
			string(models.PaymentPaid):    true,
			string(models.PaymentUnpaid):  true,
			string(models.PaymentPartial): true,
		}

		if !validPaymentStatuses[paymentStatus] {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Некорректный payment_status",
			})
			return
		}

		query = query.Where("payment_status = ?", paymentStatus)
	}

	studyStatus := ctx.Query("study_status")
	if studyStatus != "" {
		validStudyStatuses := map[string]bool{
			string(models.StudyLearning):  true,
			string(models.StudyJobSearch): true,
			string(models.StudyOffer):     true,
			string(models.StudyWorking):   true,
		}

		if !validStudyStatuses[studyStatus] {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Некорректный study_status",
			})
			return
		}

		query = query.Where("study_status = ?", studyStatus)
	}

	result := query.Find(&students)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusOK, students)
}

func GetStudentByID(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, student)
}

func CreateStudent(ctx *gin.Context) {
	var request models.CreateStudentRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректное тело запроса",
		})
		return
	}

	if strings.TrimSpace(request.FullName) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "full_name обязателен",
		})
		return
	}

	if !isValidPaymentStatus(request.PaymentStatus) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный payment_status",
		})
		return
	}

	if !isValidStudyStatus(request.StudyStatus) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректный study_status",
		})
		return
	}

	var group models.Group

	result := db.First(&group, request.GroupID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Группа не существует",
		})
		return
	}

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	student := models.Student{
		FullName:      request.FullName,
		Email:         request.Email,
		Telegram:      request.Telegram,
		GroupID:       request.GroupID,
		TuitionTotal:  request.TuitionTotal,
		TuitionPaid:   request.TuitionPaid,
		PaymentStatus: request.PaymentStatus,
		StudyStatus:   request.StudyStatus,
	}

	result = db.Create(&student)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusCreated, student)
}

func UpdateStudent(ctx *gin.Context) {
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

	var request models.UpdateStudentRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Некорректное тело запроса",
		})
		return
	}

	if request.FullName != nil {
		if strings.TrimSpace(*request.FullName) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "full_name не может быть пустым",
			})
			return
		}

		student.FullName = *request.FullName
	}

	if request.Email != nil {
		student.Email = *request.Email
	}

	if request.Telegram != nil {
		student.Telegram = *request.Telegram
	}

	if request.GroupID != nil {
		var group models.Group

		result = db.First(&group, *request.GroupID)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Группа не существует",
			})
			return
		}

		if result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка базы данных",
			})
			return
		}

		student.GroupID = *request.GroupID
	}

	if request.TuitionTotal != nil {
		student.TuitionTotal = *request.TuitionTotal
	}

	if request.TuitionPaid != nil {
		student.TuitionPaid = *request.TuitionPaid
	}

	if request.PaymentStatus != nil {
		if !isValidPaymentStatus(*request.PaymentStatus) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Некорректный payment_status",
			})
			return
		}

		student.PaymentStatus = *request.PaymentStatus
	}

	if request.StudyStatus != nil {
		if !isValidStudyStatus(*request.StudyStatus) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Некорректный study_status",
			})
			return
		}

		student.StudyStatus = *request.StudyStatus
	}

	result = db.Save(&student)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusOK, student)
}

func DeleteStudent(ctx *gin.Context) {
	id := ctx.Param("id")

	result := db.Delete(&models.Student{}, id)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Студент не найден",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Студент удалён",
	})
}

func isValidPaymentStatus(status models.PaymentStatus) bool {
	switch status {
	case models.PaymentPaid,
		models.PaymentUnpaid,
		models.PaymentPartial:
		return true
	default:
		return false
	}
}

func isValidStudyStatus(status models.StudyStatus) bool {
	switch status {
	case models.StudyLearning,
		models.StudyJobSearch,
		models.StudyOffer,
		models.StudyWorking:
		return true
	default:
		return false
	}
}

func GetGroupStudents(ctx *gin.Context) {
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

	var students []models.Student

	result = db.Where("group_id = ?", id).Find(&students)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка базы данных",
		})
		return
	}

	ctx.JSON(http.StatusOK, students)
}