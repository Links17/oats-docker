package healthz

import "github.com/gin-gonic/gin"

// healthzRouter is a router to talk with the healthz controller
type healthzRouter struct{}

// NewRouter initializes a new healthz router
func NewRouter(ginEngine *gin.Engine) {
	s := &healthzRouter{}
	s.initRoutes(ginEngine)
}

func (h *healthzRouter) initRoutes(ginEngine *gin.Engine) {
	healthzRoute := ginEngine.Group("/healthz")
	{
		// main process healthz check
		healthzRoute.GET("", h.healthz)
	}
}
