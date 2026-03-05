package manager

import (
	"errors"
	"net/http"
	repository "task/repository/rai"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type API struct {
	repo *repository.Repository
}

var (
	ErrWrongTaskID  = errors.New("Invalid task_id.")
	ErrStatusAccess = errors.New("Could not get status.")
	ErrSaveTaskID   = errors.New("Could not save task_id.")
)

func NewAPI() *API {
	return &API{
		repo: repository.NewRepository(),
	}
}

func (api *API) formTask(uuid uuid.UUID) {
	time.Sleep(30 * time.Second)
	api.repo.Update(uuid)
}

func (api *API) ExecuteTask(c *gin.Context) {
	taskID := uuid.New()

	go api.formTask(taskID)

	if err := api.repo.Save(taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": taskID.String(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID.String(),
	})
}

func (api *API) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	uuid, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrWrongTaskID.Error(),
		})
		return
	}

	taskobj, err := api.repo.Get(uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrStatusAccess.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": taskobj.Status,
	})
}

func (api *API) GetTaskResult(c *gin.Context) {
	taskID := c.Param("task_id")
	c.JSON(http.StatusOK, gin.H{
		"result":  "soon here will be result",
		"task_id": taskID,
	})
}
