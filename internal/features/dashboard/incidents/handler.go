package incidents

import (
	"log"
	"lps/internal/domain"
	"lps/internal/domain/usecases"
	"lps/pkg/templates/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type IncidentHandler struct {
	incident usecases.IncidentUseCase
}

func NewIncidentHandler(i usecases.IncidentUseCase) IncidentHandler {
	return IncidentHandler{
		incident: i,
	}
}

func (h *IncidentHandler) ServeTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		incidents, err := h.incident.GetAll(c.Request().Context())
		if err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Table(incidents).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *IncidentHandler) ServeItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		strID := c.Param("id")

		id, err := uuid.Parse(strID)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		incident, err := h.incident.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Item(incident).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *IncidentHandler) ServeEdit() echo.HandlerFunc {
	return func(c echo.Context) error {
		strID := c.Param("id")

		id, err := uuid.Parse(strID)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		incident, err := h.incident.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Form(incident).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *IncidentHandler) HandleEdit() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.FormValue("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		employeeID, err := uuid.Parse(c.FormValue("employee"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		date, err := time.Parse("2006-01-02", c.FormValue("date"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		incident := domain.Incident{
			ID:          id,
			Employee:    employeeID,
			Description: c.FormValue("description"),
			Date:        date,
		}

		if err := h.incident.Update(c.Request().Context(), incident); err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}

		updatedIncident, err := h.incident.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Item(updatedIncident).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *IncidentHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		if err := h.incident.Delete(c.Request().Context(), id); err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}

		return c.NoContent(http.StatusOK)
	}
}
