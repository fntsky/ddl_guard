package controller

import "github.com/gin-gonic/gin"

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc *UserController) RegisterUserByEmail(ctx *gin.Context) {

}
