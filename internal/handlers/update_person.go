package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/models"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"
)

// HandleUpdatePerson updates person by their firstName
func HandleUpdatePerson(logger *httplog.Logger, svsPerson *services.PersonService) http.HandlerFunc {
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

		personIn, problems, err := decodeValidateBody[inputPerson, models.Person](r)
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

		updatedPerson, err := svsPerson.UpdatePerson(ctx, firstName, models.Person{
			FirstName: personIn.FirstName,
			LastName:  personIn.LastName,
			Type:      personIn.Type,
			Age:       personIn.Age,
			Courses:   personIn.Courses,
		})
		if err != nil {
			logger.Error("error updating person", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error updating person",
			})
			return
		}

		encodeResponse(w, logger, http.StatusOK, updatedPerson)
	}
}
