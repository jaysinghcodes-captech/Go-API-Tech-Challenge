package handlers

import (
	"strconv"
	"net/http"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

// HandleGetCourseByID returns course by its ID
func HandleGetCourseByID(logger *httplog.Logger, svsCourse *services.CourseService) http.HandlerFunc {
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

		course, err := svsCourse.GetCourseById(ctx, courseIDInt)
		if err != nil {
			logger.Error("error getting course", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error getting course",
			})
			return
		}

		encodeResponse(w, logger, http.StatusOK, responseCourse{Course: mapOutputCourse(course)})
	}
}
