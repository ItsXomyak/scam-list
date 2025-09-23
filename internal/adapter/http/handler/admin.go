package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ItsXomyak/scam-list/internal/adapter/http/handler/dto"
	"github.com/ItsXomyak/scam-list/internal/domain/entity"
	"github.com/ItsXomyak/scam-list/pkg/logger"
	"github.com/ItsXomyak/scam-list/pkg/validator"
	"github.com/gin-gonic/gin"
)

type DomainRepository interface {
	CreateDomain(ctx context.Context, params *entity.CreateDomainParams) (*entity.Domain, error)
	GetAllDomains(ctx context.Context) ([]*entity.Domain, error)
	GetDomain(ctx context.Context, domain string) (*entity.Domain, error)
	UpdateDomain(ctx context.Context, updated *entity.Domain) (*entity.Domain, error)
	DeleteDomain(ctx context.Context, domain string) error
}

type AdminPanel struct {
	domain DomainRepository
	log    logger.Logger
}

func NewAdminPanel(domain DomainRepository, log logger.Logger) *AdminPanel {
	return &AdminPanel{
		domain: domain,
		log:    log,
	}
}

func (h *AdminPanel) CreateDomain(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = logger.WithAction(ctx, "admin_create_domain")

	req := &dto.CreateDomainRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		badRequestResponse(c, fmt.Errorf("pizdez: %w", err).Error())
		return
	}

	createReq := dto.FromCreateRequestToInternal(req)

	// Validate
	v := validator.New()
	dto.ValidateCreateDomain(v, createReq)
	if !v.Valid() {
		badRequestResponse(c, v.Errors)
		return
	}

	r, err := h.domain.CreateDomain(ctx, createReq)
	if err != nil {
		h.log.Error(logger.ErrorCtx(ctx, err), "failed to create domain", err)

		errCtx := dto.FromError(err)
		errorResponse(c, errCtx.Code, errCtx.Message)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domain": dto.ToDomainResponse(r),
	})
}

func (h *AdminPanel) GetAllDomains(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = logger.WithAction(ctx, "admin_get_all_domains")

	all, err := h.domain.GetAllDomains(ctx)
	if err != nil {
		h.log.Error(logger.ErrorCtx(ctx, err), "failed to get all domains", err)
		errCtx := dto.FromError(err)
		errorResponse(c, errCtx.Code, errCtx.Message)
		return
	}

	if all == nil {
		notFoundResponse(c, "not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domains": dto.ToBatchDomainResponse(all),
	})
}

func (h *AdminPanel) GetDomain(c *gin.Context) {
	ctx := logger.WithAction(c.Request.Context(), "admin_get_domain")

	domain := c.Param("domain")
	if domain == "" {
		badRequestResponse(c, "missing path param: domain")
		return
	}

	d, err := h.domain.GetDomain(ctx, domain)
	if err != nil {
		h.log.Error(logger.ErrorCtx(ctx, err), "failed to get domain", err)

		errCtx := dto.FromError(err)
		errorResponse(c, errCtx.Code, errCtx.Message)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domain": dto.ToDomainResponse(d),
	})
}

func (h *AdminPanel) PatchDomain(c *gin.Context) {
	ctx := logger.WithAction(c.Request.Context(), "admin_update_domain")

	domain := c.Param("domain")
	if domain == "" {
		badRequestResponse(c, "missing path param: domain")
		return
	}

	cur, err := h.domain.GetDomain(ctx, domain)
	if err != nil {
		h.log.Error(logger.ErrorCtx(ctx, err), "failed to get domain before update", err)
		errCtx := dto.FromError(err)
		errorResponse(c, errCtx.Code, errCtx.Message)
		return
	}

	var req dto.UpdateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequestResponse(c, err.Error())
		return
	}

	// changing only if provided
	if req.Status != nil {
		cur.Status = *req.Status
	}
	if req.CompanyName != nil {
		cur.CompanyName = req.CompanyName // *string
	}
	if req.Country != nil {
		cur.Country = req.Country // *string
	}
	if req.ScamSources != nil {
		cur.ScamSources = req.ScamSources // *[]string
	}
	if req.ScamType != nil {
		cur.ScamType = req.ScamType // *string
	}
	if req.VerifiedBy != nil {
		cur.VerifiedBy = req.VerifiedBy // *string
	}
	if req.VerificationMethod != nil {
		cur.VerificationMethod = req.VerificationMethod // *string
	}
	if req.RiskScore != nil {
		cur.RiskScore = req.RiskScore // *string (или поменяй на *float64, если перейдёшь)
	}
	if req.Reasons != nil {
		cur.Reasons = req.Reasons // *[]string
	}
	if req.Metadata != nil {
		cur.Metadata = req.Metadata
	}

	v := validator.New()
	dto.ValidatePatchDomain(v, cur)

	if !v.Valid() {
		badRequestResponse(c, v.Errors)
		return
	}

	updated, err := h.domain.UpdateDomain(ctx, cur)
	if err != nil {
		h.log.Error(logger.ErrorCtx(ctx, err), "failed to update domain", err)

		errCtx := dto.FromError(err)
		errorResponse(c, errCtx.Code, errCtx.Message)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domain": dto.ToDomainResponse(updated),
	})
}

func (h *AdminPanel) DeleteDomain(c *gin.Context) {
	ctx := logger.WithAction(c.Request.Context(), "admin_delete_domain")

	domain := c.Param("domain")
	if domain == "" {
		badRequestResponse(c, "missing path param: domain")
		return
	}

	if err := h.domain.DeleteDomain(ctx, domain); err != nil {
		h.log.Error(logger.ErrorCtx(ctx, err), "failed to delete domain", err)

		errCtx := dto.FromError(err)
		errorResponse(c, errCtx.Code, errCtx.Message)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
