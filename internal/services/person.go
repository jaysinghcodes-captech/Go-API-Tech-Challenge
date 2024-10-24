package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jaysinghcodes-captech/Go-API-Tech-Challenge/internal/models"
)

type PersonService struct {
	DB *sql.DB
}

func NewPersonService(db *sql.DB) *PersonService {
	return &PersonService{
		DB: db,
	}
}

func (p *PersonService) ListPersons(ctx context.Context) ([]models.Person, error) {
	rows, err := p.DB.Query("SELECT id, first_name, last_name, type, age FROM person ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("[in services.ListPersons] failed to get persons: %w", err)
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age)
		if err != nil {
			return nil, fmt.Errorf("[in services.ListPersons] failed to scan person from row: %w", err)
		}

		// Fetch courses for the current person
		courseIDs, err := p.getCoursesForPerson(person.ID)
		if err != nil {
			return nil, err
		}
		person.Courses = courseIDs

		persons = append(persons, person)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("[in services.ListPersons] failed to scan persons: %w", err)
	}

	return persons, nil
}

func (p *PersonService) GetPersonByFirstName(ctx context.Context, firstName string) (models.Person, error) {
	var person models.Person
	err := p.DB.QueryRow("SELECT id, first_name, last_name, type, age FROM person WHERE first_name = $1", firstName).Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Person{}, fmt.Errorf("[in services.GetPersonByFirstName] person with first name %s not found: %w", firstName, err)
		}
		return models.Person{}, fmt.Errorf("[in services.GetPersonByFirstName] failed to get person with first name %s: %w", firstName, err)
	}

	// Fetch courses for the person
	courseIDs, err := p.getCoursesForPerson(person.ID)
	if err != nil {
		return models.Person{}, err
	}
	person.Courses = courseIDs

	return person, nil
}

func (p *PersonService) UpdatePerson(ctx context.Context, firstName string, updatedPerson models.Person) (models.Person, error) {
	// Validate the updated person object
	if updatedPerson.FirstName == "" || updatedPerson.LastName == "" || updatedPerson.Type == "" || updatedPerson.Age <= 0 {
		return models.Person{}, fmt.Errorf("[in services.UpdatePerson] invalid person data")
	}

	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Person{}, fmt.Errorf("[in services.UpdatePerson] failed to begin transaction: %w", err)
	}

	// Fetch the person_id using the old firstName
	var personID int
	err = tx.QueryRowContext(ctx, "SELECT id FROM person WHERE first_name = $1", firstName).Scan(&personID)
	if err != nil {
		tx.Rollback()
		return models.Person{}, fmt.Errorf("[in services.UpdatePerson] failed to fetch person with first name %s: %w", firstName, err)
	}

	// Update the person details
	_, err = tx.ExecContext(ctx, "UPDATE person SET first_name = $1, last_name = $2, type = $3, age = $4 WHERE id = $5",
		updatedPerson.FirstName, updatedPerson.LastName, updatedPerson.Type, updatedPerson.Age, personID)
	if err != nil {
		tx.Rollback()
		return models.Person{}, fmt.Errorf("[in services.UpdatePerson] failed to update person with id %d: %w", personID, err)
	}

	// Clear existing courses
	_, err = tx.ExecContext(ctx, "DELETE FROM person_course WHERE person_id = $1", personID)
	if err != nil {
		tx.Rollback()
		return models.Person{}, fmt.Errorf("[in services.UpdatePerson] failed to clear existing courses for person with id %d: %w", personID, err)
	}

	// Associate new courses
	for _, courseID := range updatedPerson.Courses {
		_, err = tx.ExecContext(ctx, "INSERT INTO person_course (person_id, course_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", personID, courseID)
		if err != nil {
			tx.Rollback()
			return models.Person{}, fmt.Errorf("[in services.UpdatePerson] failed to associate new courses with person id %d: %w", personID, err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return models.Person{}, fmt.Errorf("[in services.UpdatePerson] failed to commit transaction: %w", err)
	}

	return updatedPerson, nil
}

func (p *PersonService) CreatePerson(ctx context.Context, person models.Person) (models.Person, error) {
	var newID int
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return models.Person{}, fmt.Errorf("[in services.CreatePerson] failed to begin transaction: %w", err)
	}

	err = tx.QueryRowContext(ctx, "INSERT INTO person (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id", person.FirstName, person.LastName, person.Type, person.Age).Scan(&newID)
	if err != nil {
		tx.Rollback()
		return models.Person{}, fmt.Errorf("[in services.CreatePerson] failed to create person: %w", err)
	}

	for _, courseID := range person.Courses {
		_, err = tx.ExecContext(ctx, "INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)", newID, courseID)
		if err != nil {
			tx.Rollback()
			return models.Person{}, fmt.Errorf("[in services.CreatePerson] failed to associate course with person: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return models.Person{}, fmt.Errorf("[in services.CreatePerson] failed to commit transaction: %w", err)
	}

	createdPerson := models.Person{
		ID:        newID,
		FirstName: person.FirstName,
		LastName:  person.LastName,
		Type:      person.Type,
		Age:       person.Age,
		Courses:   person.Courses,
	}

	return createdPerson, nil
}

func (p *PersonService) DeletePerson(ctx context.Context, firstName string) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("[in services.DeletePerson] failed to begin transaction: %w", err)
	}

	// Fetch the person_id using the firstName
	var personID int
	err = tx.QueryRowContext(ctx, "SELECT id FROM person WHERE first_name = $1", firstName).Scan(&personID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[in services.DeletePerson] failed to fetch person with first name %s: %w", firstName, err)
	}

	// Clear associated courses first
	_, err = tx.ExecContext(ctx, "DELETE FROM person_course WHERE person_id = $1", personID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[in services.DeletePerson] failed to clear associated courses for person with id %d: %w", personID, err)
	}

	// Then delete the person
	result, err := tx.ExecContext(ctx, "DELETE FROM person WHERE id = $1", personID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[in services.DeletePerson] failed to delete person with id %d: %w", personID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("[in services.DeletePerson] failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("[in services.DeletePerson] person with id %d not found", personID)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("[in services.DeletePerson] failed to commit transaction: %w", err)
	}

	return nil
}

// Helper method to retrieve courses for a person
func (p *PersonService) getCoursesForPerson(personID int) ([]int, error) {
	rows, err := p.DB.Query("SELECT course_id FROM person_course WHERE person_id = $1", personID)
	if err != nil {
		return nil, fmt.Errorf("[in services.getCoursesForPerson] failed to get courses: %w", err)
	}
	defer rows.Close()

	var courseIDs []int
	for rows.Next() {
		var courseID int
		if err := rows.Scan(&courseID); err != nil {
			return nil, fmt.Errorf("[in services.getCoursesForPerson] failed to scan course ID: %w", err)
		}
		courseIDs = append(courseIDs, courseID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("[in services.getCoursesForPerson] failed to scan course IDs: %w", err)
	}

	return courseIDs, nil
}
