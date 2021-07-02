// Author: Steve Zhang
// Date: 2020/10/16 5:49 下午

package common

import "github.com/gin-gonic/gin"

const (
	ContextKeyForError     = "error"
	ContextKeyForResponse  = "response"
	ContextKeyForUID       = "user_id"
	ContextKeyForAuthToken = "auth_token"
)

func SetLoginContext(c *gin.Context, uid, token string) {
	c.Set(ContextKeyForUID, uid)
	c.Set(ContextKeyForAuthToken, token)
	return
}

func GetLoginContext(c *gin.Context) (uid, token string) {
	uid = c.MustGet(ContextKeyForUID).(string)
	token = c.MustGet(ContextKeyForAuthToken).(string)
	return
}

func SetResponseContext(c *gin.Context, rsp *Response, err error) {
	c.Set(ContextKeyForError, err)
	c.Set(ContextKeyForResponse, rsp)
	return
}

func GetResponseContext(c *gin.Context) (rsp *Response, err error) {
	ierr := c.MustGet(ContextKeyForError)
	if ierr != nil {
		err = ierr.(error)
	}
	irsp := c.MustGet(ContextKeyForResponse)
	if irsp != nil {
		rsp = irsp.(*Response)
	}
	return
}
