package handlers

import (
	"encoding/json"
	"net/http"
	"go-api-tech-challenge/internal/models"
	"github.com/go-chi/httplog"
)

type outputCourse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
}

type outputPerson struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Type		string `json:"type"`
	Age         int    `json:"age"`
	Courses     []int  `json:"courses"`
}

func mapOutputCourse(course models.Course) outputCourse {
	return outputCourse{
		ID:   course.ID,
		Name: course.Name,
	}
}

func mapMultipleOutputCourses(courses []models.Course) []outputCourse {
	outputCourses := make([]outputCourse, 0, len(courses))
	for _, course := range courses {
		outputCourses = append(outputCourses, mapOutputCourse(course))
	}
	return outputCourses
}

func mapOutputPerson(person models.Person) outputPerson {
	return outputPerson{
		ID:        person.ID,
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Type:      person.Type,
		Age:       person.Age,
		Courses:   person.Courses,
	}
}

func mapMultipleOutputPersons(persons []models.Person) []outputPerson {
	outputPersons := make([]outputPerson, 0, len(persons))
	for _, person := range persons {
		outputPersons = append(outputPersons, mapOutputPerson(person))
	}
	return outputPersons
}

type responseCourse struct {
	Course outputCourse `json:"data"`
}

type responseCourses struct {
	Courses []outputCourse `json:"data"`
}

type responsePerson struct {
	Person outputPerson `json:"data"`
}

type responsePersons struct {
	Persons []outputPerson `json:"data"`
}

type responseMessage struct {
	Message string `json:"message"`
}

type responseErr struct {
	Error            string    `json:"error,omitempty"`
	ValidationErrors []problem `json:"validation_errors,omitempty"`
}

func encodeResponse(w http.ResponseWriter, logger *httplog.Logger, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error while marshaling data", "err", err, "data", data)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
	}
}

