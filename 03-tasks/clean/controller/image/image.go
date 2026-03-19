package controller

import (
	"net/http"
	"net/url"
	usecase "task/clean/usecase/image"

	"github.com/gin-gonic/gin"
)

type ImageRequest struct {
	Url string `json:"url"`
}

type ImageController struct {
	usecase *usecase.ImageUsecase
}

func NewImageController(uc *usecase.ImageUsecase) *ImageController {
	return &ImageController{usecase: uc}
}

func (ic *ImageController) Process(ctx *gin.Context) {
	req := ImageRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := url.ParseRequestURI(req.Url)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	img, err := ic.usecase.Process(req.Url)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, img)
}
