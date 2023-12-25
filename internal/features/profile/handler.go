package profile

import (
	"errors"
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

func (h *Handler) GetProfile(loginPage string, editAccontURL string, editEmployeeURL string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, loginPage)
			}
			return c.NoContent(http.StatusBadRequest)
		}

		userInfo, err := h.svc.GetUserInfo(c.Request().Context(), id)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return profilePage(userInfo, editAccontURL, editEmployeeURL).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) ServeProfileRequest(back string, onProfile string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, back)
			}
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

func (h *Handler) ServeAccountUpdate(saveHandler string, back string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, back)
			}
			return c.NoContent(http.StatusBadRequest)
		}
		info, err := h.svc.GetUserInfo(c.Request().Context(), id)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return editAccountInfo(info, saveHandler, back).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) ServeEmployeeUpdate(saveHandler string, back string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func (h *Handler) getId(c echo.Context) (uuid.UUID, error) {
	userId := h.sessionManager.GetString(c.Request().Context(), "user-id")
	if userId == "" {
		return uuid.UUID{}, errors.New("empty")
	}
	id, err := uuid.Parse(userId)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid")
	}
	return id, nil
}
