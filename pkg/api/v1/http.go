package v1

import (
	"strconv"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

var registers []func(e *echo.Echo, h *HTTPHandler)

func init() {
	registers = append(registers, func(e *echo.Echo, h *HTTPHandler) {
		e.GET("/v1/auth/check", h.checkAuth)
	})
}

type HTTPHandler struct {
	db *db.Connection
}

func NewHTTPHandler(db *db.Connection) *HTTPHandler {
	return &HTTPHandler{
		db: db,
	}
}

func (h *HTTPHandler) Register(e *echo.Echo) {
	for _, regFn := range registers {
		regFn(e, h)
	}

	//whoami
	e.GET("/v1/auth/check", h.checkAuth)
	e.GET("/v1/whoami/roles", h.getRoles)
}

func Paginate(c echo.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
