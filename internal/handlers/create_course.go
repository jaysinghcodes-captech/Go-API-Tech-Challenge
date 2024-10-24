package handlers

import (
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

// HandleCreateCourse creates a new course for a person
func HandleCreateCourse(logger *httplog.Logger, svsCourse *services.CourseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		courseIn, problems, err := decodeValidateBody[inputCourse](r)
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

		course, err := svsCourse.CreateCourse(ctx, courseIn.Name)
		if err != nil {
			logger.Error("error creating course", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error creating course",
			})
			return
		}

		encodeResponse(w, logger, http.StatusCreated, responseCourse{Course: mapOutputCourse(course)})
	}
}