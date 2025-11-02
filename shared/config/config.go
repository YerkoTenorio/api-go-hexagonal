package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    App      AppConfig
    Log      LogConfig
}

// DatabaseConfig configuración de la base de datos
type DatabaseConfig struct {
    Path      string
    URL       string
    AuthToken string
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Port string
	Host string
}

// AppConfig configuración general de la aplicación
type AppConfig struct {
	Environment string
	Name        string
}

// LogConfig configuración de logs
type LogConfig struct {
	Level string
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() (*Config, error) {
	// Cargar el archivo .env si existe
	if err := godotenv.Load(); err != nil {
		// No es un error crítico si no existe el archivo .env
		fmt.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

    config := &Config{
        Database: DatabaseConfig{
            Path:      getEnv("DB_PATH", "./data/tasks.db"),
            URL:       getEnv("DB_URL", ""),
            AuthToken: getEnv("DB_AUTH_TOKEN", ""),
        },
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
            Host: getEnv("SERVER_HOST", "localhost"),
        },
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Name:        getEnv("APP_NAME", "Task Manager API"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	// Validar configuración crítica
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuración inválida: %w", err)
	}

	return config, nil
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt obtiene una variable de entorno como entero
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool obtiene una variable de entorno como booleano
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// validate valida que la configuración sea correcta
func (c *Config) validate() error {
    // Validar que exista configuración de base de datos: local (Path) o remota (URL)
    if c.Database.Path == "" && c.Database.URL == "" {
        return fmt.Errorf("Debe especificarse DB_PATH (local) o DB_URL (remota)")
    }

	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT no puede estar vacío")
	}

	// Validar que el puerto sea un número válido
	if _, err := strconv.Atoi(c.Server.Port); err != nil {
		return fmt.Errorf("SERVER_PORT debe ser un número válido: %s", c.Server.Port)
	}

    return nil
}

// GetServerAddress retorna la dirección completa del servidor
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// IsDevelopment verifica si estamos en modo desarrollo
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction verifica si estamos en modo producción
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// ServerAddress retorna la dirección del servidor
func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
