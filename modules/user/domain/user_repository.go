package domain

import "context"

// UserRepository define el puerto para persistencia de usuarios
type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)             // Crea un nuevo usuario
	GetByID(ctx context.Context, id int) (*User, error)                // Obtiene un usuario por ID
	GetByUsername(ctx context.Context, username string) (*User, error) // Obtiene un usuario por nombre de usuario
	GetByEmail(ctx context.Context, email string) (*User, error)       // Obtiene un usuario por email
	GetAll(ctx context.Context) ([]*User, error)                       // Obtiene todos los usuarios
	Update(ctx context.Context, user *User) (*User, error)             // Actualiza un usuario
	Delete(ctx context.Context, id int) error                          // Elimina un usuario
	GetActiveUsers(ctx context.Context) ([]*User, error)               // Obtiene todos los usuarios activos
}
