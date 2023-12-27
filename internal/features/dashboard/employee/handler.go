package staff

import (
	"log"
	"lps/internal/domain"
	"lps/internal/domain/usecases"
	"lps/pkg/templates/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	employee    usecases.EmployeeUseCase
	departments usecases.DepartmentUseCase
	positions   usecases.PositionUseCase
}

func NewHandler(u usecases.EmployeeUseCase, d usecases.DepartmentUseCase, p usecases.PositionUseCase) Handler {
	return Handler{
		employee:    u,
		departments: d,
		positions:   p,
	}
}
func (h *Handler) ServeTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		positions, err := h.employee.GetAll(c.Request().Context())
		if err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Table(positions).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) ServeItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		strId := c.Param("id")

		id, err := uuid.Parse(strId)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		p, err := h.employee.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Item(p).Render(c.Request().Context(), c.Response().Writer)

	}
}

func (h *Handler) ServeEdit() echo.HandlerFunc {
	return func(c echo.Context) error {
		strId := c.Param("id")

		id, err := uuid.Parse(strId)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		e, err := h.employee.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		d, err := h.departments.GetAll(c.Request().Context())
		if err != nil {
			return utils.Handle500(c)
		}

		p, err := h.positions.GetAll(c.Request().Context())
		if err != nil {
			return utils.Handle500(c)

		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Form(e, d, p).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleEdit() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.FormValue("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		position, err := uuid.Parse(c.FormValue("position"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		department, err := uuid.Parse(c.FormValue("department"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		employee := domain.Employee{
			ID:         id,
			Surname:    c.FormValue("surname"),
			Name:       c.FormValue("name"),
			Position:   position,
			Department: department,
		}

		patr := c.FormValue("patronymic")
		if patr != "" {
			employee.Patronymic = &patr
		}

		if err = h.employee.Update(c.Request().Context(), employee); err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}

		d, err := h.employee.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Item(d).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		if err := h.employee.Delete(c.Request().Context(), id); err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		return c.NoContent(http.StatusOK)
	}
}
