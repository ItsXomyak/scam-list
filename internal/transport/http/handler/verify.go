package handler

import (
	"net/http"

	"github.com/ItsXomyak/scam-list/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Verify struct {
	log logger.Logger
}

func NewVerify(log logger.Logger) *Verify {
	return &Verify{log: log}
}

func (h *Verify) VerifyDomain(c *gin.Context) {
	ctx := logger.WithAction(c.Request.Context(), "verify_domain")

	domain := c.Param("domain")

	// Your verification logic here
	h.log.Info(ctx, "verifying domain", "domain", domain)

	// Example response
	c.JSON(http.StatusOK, gin.H{
		"domain":  domain,
		"is_scam": false,
		"reason":  "not found in database",
	})
}
