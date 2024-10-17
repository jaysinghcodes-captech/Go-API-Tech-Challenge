package handlers

import (
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/models"
)

type inputCourse struct {
	Name        string `json:"name"`
}

type inputPerson struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Type		string `json:"type"`
	Age         int    `json:"age"`
	Courses     []int  `json:"courses"`
}

func (course inputCourse) MapTo() (models.Course, error) {
	return models.Course{
		ID:  0,
		Name: course.Name,
	}, nil
}
func (person inputPerson) MapTo() (models.Person, error) {
	return models.Person{
		ID:  0,
		FirstName: person.FirstName,
		LastName: person.LastName,
		Type: person.Type,
		Age: person.Age,
		Courses: person.Courses,
	}, nil
}	

// valid all fields of an inputCourse struct
func (course inputCourse) Valid() []problem {
	var problems []problem

	// validate Course Name is not blank
	if course.Name == "" {
		problems = append(problems, problem{
			Name:        "name",
			Description: "must not be blank",
		})
	}

	return problems
}

func (person inputPerson) Valid() []problem {
	var problems []problem
	validTypes := map[string]bool{
		"student": true,
		"professor": true,
	}

	// validate Person FirstName is not blank
	if person.FirstName == "" {
		problems = append(problems, problem{
			Name:        "first_name",
			Description: "must not be blank",
		})
	}

	// validate Person LastName is not blank
	if person.LastName == "" {
		problems = append(problems, problem{
			Name:        "last_name",
			Description: "must not be blank",
		})
	}

	// validate Person Type is not blank
	if !validTypes[person.Type] {
		problems = append(problems, problem{
			Name:        "type",
			Description: "must be student or professor",
		})
	}
	// validate Person Age is not blank
	if person.Age <= 0 {
		problems = append(problems, problem{
			Name:        "age",
			Description: "must be greater than 0",
		})
	}

	return problems
}

type problem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Validator interface {
	Valid() (problems []problem)
}

type Mapper[T any] interface {
	MapTo() (T, error)
}

type ValidatorMapper[T any] interface {
	Validator
	Mapper[T]
}

func decodeValidateBody[I ValidatorMapper[O], O any](r *http.Request) (O, []problem, error) {
	var inputModel I

	if err := json.NewDecoder(r.Body).Decode(&inputModel); err != nil {
		return *new(O), nil, fmt.Errorf("[in decodeValidateBody] decode json: %w", err)
	}

	if problems := inputModel.Valid(); len(problems) > 0 {
		return *new(O), problems, fmt.Errorf(
			"[in decodeValidateBody] invalid %T: %d problems", inputModel, len(problems),
		)
	}

	data, err := inputModel.MapTo()
	if err != nil {
		return *new(O), nil, fmt.Errorf(
			"[in decodeValidateBody] error mapping input %T to %T: %w",
			*new(I),
			*new(O),
			err,
		)
	}

	return data, nil, nil
}