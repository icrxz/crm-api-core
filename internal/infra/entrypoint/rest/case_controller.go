package rest

import (
	"github.com/gin-gonic/gin"
)

type CaseController struct {
}

func NewCaseController() *CaseController {
	return &CaseController{}
}

func (c *CaseController) CreateCase(ctx *gin.Context) {

}

func (c *CaseController) GetCase(ctx *gin.Context) {

}

func (c *CaseController) SearchCases(ctx *gin.Context) {

}
