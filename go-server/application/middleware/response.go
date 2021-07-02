// Author: Steve Zhang
// Date: 2020/10/17 12:33 下午

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-sever/common"
)

// Response json响应中间件
func Response(c *gin.Context) {
	common.SetResponseContext(c, common.NewOKResponse(), nil)
	c.Next()
	rsp, _ := common.GetResponseContext(c)
	c.JSON(http.StatusOK, rsp)
}
