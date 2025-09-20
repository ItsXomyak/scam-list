package handler

import (
	"context"
	"net/http"

	"github.com/ItsXomyak/scam-list/internal/adapter/http/handler/dto"
	"github.com/ItsXomyak/scam-list/internal/domain/entity"
	"github.com/ItsXomyak/scam-list/pkg/logger"
	"github.com/gin-gonic/gin"
)

type DomainService interface {
	CreateDomain(ctx context.Context, params entity.CreateDomainParams) (*entity.Domain, error)
}

type adminPanel struct {
	domain DomainService
	log    logger.Logger
}

func NewAdminPanel(log logger.Logger) *adminPanel {
	return &adminPanel{
		log: log,
	}
}

func (h *adminPanel) CreateDomain(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = logger.WithAction(ctx, "admin_create_domain")

	req := &dto.CreateDomainRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		badRequestResponse(c, err.Error())
		return
	}

	r, err := h.domain.CreateDomain(ctx, *dto.FromCreateRequestToInternal(req))
	if err != nil {
		h.log.Error(logger.ErrorCtx(ctx, err), "failed to create domain", err)
		badRequestResponse(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, r)
}
