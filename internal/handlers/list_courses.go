package handlers

import (
	"net/http"

	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"
	"github.com/go-chi/httplog/v2"
)

// HandleListCourses is a handler that returns a list of courses
func HandleListCourses(logger *httplog.Logger, svsCourse *services.CourseService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		courses, err := svsCourse.ListCourses(ctx)
		if err != nil {
			logger.Error("error getting all courses", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}

		coursesOut := mapMultipleOutputCourses(courses)
		encodeResponse(w, logger, http.StatusOK, responseCourses{Courses: coursesOut})
	}
}
