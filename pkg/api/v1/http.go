package v1

import (
	"net/http"
	"strconv"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/db"
	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/moedersvoormoeders/api.mvm.digital/pkg/api/auth"
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
}

func (h *HTTPHandler) checkAuth(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.Claim)
	if claims.Name == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"status": "JWT incorrect"})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
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
