package v1

import (
	"UrlShortener/internal/entity"
	"UrlShortener/internal/usecase"
	"UrlShortener/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type urlsRoutes struct {
	ud usecase.UrlsData
	l  *slog.Logger
}

func newUrlsRoutes(handler *gin.RouterGroup, ud usecase.UrlsData, l *slog.Logger) {
	r := &urlsRoutes{ud, l}

	h := handler.Group("/url")
	{
		h.POST("/add", r.UrlsAdd)
		h.GET(":url", r.UrlsGet)
		h.DELETE(":url", r.UrlsDelete)
	}
}

type urlAddInput struct {
	Url      string `json:"url" validate:"required"`
	ShortUrl string `json:"short_url" validate:"required"`
	Expiry   string `json:"expiry" validate:"required"`
}

func (r *urlsRoutes) UrlsAdd(c *gin.Context) {
	var input urlAddInput

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	expiry, err := time.ParseDuration(input.Expiry)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = r.ud.SaveUrl(c.Request.Context(), entity.Url{
		Url:      input.Url,
		ShortUrl: input.ShortUrl,
		Expiry:   expiry,
	})

	if err != nil {
		if errors.Is(err, entity.ErrTimeout) {
			c.AbortWithStatus(http.StatusGatewayTimeout)
			return
		}

		r.l.Error("http - v1 - UrlsAdd", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	type response struct {
		Message string `json:"message"`
	}

	c.JSON(http.StatusOK, response{
		Message: "Success",
	})
}

func (r *urlsRoutes) UrlsGet(c *gin.Context) {
	shortUrl := c.Param("url")

	url, err := r.ud.GetUrl(c.Request.Context(), shortUrl)

	if err != nil {
		if errors.Is(err, entity.ErrTimeout) {
			c.AbortWithStatus(http.StatusGatewayTimeout)
			return
		}

		r.l.Error("http - v1 - UrlsGet", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	type urlGetResponse struct {
		Url      string `json:"url" validate:"required"`
		ShortUrl string `json:"short_url" validate:"required"`
		Expiry   string `json:"expiry" validate:"required"`
	}

	c.JSON(http.StatusOK, urlGetResponse{
		Url:      url.Url,
		ShortUrl: url.ShortUrl,
		Expiry:   url.Expiry.String(),
	})
}

func (r *urlsRoutes) UrlsDelete(c *gin.Context) {
	shortUrl := c.Param("url")

	err := r.ud.DeleteUrl(c.Request.Context(), shortUrl)

	if err != nil {
		if errors.Is(err, entity.ErrTimeout) {
			c.AbortWithStatus(http.StatusGatewayTimeout)
			return
		}

		r.l.Error("http - v1 - UrlsDelete", logger.Err(err))
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	type response struct {
		Message string `json:"message"`
	}

	c.JSON(http.StatusOK, response{
		Message: "Success",
	})
}
