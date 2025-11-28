package handlers

import (
	"bonds-report-service/internal/service"
	"bonds-report-service/lib/valuefromcontext"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Client struct {
	service *service.Client
}

func NewHandlers(service *service.Client) *Client {
	return &Client{service: service}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		chatID := c.GetHeader(valuefromcontext.HeaderChatID)
		if chatID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing X-Chat-ID header"})
			return
		}
		ctx := context.WithValue(c.Request.Context(), valuefromcontext.ChatIdKey, chatID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (h *Client) GetAccountsList(c *gin.Context) {
	ctx := c.Request.Context()
	accountsResponce, err := h.service.GetAccountsList(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get accounts"})
		return
	}

	c.JSON(http.StatusOK, accountsResponce)
	return
}

func (h *Client) GetBondReportsByFifo(c *gin.Context) {
	ctx := c.Request.Context()
	chatID, err := valuefromcontext.GetChatIDFromCtxInt(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect X-ChatId header"})
		return
	}
	err = h.service.GetBondReportsByFifo(ctx, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
	return

}

func (h *Client) GetUSD(c *gin.Context) {
	ctx := c.Request.Context()

	usdResponce, err := h.service.GetUsd(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, usdResponce)
	return

}
func (h *Client) GetBondReports(c *gin.Context) {
	ctx := c.Request.Context()
	chatID, err := valuefromcontext.GetChatIDFromCtxInt(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect X-ChatId header"})
		return
	}
	getBondReportsResponce, err := h.service.GetBondReports(ctx, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, getBondReportsResponce)

}
func (h *Client) GetPortfolioStructure(c *gin.Context) {
	ctx := c.Request.Context()

	portfolioStructuresResonce, err := h.service.GetPortfolioStructureForEachAccount(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, portfolioStructuresResonce)

}
func (h *Client) GetUnionPortfolioStructure(c *gin.Context) {
	ctx := c.Request.Context()

	portgolioStructure, err := h.service.GetUnionPortfolioStructureForEachAccount(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, portgolioStructure)
}
func (h *Client) GetUnionPortfolioStructureWithSber(c *gin.Context) {
	ctx := c.Request.Context()

	portgolioStructure, err := h.service.GetUnionPortfolioStructureWithSber(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, portgolioStructure)
}
