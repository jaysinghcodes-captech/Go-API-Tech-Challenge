package handlers

import (
	"net/http"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"

	
	"github.com/go-chi/httplog/v2"
)

// HandleListPersons is a handler that returns a list of persons
func HandleListPersons(logger *httplog.Logger, svsPerson *services.PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		persons, err := svsPerson.ListPersons(ctx)
		if err != nil {
			logger.Error("error getting all persons", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error retrieving data",
			})
			return
		}
		
		personsOut := mapMultipleOutputPersons(persons)
		encodeResponse(w, logger, http.StatusOK, responsePersons{Persons: personsOut})
	}
}