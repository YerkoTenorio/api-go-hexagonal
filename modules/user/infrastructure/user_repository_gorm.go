package infrastructure

import (
	"context"
	"fmt"

	"github.com/YerkoTenorio/api-go-hexagonal/modules/user/domain"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		db: db,
	}
}

// Create crea un nuevo usuario en la base de datos
func (r *GormUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	gormUser := r.toGormModel(user)

	result := r.db.WithContext(ctx).Create(gormUser)
	if result.Error != nil {
		return nil, fmt.Errorf("error al crear usuario: %w", result.Error)
	}

	return r.toDomainModel(gormUser), nil
}

// GetByID obtiene un usuario por su ID

func (r *GormUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var gormUser GormUserModel

	result := r.db.WithContext(ctx).First(&gormUser, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("usuario con ID %d no encontrado", id)
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", result.Error)
	}

	return r.toDomainModel(&gormUser), nil
}

// GetByUsername obtiene un usuario por su username

func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var gormUser GormUserModel

	result := r.db.WithContext(ctx).First(&gormUser, username)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("usuario con username %s no encontrado", username)
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", result.Error)
	}

	return r.toDomainModel(&gormUser), nil

}

// GetByEmail obtiene un usuario por su email
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var gormUser GormUserModel

	result := r.db.WithContext(ctx).Where("email = ?", email).First(&gormUser)
	if result.Error != nil {
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("usuario con email %s no encontrado", email)
			}
			return nil, fmt.Errorf("error al obtener el usuario: %w", result.Error)
		}
	}

	return r.toDomainModel(&gormUser), nil
}

// GetAll obtiene todos los usuarios
func (r *GormUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var gormUsers []GormUserModel

	result := r.db.WithContext(ctx).Find(&gormUsers)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", result.Error)
	}

	users := make([]*domain.User, len(gormUsers))
	for i, gormUser := range gormUsers {
		users[i] = r.toDomainModel(&gormUser)
	}
	return users, nil
}

// GetActiveUsers obtiene solo los usuarios activos

func (r *GormUserRepository) GetActiveUsers(ctx context.Context) ([]*domain.User, error) {
	var gormUsers []GormUserModel

	result := r.db.WithContext(ctx).Where("active = ?", true).Find(&gormUsers)
	if result.Error != nil {
		return nil, fmt.Errorf("error al obtener usuarios activos: %w", result.Error)
	}

	users := make([]*domain.User, len(gormUsers))
	for i, gormUser := range gormUsers {
		users[i] = r.toDomainModel(&gormUser)
	}

	return users, nil
}

// Update actualiza un usuario existente

func (r *GormUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	gormUser := r.toGormModel(user)

	result := r.db.WithContext(ctx).Save(gormUser)
	if result.Error != nil {
		return nil, fmt.Errorf("error al actualizar usuario: %w", result.Error)
	}

	return r.toDomainModel(gormUser), nil
}

// Delete elimina un usuario por su ID

func (r *GormUserRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&GormUserModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("error al eliminar usuario: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("usuario con ID %d no encontrado", id)
	}

	return nil
}

// toGormModel convierte un domain.User a GormUserModel
func (r *GormUserRepository) toGormModel(user *domain.User) *GormUserModel {
	return &GormUserModel{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// toDomainModel convierte un GormUserModel a domain.User
func (r *GormUserRepository) toDomainModel(gormUser *GormUserModel) *domain.User {
	return &domain.User{
		ID:        gormUser.ID,
		Username:  gormUser.Username,
		Email:     gormUser.Email,
		Password:  gormUser.Password,
		FirstName: gormUser.FirstName,
		LastName:  gormUser.LastName,
		Active:    gormUser.Active,
		CreatedAt: gormUser.CreatedAt,
		UpdatedAt: gormUser.UpdatedAt,
	}

}
