package server

import (
	"taskserver/clean/composure"
	"taskserver/manager"

	"github.com/gin-gonic/gin"
)

func CreateAndRun() error {
	r := gin.Default()

	auth := composure.NewAuthController()
	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)

	api := manager.NewAPI()
	authorized := r.Group("/")
	authorized.Use(auth.Middleware())

	authorized.POST("/task", api.ExecuteTask)
	authorized.GET("/status/:task_id", api.GetTaskStatus)
	authorized.GET("/result/:task_id", api.GetTaskResult)

	return r.RunTLS(":8080", "cert.pem", "key.pem")
}
