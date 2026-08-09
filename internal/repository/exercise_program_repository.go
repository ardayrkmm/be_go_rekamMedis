package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type ExerciseProgramRepository interface {
	FindAll(offset, limit int) ([]models.ExerciseProgram, int64, error)
	FindByID(id string) (*models.ExerciseProgram, error)
	Create(program *models.ExerciseProgram) error
	Update(program *models.ExerciseProgram) error
	Delete(id string) error
}

type exerciseProgramRepository struct {
	db *gorm.DB
}

func NewExerciseProgramRepository(db *gorm.DB) ExerciseProgramRepository {
	return &exerciseProgramRepository{db}
}

func (r *exerciseProgramRepository) FindAll(offset, limit int) ([]models.ExerciseProgram, int64, error) {
	var programs []models.ExerciseProgram
	var total int64
	err := r.db.Model(&models.ExerciseProgram{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Preload("TherapySession").Offset(offset).Limit(limit).Find(&programs).Error
	return programs, total, err
}

func (r *exerciseProgramRepository) FindByID(id string) (*models.ExerciseProgram, error) {
	var program models.ExerciseProgram
	err := r.db.Preload("TherapySession").First(&program, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &program, nil
}

func (r *exerciseProgramRepository) Create(program *models.ExerciseProgram) error {
	return r.db.Create(program).Error
}

func (r *exerciseProgramRepository) Update(program *models.ExerciseProgram) error {
	return r.db.Save(program).Error
}

func (r *exerciseProgramRepository) Delete(id string) error {
	return r.db.Delete(&models.ExerciseProgram{}, id).Error
}
