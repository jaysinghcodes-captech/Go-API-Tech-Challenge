package routes

import (
	"go-api-tech-challenge/internal/handlers"
	"go-api-tech-challenge/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
)

type Option func(*routerOptions)

type routerOptions struct {
	registerHealthRoute bool
}

func WithRegisterHealthRoute(registerHealthRoute bool) Option {
	return func(options *routerOptions) {
		options.registerHealthRoute = registerHealthRoute
	}
}

func RegisterRoutes(router *chi.Mux, logger *httplog.logger, svsCourse *services.CourseService, svsPerson *services.PersonService, opts ...Option) {
	options := routerOptions{
		registerHealthRoute: true,
	}
	for _, opt := range opts {
		opt(&options)
	}

	router.Route("/api", func(router chi.Router) {
		if options.registerHealthRoute {
			router.Get("/health-check", handlers.HandleHealth(logger))
		}

		router.Route("/course", func(router chi.Router) {
			router.Get("/", handlers.HandleListCourses(logger, svsCourse))
			router.Post("/", handlers.HandleCreateCourse(logger, svsCourse))
			router.Get("/{id}", handlers.HandleGetCourseByID(logger, svsCourse))
			router.Put("/{id}", handlers.HandleUpdateCourse(logger, svsCourse))
			router.Delete("/{id}", handlers.HandleDeleteCourse(logger, svsCourse))
		})

		router.Route("/person", func(router chi.Router) {
			router.Get("/", handlers.HandleListPersons(logger, svsPerson))
			router.Post("/", handlers.HandleCreatePerson(logger, svsPerson))
			router.Get("/{firstName}", handlers.HandleGetPersonByName(logger, svsPerson))
			router.Put("/{firstName}", handlers.HandleUpdatePerson(logger, svsPerson))
			router.Delete("/{firstName}", handlers.HandleDeletePerson(logger, svsPerson))
		})
	
	})
	
}


