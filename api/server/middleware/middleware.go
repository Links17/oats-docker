package middleware

import (
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"

	"oats-docker/pkg/types"
)

var AlwaysAllowPath sets.String

func InitMiddlewares(ginEngine *gin.Engine) {
	// 初始化可忽略的请求路径
	AlwaysAllowPath = sets.NewString(types.HealthURL, types.LoginURL, types.LogoutURL)

	// 依次进行跨域，日志，单用户限速，总量限速，验证，和鉴权
	//ginEngine.Use(Cors(), LoggerToFile(), UserRateLimiter(), Limiter(), Authentication())
	// TODO: 临时关闭
	//if os.Getenv("DEBUG") != "true" {
	//	ginEngine.Use(Rbac())
	//}
}
