package handlers

import (
	"net/http"
	"strconv"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

// HandleUpdateCourse updates course by its ID
func HandleUpdateCourse(logger *httplog.Logger, svsCourse *services.CourseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		courseID := chi.URLParam(r, "id")
		if courseID == "" {
			logger.Error("missing course ID")
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "missing course ID",
			})
			return
		}

		// convert stringID to intID
		courseIDInt, err := strconv.Atoi(courseID)
		if err != nil {
			logger.Error("invalid course ID", "error", err)
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "invalid course ID",
			})
			return
		}

		courseIn, problems, err := decodeValidateBody[inputCourse, models.Course](r)
		if err != nil {
			switch {
			case len(problems) > 0:
				logger.Error("Problems validating input", "error", err, "problems", problems)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					ValidationErrors: problems,
				})
			default:
				logger.Error("BodyParser error", "error", err)
				encodeResponse(w, logger, http.StatusBadRequest, responseErr{
					Error: "missing values or malformed body",
				})
			}
			return
		}

		updatedCourse, err := svsCourse.UpdateCourse(ctx, courseIDInt, courseIn.Name)
		if err != nil {
			logger.Error("error updating course", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error updating course",
			})
			return
		}

		encodeResponse(w, logger, http.StatusOK, responseCourse{Course: mapOutputCourse(updatedCourse)})
	}
}
