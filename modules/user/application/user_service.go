package application

import (
	"context"
	"fmt"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/user/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserService maneja los casos de uso relacionados con usuarios
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService crea una nueva instancia de UserService

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser crea un nuevo usuario
func (s *UserService) CreateUser(ctx context.Context, username, email, password, firstname, lastname string) (*domain.User, error) {
	//Validaciones basicas
	if username == "" {
		return nil, fmt.Errorf("el username es requerido")
	}
	if email == "" {
		return nil, fmt.Errorf("el email es requerido")
	}
	if password == "" {
		return nil, fmt.Errorf("el password es requerido")
	}
	if firstname == "" {
		return nil, fmt.Errorf("el nombre es requerido")
	}
	if lastname == "" {
		return nil, fmt.Errorf("el apellido es requerido")
	}

	//Verificar si el username ya existe
	existingUserByUsername, _ := s.userRepo.GetByUsername(ctx, username)
	if existingUserByUsername != nil {
		return nil, fmt.Errorf("el username ya esta en uso")
	}

	// Verificar si el email ya existe
	existingUserByEmail, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUserByEmail != nil {
		return nil, fmt.Errorf("el email ya esta en uso")
	}

	// Hashear la contrasena
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al hashear la contrasena: %w", err)
	}

	//Crear nuevo usuario
	user, err := domain.NewUser(username, email, string(hashedPassword), firstname, lastname)
	if err != nil {
		return nil, fmt.Errorf("error al crear el usuario: %w", err)
	}

	// pesistir el usuario
	return s.userRepo.Create(ctx, user)
}

// AuthenticateUser autentica un usuario
func (s *UserService) AuthenticateUser(ctx context.Context, username, password string) (*domain.User, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username y contrasena son requeridos")
	}

	//Buscar usuario por username
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Verificar si el usuario esta activo

	if !user.Active {
		return nil, fmt.Errorf("usuario inactivo")
	}

	// Verificar contrasena
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("contrasena incorrecta")
	}

	return user, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("el ID del usuario no puede ser cero")
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el usuario con ID %d: %w", id, err)
	}

	return user, nil
}

// GetAllUsers obtiene todos los usuarios
func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron obtener los usuarios: %w", err)
	}

	return users, nil
}

// UpdateUser actualiza un usuario existente
func (s *UserService) UpdateUser(ctx context.Context, id int, firstName, lastName string) (*domain.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("el ID del usuario es requerido")
	}

	// Obtener el usuario existente
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("no se pudo encontrar el usuario con ID %d: %w", id, err)
	}

	// Actualizar campos
	user.Update(firstName, lastName)

	// Validar el usuario actualizado
	if !user.IsValid() {
		return nil, fmt.Errorf("el usuario actualizado no es válido")
	}

	// Persistir los cambios
	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("no se pudo actualizar el usuario: %w", err)
	}

	return updatedUser, nil
}

// DeactivateUser desactiva un usuario
func (s *UserService) DeactivateUser(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("el ID del usuario es requerido")
	}

	// Obtener el usuario existente
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("no se pudo encontrar el usuario con ID %d: %w", id, err)
	}

	// Desactivar usuario
	user.Deactivate()

	// Persistir los cambios
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("no se pudo desactivar el usuario: %w", err)
	}

	return nil
}

// ActivateUser activa un usuario
func (s *UserService) ActivateUser(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("el ID del usuario es requerido")
	}

	// Obtener el usuario existente
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("no se pudo encontrar el usuario con ID %d: %w", id, err)
	}

	// Activar usuario
	user.Activate()

	// Persistir los cambios
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("no se pudo activar el usuario: %w", err)
	}

	return nil
}

// DeleteUser elimina un usuario por su ID
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("el ID del usuario es requerido")
	}

	// Verificar que el usuario existe antes de eliminarlo
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("no se pudo encontrar el usuario con ID %d: %w", id, err)
	}

	// Eliminar el usuario
	err = s.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el usuario con ID %d: %w", id, err)
	}

	return nil
}

// GetActiveUsers obtiene solo los usuarios activos
func (s *UserService) GetActiveUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.GetActiveUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron obtener los usuarios activos: %w", err)
	}

	return users, nil
}
