package auth

import (
	"errors"
	"net/http"
	"strings"
	usecase "taskserver/clean/usecase/auth"

	"github.com/gin-gonic/gin"
)

var (
	ErrWithJSON            = errors.New("Couldn't parse json.")
	ErrWithJSONFields      = errors.New("Wrong json.")
	ErrWithAuthToken       = errors.New("No auth token.")
	ErrWithAuthTokenFormat = errors.New("Wrong auth token format.")
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthController struct {
	handler *usecase.AuthHandler
}

func NewAuthController(handler *usecase.AuthHandler) *AuthController {
	return &AuthController{
		handler: handler,
	}
}

func (controller *AuthController) Register(c *gin.Context) {
	var credentials Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrWithJSON.Error(),
		})
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrWithJSONFields.Error(),
		})
		return
	}

	token, err := controller.handler.Register(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})
}

func (controller *AuthController) Login(c *gin.Context) {
	var credentials Credentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrWithJSON.Error(),
		})
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": ErrWithJSONFields.Error(),
		})
		return
	}

	token, err := controller.handler.Login(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (controller *AuthController) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrWithAuthToken.Error(),
			})
			c.Abort()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": ErrWithAuthTokenFormat.Error(),
			})
			c.Abort()
			return
		}

		token := headerParts[1]
		_, err := controller.handler.GetSession(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
