package domain

import "github.com/gin-gonic/gin"

type Controller interface {

	Setup(r *gin.Engine)

}
