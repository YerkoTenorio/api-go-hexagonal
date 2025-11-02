package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/YerkoTenorio/api-go-hexagonal/shared/config"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

// SQLiteDB encapsula la conexion a SQLite

type SQLiteDB struct {
	DB *sql.DB
}

// NewSQLiteDB crea una nueva conexion a SQLite usando la configuracion

func NewSQLiteDB(cfg *config.Config) (*SQLiteDB, error) {
	var db *sql.DB
	var err error

	// Si hay URL configurada, conectar a libSQL (Turso) remoto
	if cfg.Database.URL != "" {
		dsn := cfg.Database.URL
		if cfg.Database.AuthToken != "" {
			sep := "?"
			if strings.Contains(dsn, "?") {
				sep = "&"
			}
			dsn = fmt.Sprintf("%s%sauthToken=%s", dsn, sep, url.QueryEscape(cfg.Database.AuthToken))
		}

		// Log: mostrar URL sin token
		fmt.Println("[DB] Conectando a libSQL remoto:", cfg.Database.URL)

		db, err = sql.Open("libsql", dsn)
		if err != nil {
			return nil, fmt.Errorf("error abriendo base de datos remota (libSQL): %w", err)
		}
	} else {
		// Conexión local SQLite usando modernc.org/sqlite
		dbPath := cfg.Database.Path

		// Crear el directorio si no existe
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("error creando el directorio de BD: %w", err)
		}

		// Abrir conexión a SQLite local
		// Log: mostrar ruta local
		fmt.Println("[DB] Conectando a SQLite local:", dbPath)

		db, err = sql.Open("sqlite", dbPath)
		if err != nil {
			return nil, fmt.Errorf("error abriendo base de datos local: %w", err)
		}
	}

	// Verificar que la conexion funciona
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verificando conexion a BD: %w", err)
	}

	// Log de versión para confirmar conexión
	var version string
	if err := db.QueryRow("SELECT sqlite_version()").Scan(&version); err == nil {
		fmt.Println("[DB] Conexión OK. sqlite_version():", version)
	} else {
		fmt.Println("[DB] Conexión OK. (No se pudo leer sqlite_version)")
	}

	sqliteDB := &SQLiteDB{DB: db}

	//Crear las tablas si no existen
	if err := sqliteDB.createTables(); err != nil {
		return nil, fmt.Errorf("error creando tablas: %w", err)
	}
	return sqliteDB, nil
}

func (s *SQLiteDB) createTables() error {
	createTasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
	   id INTEGER PRIMARY KEY AUTOINCREMENT,
	   title TEXT NOT NULL,
	   description TEXT NOT NULL,
	   completed BOOLEAN NOT NULL DEFAULT FALSE,
	   created_at DATETIME NOT NULL,
	   updated_at DATETIME NOT NULL 
	   );`

	_, err := s.DB.Exec(createTasksTable)
	if err != nil {
		return fmt.Errorf("error creando tabla tasks: %w", err)
	}

	return nil
}

// Close cierra la conexion a la base de datos
func (s *SQLiteDB) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}

// GetDB retorna la instancia de la base de datos
func (s *SQLiteDB) GetDB() *sql.DB {
	return s.DB
}
