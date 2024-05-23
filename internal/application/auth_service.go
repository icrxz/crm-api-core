package application

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type authService struct {
	userRepository domain.UserRepository
	jwtSecretKey   string
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context) error
	CreateToken(userID string) (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
	VerifyUserSession(ctx context.Context, userID, clientIP string) error
}

func NewAuthService(userRepository domain.UserRepository, jwtSecretKey string) AuthService {
	return &authService{
		userRepository: userRepository,
		jwtSecretKey:   jwtSecretKey,
	}
}

func (a *authService) Login(ctx context.Context, email, password string) (string, error) {
	userEmailFilter := domain.UserFilters{
		Email: []string{
			email,
		},
	}

	user, err := a.userRepository.Search(ctx, userEmailFilter)
	if err != nil {
		return "", err
	}

	if len(user) <= 0 {
		return "", domain.NewValidationError("no user found", nil)
	}

	if !user[0].ComparePassword(password) {
		return "", domain.NewValidationError("password is incorrect", nil)
	}

	return a.CreateToken(user[0].UserID)
}

// Logout implements AuthService.
func (a *authService) Logout(ctx context.Context) error {
	panic("unimplemented")
}

func (a *authService) CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(a.jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("Invalid signing method")
		}

		return a.jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}

func (a *authService) VerifyUserSession(ctx context.Context, userID, clientIP string) error {
	fmt.Println("userID", userID)
	fmt.Println("clientIP", clientIP)

	user, err := a.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.LastLoggedIP != clientIP {
		return domain.NewUnauthorizedError("client ip is different from the logged one, please login again!")
	}

	return nil
}
