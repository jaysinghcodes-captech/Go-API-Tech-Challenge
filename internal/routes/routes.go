package routes

import (
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/handlers"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

// RegisterRoutes sets up all the API routes
func RegisterRoutes(router *chi.Mux, logger *httplog.Logger, svsCourse *services.CourseService, svsPerson *services.PersonService) {
	// Course-related routes
	router.Route("/api/course", func(router chi.Router) {
		router.Get("/", handlers.HandleListCourses(logger, svsCourse))
		// router.Post("/", handlers.HandleCreateCourse(logger, svsCourse))
		// router.Get("/{id}", handlers.HandleGetCourseByID(logger, svsCourse))
		// router.Put("/{id}", handlers.HandleUpdateCourse(logger, svsCourse))
		// router.Delete("/{id}", handlers.HandleDeleteCourse(logger, svsCourse))
	})

	// Person-related routes
	router.Route("/api/person", func(router chi.Router) {
		router.Get("/", handlers.HandleListPersons(logger, svsPerson))
		// router.Post("/", handlers.HandleCreatePerson(logger, svsPerson))
		// router.Get("/{firstName}", handlers.HandleGetPersonByName(logger, svsPerson))
		// router.Put("/{firstName}", handlers.HandleUpdatePerson(logger, svsPerson))
		// router.Delete("/{firstName}", handlers.HandleDeletePerson(logger, svsPerson))
	})
}
