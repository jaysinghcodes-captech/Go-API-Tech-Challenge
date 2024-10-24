package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"
)

// HandleDeletePerson deletes person by their firstName
func HandleDeletePerson(logger *httplog.Logger, svsPerson *services.PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		firstName := chi.URLParam(r, "firstName")
		if firstName == "" {
			logger.Error("missing person firstName")
			encodeResponse(w, logger, http.StatusBadRequest, responseErr{
				Error: "missing person firstName",
			})
			return
		}

		err := svsPerson.DeletePerson(ctx, firstName)
		if err != nil {
			logger.Error("error deleting person", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error deleting person",
			})
			return
		}

		encodeResponse(w, logger, http.StatusOK, nil)
	}
}
