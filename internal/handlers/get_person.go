package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"
)

// HandleGetPersonByName returns person by their first name
func HandleGetPersonByName(logger *httplog.Logger, svsPerson *services.PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		firstName := chi.URLParam(r, "firstName")
		if firstName == "" {
			logger.Error("missing first name")
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "missing first name",
			})
			return
		}

		person, err := svsPerson.GetPersonByFirstName(ctx, firstName)
		if err != nil {
			logger.Error("error getting person", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error getting person",
			})
			return
		}

		encodeResponse(w, logger, http.StatusOK, responsePerson{Person: mapOutputPerson(person)})
	}
}
