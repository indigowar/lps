package departments

import (
	"log"
	"lps/internal/domain/usecases"
	"lps/pkg/templates/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase usecases.DepartmentUseCase
}

func NewHandler(usecase usecases.DepartmentUseCase) Handler {
	return Handler{
		usecase: usecase,
	}
}

func (h *Handler) ServeTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		departments, err := h.usecase.GetAll(c.Request().Context())
		if err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return TableDepartments(departments).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) ServeItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		strId := c.Param("id")

		id, err := uuid.Parse(strId)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		d, err := h.usecase.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return ItemDepartment(d).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) ServeEdit() echo.HandlerFunc {
	return func(c echo.Context) error {
		strId := c.Param("id")

		id, err := uuid.Parse(strId)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		d, err := h.usecase.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return FormDepartment(d).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleEdit() echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.FormValue("name")

		id, err := uuid.Parse(c.FormValue("id"))
		if err != nil || len(name) == 0 {
			return c.NoContent(http.StatusBadRequest)
		}

		if err = h.usecase.Update(c.Request().Context(), id, name); err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}

		d, err := h.usecase.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return ItemDepartment(d).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		if err := h.usecase.Delete(c.Request().Context(), id); err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		return c.NoContent(http.StatusOK)
	}
}
