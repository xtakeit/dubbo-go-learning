// Author: Steve Zhang
// Date: 2020/10/17 11:45 上午

package middleware

import (
	"go-sever/common"

	"github.com/gin-gonic/gin"
)

// mock用户会话
var sessionList = map[string]string{
	"token_001": "admin",
}

func Auth(c *gin.Context) {
	token := common.GetAuthTokenHeader(c)
	uid, ok := sessionList[token]

	if !ok {
		common.SetResponseContext(c, &common.Response{
			Code:    common.ResponseCodeAuthFailed,
			Message: "登陆失效, 请重新登陆",
		}, nil)
		c.Abort()
		return
	}

	// 设置登陆上下文
	common.SetLoginContext(c, uid, token)

	c.Next()
}
