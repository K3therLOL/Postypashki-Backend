package server

import (
	"taskserver/manager"

	"github.com/gin-gonic/gin"
)

func CreateAndRun() error {
	r := gin.Default()

	api := manager.NewAPI()

	r.POST("/task", api.ExecuteTask)
	r.GET("/status/:task_id", api.GetTaskStatus)
	r.GET("/result/:task_id", api.GetTaskResult)

	return r.RunTLS(":8080", "cert.pem", "key.pem")
}
