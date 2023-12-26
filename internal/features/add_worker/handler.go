package addworker

import (
	"errors"
	"log"
	"lps/pkg/templates/utils"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

func (h *Handler) ServePage(handlerURL string) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return c.Redirect(http.StatusSeeOther, "/auth/login")
			}
			return c.NoContent(http.StatusBadRequest)
		}

		canPerfom, err := h.svc.ConfirmUserCanAddWorker(c.Request().Context(), id)
		if err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		if !canPerfom {
			return c.Redirect(http.StatusSeeOther, "/denied")
		}

		deparments, err := h.svc.GetDepartments(c.Request().Context())
		if err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}

		positions, err := h.svc.GetPositions(c.Request().Context())
		if err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
		return addWorkerPage(handlerURL, positions, deparments).Render(c.Request().Context(), c.Response().Writer)
	}
}

func (h *Handler) HandleRequest() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getId(c)
		if err != nil {
			if err.Error() == "empty" {
				return utils.Handle404(c)
			}
			return c.NoContent(http.StatusBadRequest)
		}

		canPerfom, err := h.svc.ConfirmUserCanAddWorker(c.Request().Context(), id)
		if err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		if !canPerfom {
			return c.NoContent(http.StatusUnauthorized)
		}

		login := c.FormValue("login")
		surname := c.FormValue("surname")
		name := c.FormValue("name")
		patronymicValue := c.FormValue("patronymic")
		phone := c.FormValue("phone_number")
		positionStr := c.FormValue("position")
		departmentStr := c.FormValue("department")

		var patronymic *string = nil
		if patronymicValue != "" {
			patronymic = &patronymicValue
		}
		position, err := uuid.Parse(positionStr)
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusBadRequest)
		}
		department, err := uuid.Parse(departmentStr)
		if err != nil {
			log.Println(err)
			return c.NoContent(http.StatusBadRequest)
		}

		log.Println(position)
		log.Println(department)

		err = h.svc.CreateWorker(c.Request().Context(), login, surname, name, patronymic, phone, position, department)
		if err != nil {
			log.Println(err)
			return utils.Handle500(c)
		}
		return c.Redirect(http.StatusAccepted, "/")
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
