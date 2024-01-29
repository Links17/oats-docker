package healthz

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *healthzRouter) healthz(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
