package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/task/domain"
	"github.com/YerkoTenorio/api-go-hexagonal/shared/database"
)

// SQLiteTaskRepository implementa TaskRepository usando SQLite
type SQLiteTaskRepository struct {
	db *database.SQLiteDB
}

// NewSQLiteTaskRepository crea una nueva instancia del repositorio
func NewSQLiteTaskRepository(db *database.SQLiteDB) domain.TaskRepository {
	return &SQLiteTaskRepository{
		db: db,
	}
}

// Create inserta una nueva tarea en la base de datos
func (r *SQLiteTaskRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	query := `INSERT INTO tasks (title, description, completed, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now().UTC()
	result, err := r.db.GetDB().ExecContext(ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("error insertando tarea: %w", err)
	}
	// Obtener el ID generado

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ID de tarea insertada: %w", err)
	}
	// Actualizar la tarea con el ID y timestamps
	task.ID = int(id)
	task.CreatedAt = now
	task.UpdatedAt = now

	return task, nil
}

// GetByID obtiene una tarea por su ID
func (r *SQLiteTaskRepository) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = ?`

	row := r.db.GetDB().QueryRowContext(ctx, query, id)

	task := &domain.Task{}
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tarea con ID %d no encontrada", id)
		}
		return nil, fmt.Errorf("error obteniendo tarea: %w", err)
	}
	return task, nil
}

// GetAll obtiene todas las tareas
func (r *SQLiteTaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	// Definir la consulta SQL
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks`
	// Obtener todas las filas
	rows, err := r.db.GetDB().QueryContext(ctx, query)
	// Manejar el error de la consulta
	if err != nil {
		return nil, fmt.Errorf("error obteniendo todas las tareas: %w", err)
	}

	defer rows.Close()         // Cerrar las filas después de usarlas
	var tasks []*domain.Task // Slice para almacenar las tareas

	for rows.Next() { // Iterar sobre cada fila
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando tarea: %w", err) // Manejar error de escaneo
		}
		tasks = append(tasks, task) // Añadir la tarea al slice
	}

	if err = rows.Err(); err != nil { // Manejar error de iteración
		return nil, fmt.Errorf("error iterando sobre filas: %w", err)
	}

	return tasks, nil

}

// Update actualiza una tarea existente en la base de datos
func (r *SQLiteTaskRepository) Update(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	query := `UPDATE tasks SET title = ?, description = ?, completed = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	result, err := r.db.GetDB().ExecContext(ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		now,
		task.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("error actualizando tarea: %w", err)
	}
	// Verificar si se actualizo alguna fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error verificando filas afectadas: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("tarea con ID %d no encontrada", task.ID)
	}
	// Actualizar timestamp
	task.UpdatedAt = now

	return task, nil

}

// Delete elimina una tarea de la base de datos
func (r *SQLiteTaskRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := r.db.GetDB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error eliminando tarea: %w", err)
	}

	// Verificar que se elimino al menos una fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error verificando eliminacion: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("tarea con ID %d no encontrada", id)
	}

	return nil
}

// GetByStatus obtiene tareas por su estado (completadas o no)
func (r *SQLiteTaskRepository) GetByStatus(ctx context.Context, completed bool) ([]*domain.Task, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE completed = ? ORDER BY created_at DESC`

	rows, err := r.db.GetDB().QueryContext(ctx, query, completed)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo tareas por estado: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando tarea: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando sobre filas: %w", err)
	}
	return tasks, nil
}
