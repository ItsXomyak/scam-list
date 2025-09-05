package handler

import (
	"net/http"

	"github.com/ItsXomyak/scam-list/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Verify struct {
	verifier Verifier
	log      logger.Logger
}

func NewVerify(verifier Verifier, log logger.Logger) *Verify {
	return &Verify{
		verifier: verifier,
		log:      log,
	}
}

func (h *Verify) VerifyDomain(c *gin.Context) {
	ctx := logger.WithAction(c.Request.Context(), "verify_domain")

	domain := c.Param("domain")

	err := h.verifier.VerifyDomain(domain)
	if err != nil {
		h.log.Error(logger.ErrorCtx(ctx, err), "error verifying domain", err, "domain", domain)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	// Example response
	c.JSON(http.StatusOK, gin.H{
		"domain":  domain,
		"is_scam": false,
		"reason":  "not found in database",
	})
}
