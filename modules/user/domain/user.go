package domain

import (
	"errors"
	"regexp"
	"time"
)

// User representa un usuario en el sistema
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // No se expone en JSON
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewUser crea una nueva instancia de User
func NewUser(username, email, password, firstName, lastName string) (*User, error) {
	now := time.Now().UTC()
	user := &User{
		Username:  username,
		Email:     email,
		Password:  password, // Debe venir hasheado
		FirstName: firstName,
		LastName:  lastName,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if !user.IsValid() {
		return nil, errors.New("el usuario no es válido")
	}

	return user, nil
}

// IsValid valida que el usuario tenga los campos requeridos
func (u *User) IsValid() bool {
	if u.Username == "" || len(u.Username) < 3 {
		return false
	}

	if !isValidEmail(u.Email) {
		return false
	}

	if u.Password == "" || len(u.Password) < 6 {
		return false
	}

	if u.FirstName == "" || u.LastName == "" {
		return false
	}

	return true
}

// isValidEmail valida el formato del email
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Update actualiza los campos del usuario
func (u *User) Update(firstName, lastName string) {
	if firstName != "" {
		u.FirstName = firstName
	}
	if lastName != "" {
		u.LastName = lastName
	}
	u.UpdatedAt = time.Now().UTC()
}

// Deactivate desactiva el usuario
func (u *User) Deactivate() {
	u.Active = false
	u.UpdatedAt = time.Now().UTC()
}

// Activate activa el usuario
func (u *User) Activate() {
	u.Active = true
	u.UpdatedAt = time.Now().UTC()
}
