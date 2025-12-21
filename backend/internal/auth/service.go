package auth

import (
	"errors"
	"os"
	"time"

	// <--- Importante: v5
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// RegisterRequest define qué datos necesitamos del Frontend (DTO)
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (s *Service) RegisterUser(req RegisterRequest) error {
	// 1. Validar si el usuario ya existe
	existingUser, _ := s.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return errors.New("el correo electrónico ya está registrado")
	}

	// 2. Hashear la contraseña (Seguridad Crítica)
	// Costo 10 es un buen balance entre seguridad y velocidad
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return err
	}

	// 3. Crear la entidad User
	newUser := User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Bankroll:     1000.00, // <--- Aquí asignamos el bono al campo correcto
	}

	// 4. Guardar en DB
	return s.repo.CreateUser(&newUser)
}

// LoginRequest define los datos para iniciar sesión
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) LoginUser(req LoginRequest) (string, error) {
	// 1. Buscar al usuario
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("credenciales inválidas") // No digas "email no existe" por seguridad
	}

	// 2. Verificar contraseña (Hash vs Plano)
	// bcrypt hace el trabajo sucio de comparar el hash guardado con lo que envían
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	// 3. Generar el JWT
	// Creamos los "Claims" (la información que va dentro del token)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Expira en 3 días
	}

	// Creamos el token sin firmar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Firmar el token con nuestro secreto del .env
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GetUserByID busca un usuario por su UUID (sin devolver la contraseña)
func (s *Service) GetUserProfile(id string) (*User, error) {
	var user User
	// Usamos First para buscar por ID
	// Omit("password_hash") asegura que NO traigamos el hash de la DB por seguridad extra
	err := s.repo.db.Omit("password_hash").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
