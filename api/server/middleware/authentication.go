package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authentication 身份认证
/*func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if AlwaysAllowPath.Has(c.Request.URL.Path) {
			return
		}

		r := httputils.NewResponse()
		token, err := extractToken(c, c.Request.URL.Path == types.WebShellURL)
		if err != nil {
			r.SetCode(http.StatusUnauthorized)
			httputils.SetFailed(c, r, err)
			c.Abort()
			return
		}
		claims, err := httputils.ParseToken(token, oats.CoreV1.User().GetJWTKey())
		if err != nil {
			r.SetCode(http.StatusUnauthorized)
			httputils.SetFailed(c, r, err)
			c.Abort()
			return
		}

		// 保存用户id
		c.Set(types.UserId, claims.Id)
	}
}*/

// 从请求头中获取 token
func extractToken(c *gin.Context, ws bool) (string, error) {
	emptyFunc := func(t string) bool { return len(t) == 0 }
	if ws {
		wsToken := c.GetHeader("Sec-WebSocket-Protocol")
		if emptyFunc(wsToken) {
			return "", fmt.Errorf("authorization header is not provided")
		}
		return wsToken, nil
	}

	token := c.GetHeader("Authorization")
	if emptyFunc(token) {
		return "", fmt.Errorf("authorization header is not provided")
	}
	fields := strings.Fields(token)
	if len(fields) != 2 {
		return "", fmt.Errorf("invalid authorization header format")
	}
	if fields[0] != "Bearer" {
		return "", fmt.Errorf("unsupported authorization type")
	}

	return fields[1], nil
}
