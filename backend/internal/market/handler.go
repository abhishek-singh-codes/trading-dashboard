package market

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{service: NewService(NewRepository(db))}
}

func (h *Handler) Search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query must be at least 2 characters"})
		return
	}
	results, err := h.service.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"query": query, "results": results, "count": len(results)})
}

func (h *Handler) GetQuote(c *gin.Context) {
	symbol := strings.ToUpper(strings.TrimSpace(c.Param("symbol")))
	quote, err := h.service.GetQuote(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, quote)
}

func (h *Handler) GetHistory(c *gin.Context) {
	symbol := strings.ToUpper(strings.TrimSpace(c.Param("symbol")))
	period := c.DefaultQuery("period", "30d")
	if period != "1d" && period != "30d" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period must be '1d' or '30d'"})
		return
	}
	history, err := h.service.FetchAndStoreHistory(symbol, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}
