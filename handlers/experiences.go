package handlers

import (
	"errors"
	"net/http"

	"emailaddress.horse/thousand/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type newExperienceRenderer interface {
	NewExperience(http.ResponseWriter, models.Memory) error
}

func NewExperience(r *chi.Mux, l *zap.Logger, t newExperienceRenderer, mg memoryGetter) {
	r.Get("/vampires/{vampireID}/memories/{id}/experiences/new", func(w http.ResponseWriter, r *http.Request) {
		vampireID, err := uuid.Parse(chi.URLParam(r, "vampireID"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		memoryID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			l.Error("failed to parse id as UUID", zap.Error(err))
			handleError(w, err)
			return
		}

		memory, err := mg.GetMemory(r.Context(), vampireID, memoryID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				err = NotFoundError.Cause(err)
			}

			l.Error("failed to find memory", zap.Stringer("vampireID", vampireID), zap.Stringer("memoryID", memoryID), zap.Error(err))
			handleError(w, err)
			return
		}

		err = t.NewExperience(w, memory)
		if err != nil {
			l.Error("failed to render", zap.Error(err))
			handleError(w, err)
		}
	})
}

func CreateExperience(e *echo.Echo, ec experienceCreator) {
	e.POST("/vampires/:vampireID/memories/:id/experiences", func(c echo.Context) error {
		vampireID, err := uuid.Parse(c.Param("vampireID"))
		if err != nil {
			return err
		}

		memoryID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return err
		}

		description := c.FormValue("description")

		_, err = ec.CreateExperience(c.Request().Context(), vampireID, memoryID, description)
		if errors.Is(err, models.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Memory could not be found").SetInternal(err)
		} else if err != nil {
			return err
		}

		return c.Redirect(http.StatusSeeOther, "/vampires/"+vampireID.String())
	}).Name = "create-experience"
}
