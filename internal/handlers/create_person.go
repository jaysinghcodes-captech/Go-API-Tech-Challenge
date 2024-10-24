package handlers

import (
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/services"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/models"

	"net/http"

	"github.com/go-chi/httplog/v2"
)

// HandleCreatePerson creates a new person
func HandleCreatePerson(logger *httplog.Logger, svsPerson *services.PersonService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		personIn, problems, err := decodeValidateBody[inputPerson](r)
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

		person, err := svsPerson.CreatePerson(ctx, models.Person{
			FirstName: personIn.FirstName,
			LastName:  personIn.LastName,
			Type:      personIn.Type,
			Age:       personIn.Age,
			Courses:   personIn.Courses,
		})
		if err != nil {
			logger.Error("error creating person", "error", err)
			encodeResponse(w, logger, http.StatusInternalServerError, responseErr{
				Error: "Error creating person",
			})
			return
		}

		encodeResponse(w, logger, http.StatusCreated, responsePerson{Person: mapOutputPerson(person)})
	}
}