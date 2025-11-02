package infrastructure

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/stretchr/testify/require"

    "github.com/YerkoTenorio/api-go-hexagonal/modules/task/domain"
    "github.com/YerkoTenorio/api-go-hexagonal/shared/config"
    "github.com/YerkoTenorio/api-go-hexagonal/shared/database"
)

func newTestSQLiteDB(t *testing.T) (*database.SQLiteDB, string) {
    t.Helper()
    // Crear archivo temporal para la BD de prueba
    dbPath := filepath.Join(os.TempDir(), fmt.Sprintf("tasks_test_%d.db", time.Now().UnixNano()))
    cfg := &config.Config{Database: config.DatabaseConfig{Path: dbPath}}
    sqliteDB, err := database.NewSQLiteDB(cfg)
    require.NoError(t, err)
    t.Cleanup(func() {
        _ = sqliteDB.Close()
        _ = os.Remove(dbPath)
    })
    return sqliteDB, dbPath
}

func TestSQLiteTaskRepository_CreateAndGetByID(t *testing.T) {
    ctx := context.Background()
    sqliteDB, _ := newTestSQLiteDB(t)
    repo := NewSQLiteTaskRepository(sqliteDB)

    // Crear tarea
    input := &domain.Task{Title: "Comprar pan", Description: "Ir a la panadería", Completed: false}
    created, err := repo.Create(ctx, input)
    require.NoError(t, err)
    require.True(t, created.ID > 0)
    require.Equal(t, input.Title, created.Title)
    require.Equal(t, input.Description, created.Description)
    require.False(t, created.Completed)
    require.False(t, created.CreatedAt.IsZero())
    require.False(t, created.UpdatedAt.IsZero())

    // Obtener por ID
    fetched, err := repo.GetByID(ctx, created.ID)
    require.NoError(t, err)
    require.Equal(t, created.ID, fetched.ID)
    require.Equal(t, created.Title, fetched.Title)
    require.Equal(t, created.Description, fetched.Description)
    require.Equal(t, created.Completed, fetched.Completed)
}

func TestSQLiteTaskRepository_GetAllAndGetByStatus(t *testing.T) {
    ctx := context.Background()
    sqliteDB, _ := newTestSQLiteDB(t)
    repo := NewSQLiteTaskRepository(sqliteDB)

    // Sembrar varias tareas
    _, err := repo.Create(ctx, &domain.Task{Title: "T1", Description: "D1", Completed: false})
    require.NoError(t, err)
    _, err = repo.Create(ctx, &domain.Task{Title: "T2", Description: "D2", Completed: true})
    require.NoError(t, err)
    _, err = repo.Create(ctx, &domain.Task{Title: "T3", Description: "D3", Completed: false})
    require.NoError(t, err)

    all, err := repo.GetAll(ctx)
    require.NoError(t, err)
    require.Len(t, all, 3)

    completed, err := repo.GetByStatus(ctx, true)
    require.NoError(t, err)
    require.Len(t, completed, 1)

    pending, err := repo.GetByStatus(ctx, false)
    require.NoError(t, err)
    require.Len(t, pending, 2)
}

func TestSQLiteTaskRepository_Update(t *testing.T) {
    ctx := context.Background()
    sqliteDB, _ := newTestSQLiteDB(t)
    repo := NewSQLiteTaskRepository(sqliteDB)

    created, err := repo.Create(ctx, &domain.Task{Title: "Inicial", Description: "Desc", Completed: false})
    require.NoError(t, err)

    // Modificar campos
    created.Title = "Actualizado"
    created.Description = "Desc Actualizada"
    created.Completed = true
    updated, err := repo.Update(ctx, created)
    require.NoError(t, err)
    require.Equal(t, created.ID, updated.ID)
    require.Equal(t, "Actualizado", updated.Title)
    require.Equal(t, "Desc Actualizada", updated.Description)
    require.True(t, updated.Completed)
    require.True(t, updated.UpdatedAt.After(updated.CreatedAt) || updated.UpdatedAt.Equal(updated.CreatedAt))

    // Verificar persistencia
    fetched, err := repo.GetByID(ctx, created.ID)
    require.NoError(t, err)
    require.Equal(t, "Actualizado", fetched.Title)
    require.True(t, fetched.Completed)
}

func TestSQLiteTaskRepository_Delete(t *testing.T) {
    ctx := context.Background()
    sqliteDB, _ := newTestSQLiteDB(t)
    repo := NewSQLiteTaskRepository(sqliteDB)

    created, err := repo.Create(ctx, &domain.Task{Title: "Borrar", Description: "Desc", Completed: false})
    require.NoError(t, err)

    err = repo.Delete(ctx, created.ID)
    require.NoError(t, err)

    _, err = repo.GetByID(ctx, created.ID)
    require.Error(t, err)
}