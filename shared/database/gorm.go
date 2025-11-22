package database

import (
	"fmt"

	"github.com/YerkoTenorio/api-go-hexagonal/shared/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormDB encapsula la conexión a GORM con PostgreSQL
type GormDB struct {
	DB *gorm.DB
}

// NewGormDB crea una nueva conexión a GORM usando PostgreSQL
func NewGormDB(cfg *config.Config) (*GormDB, error) {
	// Configurar GORM en modo silencioso para logs limpios
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	// Abrir conexión con PostgreSQL
	fmt.Println("[GORM] Conectando a PostgreSQL:", cfg.Database.URL)

	db, err := gorm.Open(postgres.Open(cfg.Database.URL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("error abriendo base de datos PostgreSQL con GORM: %w", err)
	}

	// Obtener la conexión subyacente para verificar
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo conexión SQL: %w", err)
	}

	// Verificar que la conexión funciona
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error verificando conexión a BD: %w", err)
	}

	fmt.Println("[GORM] Conexión OK")

	gormDB := &GormDB{DB: db}

	// Auto-migrar las tablas
	if err := gormDB.autoMigrate(); err != nil {
		return nil, fmt.Errorf("error migrando tablas: %w", err)
	}

	return gormDB, nil
}

func (g *GormDB) autoMigrate() error {
	fmt.Println("[GORM] Ejecutando auto-migración...")

	// Crear la tabla tasks para PostgreSQL
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	if err := g.DB.Exec(createTableSQL).Error; err != nil {
		return fmt.Errorf("error creando tabla tasks con GORM: %w", err)
	}

	fmt.Println("[GORM] Auto-migración completada")
	return nil
}

// Close cierra la conexión a la base de datos
func (g *GormDB) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB retorna la instancia de GORM
func (g *GormDB) GetDB() *gorm.DB {
	return g.DB
}
