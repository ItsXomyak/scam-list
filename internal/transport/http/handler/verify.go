package handler

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/ItsXomyak/scam-list/pkg/logger"
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

// VerifyDomain
func (h *Verify) VerifyDomain(c *gin.Context) {
	ctx := logger.WithAction(c.Request.Context(), "handler_verify_domain")

	domain := c.Param("domain")

	if err := validateDomain(domain); err != nil {
		badRequestResponse(c, err.Error())
		return
	}

	// ProcessDomain
	_, err := h.verifier.ProcessDomain(ctx, domain)
	if err != nil {
		h.log.Error(logger.ErrorCtx(ctx, err), "error verifying domain", err, "domain", domain)
		internalErrorResponse(c, "internal server error")
		return
	}

	// Example response
	c.JSON(http.StatusOK, gin.H{
		"domain":  domain,
		"is_scam": false,
	})
}

// checks if the provided domain is valid
func validateDomain(raw string) error {
	if len(raw) == 0 {
		return errors.New("URL must be provided")
	}

	// Parse the URL
	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return errors.New("invalid URL format")
	}

	// Must have scheme
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("URL must have http or https scheme")
	}

	// Must have host
	if parsed.Host == "" {
		return errors.New("URL must have a host")
	}

	// Validate domain with regex (letters, digits, dash, dot)
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	host := parsed.Hostname()
	if !domainRegex.MatchString(host) {
		return errors.New("invalid domain format")
	}
	return nil
}
