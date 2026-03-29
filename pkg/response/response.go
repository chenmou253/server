package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data any) {
	c.JSON(200, gin.H{"code": 0, "message": "ok", "data": data})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"code": code, "message": message})
}

func ErrorWithData(c *gin.Context, code int, message string, data any) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(code, gin.H{"code": code, "message": message, "data": data})
}
