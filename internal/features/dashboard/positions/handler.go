package positions

import (
	"log"
	"lps/internal/domain/usecases"
	"lps/pkg/templates/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase usecases.PositionUseCase
}

func NewHandler(u usecases.PositionUseCase) Handler {
	return Handler{
		usecase: u,
	}
}
func (h *Handler) ServeTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		positions, err := h.usecase.GetAll(c.Request().Context())
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

		p, err := h.usecase.Get(c.Request().Context(), id)
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

		p, err := h.usecase.Get(c.Request().Context(), id)
		if err != nil {
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return Form(p).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleEdit() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.FormValue("id"))
		title := c.FormValue("title")
		level := c.FormValue("level")
		if err != nil || len(title) == 0 {
			return c.NoContent(http.StatusBadRequest)
		}

		if err = h.usecase.Update(c.Request().Context(), id, title, level); err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}

		d, err := h.usecase.Get(c.Request().Context(), id)
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
		if err := h.usecase.Delete(c.Request().Context(), id); err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		return c.NoContent(http.StatusOK)
	}
}
