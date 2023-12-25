package profile

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	sessionManager *scs.SessionManager
	svc            Service
}

func NewHandler(svc Service, sm *scs.SessionManager) Handler {
	return Handler{
		sessionManager: sm,
		svc:            svc,
	}
}

func (h *Handler) GetProfile(loginPage string, editFormUrl string) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := h.sessionManager.GetString(c.Request().Context(), "user-id")
		if userId == "" {
			return c.Redirect(http.StatusSeeOther, loginPage)
		}
		id, err := uuid.Parse(userId)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		accountInfo, _, err := h.svc.GetUserInfo(c.Request().Context(), id)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Page(accountInfo, editFormUrl).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) ServeProfileRequest(back string, onProfile string) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := h.sessionManager.GetString(c.Request().Context(), "user-id")
		if userId == "" {
			return c.Redirect(http.StatusSeeOther, back)
		}
		id, err := uuid.Parse(userId)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		login := c.FormValue("login")
		password := c.FormValue("password")
		oldPassword := c.FormValue("old_password")
		if err := h.svc.UpdateAccount(c.Request().Context(), id, oldPassword, login, password); err != nil {
			log.Println("err: ", err)
			return c.Redirect(http.StatusSeeOther, onProfile)
		}
		return c.Redirect(http.StatusSeeOther, onProfile)
	}
}

func (h *Handler) ServeUpdateForm(saveHandler string, back string) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := h.sessionManager.GetString(c.Request().Context(), "user-id")
		if userId == "" {
			return c.Redirect(http.StatusSeeOther, back)
		}
		id, err := uuid.Parse(userId)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		accountInfo, _, err := h.svc.GetUserInfo(c.Request().Context(), id)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return editAccountInfo(accountInfo, saveHandler, back).Render(c.Request().Context(), c.Response().Writer)
	}
}
