package middleware


import(
	"github.com/gin-gonic/gin"
	"net/http"
)

type role string

const(
	Admin role = "admin"
	User role = "user"
)

func AdminMiddleware()gin.HandlerFunc{
	return func(c *gin.Context){
		AdminRole:=c.GetString("role")
		if AdminRole!=string(Admin){
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}


func UserMiddleware()gin.HandlerFunc{
	return func(c *gin.Context){
		UserRole:=c.GetString("role")
		if UserRole!=string(User){
			c.JSON(http.StatusUnauthorized,gin.H{"error":"Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}


