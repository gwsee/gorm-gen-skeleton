package controller

import (
	"github.com/gin-gonic/gin"
	"gorm-gen-skeleton/app/request"
	"gorm-gen-skeleton/internal/variable/consts"
	"log"
)

var validator *request.Request

type base struct {
}

func init() {
	var err error
	validator, err = request.New()
	if err != nil {
		log.Fatal(consts.ErrorInitConfig)
	}
}

func (base) Validate(ctx *gin.Context, param any) map[string]string {
	return validator.Validator(ctx, param)
}
