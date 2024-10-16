package services

import (
	"database/sql"
	"context"
	"fmt"
	"go-api-tech-challenge/internal/models"
)

type CourseService struct {
	DB *sql.DB
}

func NewCourseService(db *sql.DB) *CourseService {
	return &CourseService{
		DB: db,
	}
}

func (c *CourseService) ListCourses(ctx context.Context) ([]models.Course, error) {
	rows, err := c.DB.Query("SELECT * FROM courses ORDER BY id")
	if err != nil {
		return []models.Course{}, fmt.Errorf("[in services.ListCourses] failed to get courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Name)
		if err != nil {
			return []models.Course{}, fmt.Errorf("[in services.ListCourses] failed to scan course from row: %w", err)
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return []models.Course{}, fmt.Errorf("[in services.ListCourses] failed to scan courses: %w", err)
	}

	return courses, nil
}

func (c *CourseService) GetCourseById(ctx context.Context, id int) (models.Course, error) {
	var course models.Course
	err := c.DB.QueryRow("SELECT * FROM courses WHERE id = $1", id).Scan(&course.ID, &course.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Course{}, fmt.Errorf("[in services.GetCourseById] course with id %d not found: %w", id, err)
		}
		return models.Course{}, fmt.Errorf("[in services.GetCourseById] failed to get course with id %d: %w", id, err)
	}

	return course, nil
}

func (c *CourseService) CreateCourse(ctx context.Context, courseName string) (models.Course, error) {
	var newID int
	err := c.DB.QueryRow("INSERT INTO courses (name) VALUES ($1) RETURNING id", courseName).Scan(&newID)
	if err != nil {
		return models.Course{}, fmt.Errorf("[in services.CreateCourse] failed to create course: %w", err)
	}

	return models.Course{ID: newID, Name: courseName}, nil
}

func (c *CourseService) UpdateCourse(ctx context.Context, courseID int, newCourseName string) (models.Course, error) {
	result, err := c.DB.Exec("UPDATE courses SET name = $1 WHERE id = $2", newCourseName, courseID)
	if err != nil {
		return models.Course{}, fmt.Errorf("[in services.UpdateCourse] failed to update course with id %d: %w", courseID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Course{}, fmt.Errorf("[in services.UpdateCourse] failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.Course{}, fmt.Errorf("[in services.UpdateCourse] course with id %d not found", courseID)
	}

	return models.Course{ID: courseID, Name: newCourseName}, nil
}

func (c *CourseService) DeleteCourse(ctx context.Context, id int) error {
	result, err := c.DB.Exec("DELETE FROM courses WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("[in services.DeleteCourse] failed to delete course with id %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("[in services.DeleteCourse] failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("[in services.DeleteCourse] course with id %d not found", id)
	}

	return nil
}

