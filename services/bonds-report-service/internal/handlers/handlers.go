package handlers

import (
	"bonds-report-service/internal/service"
	"bonds-report-service/lib/valuefromcontext"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Client struct {
	logger  *slog.Logger
	service *service.Client
}

func NewHandlers(logger *slog.Logger, service *service.Client) *Client {
	return &Client{
		logger:  logger,
		service: service}
}

func (h *Client) GetAccountsList(c *gin.Context) {
	const op = "handlers.GetAccountsList"
	ctx := c.Request.Context()
	accountsResponce, err := h.service.GetAccountsList(ctx)
	if err != nil {
		h.logger.Error("failed to get accounts list",
			slog.String("op", op),
			slog.Any("error", err),
			slog.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get accounts"})
		return
	}

	c.JSON(http.StatusOK, accountsResponce)
}

func (h *Client) GetBondReportsByFifo(c *gin.Context) {
	const op = "handlers.GetBondReportsByFifo"
	ctx := c.Request.Context()
	chatID, err := valuefromcontext.GetChatIDFromCtxInt(ctx)
	if err != nil {
		h.logger.Error("incorrect X-ChatId header",
			slog.String("op", op),
			slog.Any("error", err),
			slog.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect X-ChatId header"})
		return
	}
	err = h.service.GetBondReportsByFifo(ctx, chatID)
	if err != nil {
		h.logger.Error("internal server error",
			slog.String("op", op),
			slog.Any("error", err),
			slog.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Client) GetUSD(c *gin.Context) {
	const op = "handlers.GetUSD"
	ctx := c.Request.Context()

	usdResponce, err := h.service.GetUsd(ctx)
	if err != nil {
		h.logger.Error("internal server error",
			slog.String("op", op),
			slog.Any("error", err),
			slog.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, usdResponce)
}
func (h *Client) GetBondReports(c *gin.Context) {
	const op = "handlers.GetBondReports"
	ctx := c.Request.Context()
	logg := h.logger.With(
		slog.String("op", op),
		slog.String("path", c.Request.URL.Path))

	chatID, err := valuefromcontext.GetChatIDFromCtxInt(ctx)
	if err != nil {
		logg.Error(
			"incorrect X-ChatId header",
			slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect X-ChatId header"})
		return
	}
	getBondReportsResponce, err := h.service.GetBondReports(ctx, chatID)
	if err != nil {
		logg.Error("internal server error",
			slog.Any("error", err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, getBondReportsResponce)

}
func (h *Client) GetPortfolioStructure(c *gin.Context) {
	const op = "handlers.GetPortfolioStructure"
	ctx := c.Request.Context()

	portfolioStructuresResonce, err := h.service.GetPortfolioStructureForEachAccount(ctx)
	if err != nil {
		h.logger.Error("internal server error",
			slog.String("op", op),
			slog.Any("error", err),
			slog.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, portfolioStructuresResonce)

}
func (h *Client) GetUnionPortfolioStructure(c *gin.Context) {
	const op = "handlers.GetUnionPortfolioStructure"
	ctx := c.Request.Context()

	portgolioStructure, err := h.service.GetUnionPortfolioStructureForEachAccount(ctx)
	if err != nil {
		h.logger.Error("internal server error",
			slog.String("op", op),
			slog.Any("error", err),
			slog.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, portgolioStructure)
}
func (h *Client) GetUnionPortfolioStructureWithSber(c *gin.Context) {
	const op = "handlers.GetUnionPortfolioStructureWithSber"
	ctx := c.Request.Context()

	portgolioStructure, err := h.service.GetUnionPortfolioStructureWithSber(ctx)
	if err != nil {
		h.logger.Error("internal server error",
			slog.String("op", op),
			slog.Any("error", err),
			slog.String("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, portgolioStructure)
}
