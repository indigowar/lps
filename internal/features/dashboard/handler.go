package dashboard

import (
	"errors"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"lps/pkg/templates/utils"
)

type Handler struct {
	svc Service
	sm  *scs.SessionManager
}

func NewHandler(svc Service, sm *scs.SessionManager) Handler {
	return Handler{
		svc: svc,
		sm:  sm,
	}
}

func (h *Handler) ShowDashboard() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, "/auth/login")
			}
			return utils.Handle500(c)
		}

		role, err := h.svc.GetUserRole(c.Request().Context(), id)
		if err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		return h.showDashboardForRole(c, id, role)
	}
}

func (h *Handler) ServePositionsTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		positions, err := h.svc.GetPositions(c.Request().Context())
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		if err != nil {
			return utils.Handle500(c)
		}
		return tablePositions(positions).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) showDashboardForRole(c echo.Context, id uuid.UUID, role UserRole) error {
	switch role {
	case UserRoleAdmin:
		return h.showAdminDashboard(c, id)
	case UserRoleHead:
		return h.showHeadDashboard(c, id)
	case UserRoleStaff:
		return h.showStaffDashboard(c, id)
	default:
		return utils.Handle500(c)
	}
}

func (h *Handler) showAdminDashboard(c echo.Context, id uuid.UUID) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return adminDashboard().Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) showHeadDashboard(c echo.Context, id uuid.UUID) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return headDashboard().Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) showStaffDashboard(c echo.Context, id uuid.UUID) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return staffDashboard().Render(c.Request().Context(), c.Response().Writer)
}

func (h *Handler) getId(c echo.Context) (uuid.UUID, error) {
	userId := h.sm.GetString(c.Request().Context(), "user-id")
	if userId == "" {
		return uuid.UUID{}, errors.New("empty")
	}
	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid")
	}
	return id, nil
}
