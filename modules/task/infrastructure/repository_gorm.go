package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/domain"
	"gorm.io/gorm"
)

// GormTaskModel es el modelo de GORM para la tabla tasks (PostgreSQL)
type GormTaskModel struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`    // SERIAL en PostgreSQL
	Title       string    `gorm:"not null;size:255" json:"title"`        // VARCHAR(255)
	Description string    `gorm:"not null;type:text" json:"description"` // TEXT
	Completed   bool      `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (GormTaskModel) TableName() string {
	return "tasks"
}

// ToDomain convierte el modelo GORM a entidad de dominio
func (g *GormTaskModel) ToDomain() *domain.Task {
	return &domain.Task{
		ID:          g.ID,
		Title:       g.Title,
		Description: g.Description,
		Completed:   g.Completed,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

// FromDomain convierte entidad de dominio a modelo GORM
func (g *GormTaskModel) FromDomain(task *domain.Task) {
	g.ID = task.ID
	g.Title = task.Title
	g.Description = task.Description
	g.Completed = task.Completed
	g.CreatedAt = task.CreatedAt
	g.UpdatedAt = task.UpdatedAt
}

// GormTaskRepository implementa TaskRepository usando GORM
type GormTaskRepository struct {
	db *gorm.DB
}

// NewGormTaskRepository crea una nueva instancia del repositorio GORM
func NewGormTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &GormTaskRepository{
		db: db,
	}
}

// Create inserta una nueva tarea con Gorm
func (r *GormTaskRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {

	gormTask := &GormTaskModel{} // Inicializa el modelo GORM
	gormTask.FromDomain(task)    // Convierte la entidad de dominio a modelo GORM

	if err := r.db.WithContext(ctx).Create(gormTask).Error; err != nil {
		return nil, fmt.Errorf("error creando tarea con GORM: %w", err)
	}

	return gormTask.ToDomain(), nil
}

func (r *GormTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	var gormTasks []GormTaskModel

	if err := r.db.WithContext(ctx).Find(&gormTasks).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo todas las tareas con GORM: %w", err)
	}

	tasks := make([]*domain.Task, len(gormTasks))
	for i, gormTask := range gormTasks {
		tasks[i] = gormTask.ToDomain()
	}

	return tasks, nil

}

// GetByID obtiene una tarea por su ID usando GORM
func (r *GormTaskRepository) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	var gormTask GormTaskModel

	if err := r.db.WithContext(ctx).First(&gormTask, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tarea con ID %d no encontrada", id)
		}
		return nil, fmt.Errorf("error obteniendo tarea con GORM: %w", err)
	}

	return gormTask.ToDomain(), nil
}

func (r *GormTaskRepository) Update(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	gormTask := &GormTaskModel{}
	gormTask.FromDomain(task)

	result := r.db.WithContext(ctx).Model(&GormTaskModel{}).Where("id = ?", task.ID).Updates(gormTask)
	if result.Error != nil {
		return nil, fmt.Errorf("error actualizando tarea con GORM: %w", &result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("tarea con id %d no encontrada", task.ID)
	}

	var updatedTask GormTaskModel
	if err := r.db.WithContext(ctx).First(&updatedTask, task.ID).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo tarea actualizada: %w", err)
	}

	return updatedTask.ToDomain(), nil
}

func (r *GormTaskRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&GormTaskModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("error eliminando tarea con GORM: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("tarea con id %d no encontrada", id)
	}
	return nil
}

func (r *GormTaskRepository) GetByStatus(ctx context.Context, completed bool) ([]*domain.Task, error) {
	var gormTasks []GormTaskModel
	if err := r.db.WithContext(ctx).Where("completed = ?", completed).Order("created_at DESC").Find(&gormTasks).Error; err != nil {
		return nil, fmt.Errorf("error obteniendo tareas por estado con GORM: %w", err)
	}

	tasks := make([]*domain.Task, len(gormTasks))
	for i, gormtask := range gormTasks {
		tasks[i] = gormtask.ToDomain()
	}
	return tasks, nil
}
