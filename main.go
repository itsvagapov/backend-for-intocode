package main

import (
	"github.com/gin-gonic/gin"
	"github.com/itsvagapov/service"
)

func main() {
	router := gin.Default()

	router.GET("/students", service.GetAllStudents)
	router.GET("/students/:id", service.GetStudentByID)
	router.POST("/students", service.CreateStudent)
	router.PATCH("/students/:id", service.UpdateStudent)
	router.DELETE("/students/:id", service.DeleteStudent)
	router.GET("/students/:id/notes", service.GetNotesByStudentID)
	router.GET("/students/:id/groups", service.GetStudentsByGroupID)

	router.GET("/groups", service.GetAllGroups)
	router.GET("/groups/:id", service.GetGroupByID)
	router.POST("/groups", service.CreateGroup)
	router.PATCH("/groups/:id", service.UpdateGroup)
	router.DELETE("/groups/:id", service.DeleteGroup)
	router.GET("/groups/:id/stats/offer", service.GetOfferStats)

	router.GET("/notes/:id", service.GetNoteByID)
	router.POST("/notes", service.CreateNote)
	router.PATCH("/notes/:id", service.UpdateNote)
	router.DELETE("/notes/:id", service.DeleteNote)
	router.GET("/students/:id/notes", service.GetNotesByStudentID)

	router.Run(":8080")
}
