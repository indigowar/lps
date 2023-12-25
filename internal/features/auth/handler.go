package auth

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc            Service
	sessionManager *scs.SessionManager
}

func NewHandler(svc Service, sessionManager *scs.SessionManager) *Handler {
	return &Handler{
		svc:            svc,
		sessionManager: sessionManager,
	}
}

func (h *Handler) ServeLoginPage(handler string) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := h.sessionManager.GetString(c.Request().Context(), "user-id")
		log.Println("USER ID IS ", userId)
		if userId != "" {
			return c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return loginPage(handler).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleLoginRequest() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := h.sessionManager.GetString(c.Request().Context(), "user-id")
		if userId != "" {
			return c.NoContent(http.StatusBadRequest)
		}

		login := c.FormValue("login")
		password := c.FormValue("password")

		id, err := h.svc.Login(c.Request().Context(), login, password)
		if err != nil {
			log.Printf("Failed to login, due to: %s\n", err)
			return echo.NewHTTPError(http.StatusBadRequest, "credentials are invalid")
		}

		h.sessionManager.Put(c.Request().Context(), "user-id", id.String())

		c.Response().Header().Add("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	}
}

func (h *Handler) ServeRegisterPage(handler string) echo.HandlerFunc {
	return func(c echo.Context) error {
		// currentSession, _ := session.Get("user-session", c)
		// _, exists := currentSession.Values["user-id"]
		// if exists {
		// 	return c.Redirect(http.StatusPermanentRedirect, "/")
		// }

		login := c.Param("login")

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return registrationPage(handler, login).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleRegisterRequest() echo.HandlerFunc {
	return func(c echo.Context) error {
		login := c.FormValue("login")
		password := c.FormValue("password")

		log.Println("Found: ", login, " and ", password)

		return nil
	}
}

func (h *Handler) HandleLogout() echo.HandlerFunc {
	return func(c echo.Context) error {
		// currentSession, _ := session.Get("user-session", c)

		// if _, exists := currentSession.Values["user-id"]; !exists {
		// 	return echo.NewHTTPError(http.StatusUnauthorized, "<h1>unauthorized</h1>")
		// }
		// delete(currentSession.Values, "user-id")
		// currentSession.Save(c.Request(), c.Response().Writer)
		return c.NoContent(http.StatusOK)
	}
}
