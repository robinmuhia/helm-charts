package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/domain"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/usecases"
)

type HandlersInterfacesImpl struct {
	usecase *usecases.UsecaseHelmService
}

func NewHandlersInterfaces(usecases *usecases.UsecaseHelmService) *HandlersInterfacesImpl {
	return &HandlersInterfacesImpl{
		usecase: usecases,
	}
}

func (h HandlersInterfacesImpl) ParseHelmLink(c *gin.Context) {
	urlLink := domain.HelmLinkInput{}

	err := c.BindJSON(&urlLink)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	images, err := h.usecase.ProcessHelmChart(c.Request.Context(), urlLink)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, images)
}
