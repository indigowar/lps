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

func (h *Handler) HandleAccountUpdate(onProfile string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, onProfile)
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

		c.Response().Header().Add("Hx-Redirect", onProfile)
		return c.NoContent(http.StatusAccepted)
	}
}

func (h *Handler) ServeAccountUpdateForm(saveHandler string, cancelHandler string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, cancelHandler)
			}
			return c.NoContent(http.StatusBadRequest)
		}
		info, err := h.svc.GetUserInfo(c.Request().Context(), id)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return editAccountInfo(info, saveHandler, cancelHandler).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) ServeEmployeeUpdateForm(saveHandler string, handlerCancel string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, handlerCancel)
			}
			return c.NoContent(http.StatusBadRequest)
		}
		info, err := h.svc.GetUserInfo(c.Request().Context(), id)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return editEmployeeInfo(info, saveHandler, handlerCancel).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleEmployeeUpdate(onProfile string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, onProfile)
			}
			return c.Redirect(http.StatusSeeOther, onProfile)
		}

		surname := c.FormValue("surname")
		name := c.FormValue("name")
		phone_number := c.FormValue("phone_number")

		var patronymic *string = nil
		if c.FormValue("patronymic") != "" {
			p := c.FormValue("patronymic")
			patronymic = &p
		}

		if err := h.svc.UpdateEmployee(c.Request().Context(), id, surname, name, patronymic, phone_number); err != nil {
			log.Println(err)
			return c.Redirect(http.StatusBadRequest, onProfile)
		}

		c.Response().Header().Add("HX-Redirect", onProfile)
		return c.NoContent(http.StatusAccepted)
	}
}

func (h *Handler) HandleEditCancellation(onProfile string) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, onProfile)
			}
			return c.Redirect(http.StatusSeeOther, onProfile)
		}

		c.Response().Header().Add("HX-Redirect", onProfile)
		return c.NoContent(http.StatusAccepted)
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
