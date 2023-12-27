package addincident

import (
	"errors"
	"log"
	"lps/internal/domain/usecases"
	"lps/pkg/templates/utils"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	sm *scs.SessionManager

	create usecases.CreateIncidentUseCase
}

func NewHandler(sm *scs.SessionManager, create usecases.CreateIncidentUseCase) Handler {
	return Handler{
		sm:     sm,
		create: create,
	}
}

func (h *Handler) ServePage() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, "/auth/login")
			}
			return c.NoContent(http.StatusBadRequest)
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return formPage(id).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleRequest() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.FormValue("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		description := c.FormValue("description")
		date := c.FormValue("date")

		if err := h.create.CreateIncident(c.Request().Context(), id, description, date); err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}
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
