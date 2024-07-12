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
	Login(ctx context.Context, email, password, clientIP string) (string, *domain.User, error)
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

func (a *authService) Login(ctx context.Context, email, password, clientIP string) (string, *domain.User, error) {
	if email == "" || password == "" {
		return "", nil, domain.NewValidationError("email and password are required", nil)
	}

	userEmailFilter := domain.UserFilters{
		Email: []string{
			email,
		},
	}

	users, err := a.userRepository.Search(ctx, userEmailFilter)
	if err != nil {
		return "", nil, err
	}

	if len(users) <= 0 {
		return "", nil, domain.NewValidationError("no user found", nil)
	}
	loggedUser := users[0]

	if !loggedUser.ComparePassword(password) {
		return "", nil, domain.NewValidationError("password is incorrect", nil)
	}

	createdToken, err := a.CreateToken(loggedUser.UserID)
	if err != nil {
		return "", nil, err
	}

	loggedUser.MergeUpdate(domain.UserUpdate{LastLoggedIP: &clientIP}, "")
	err = a.userRepository.Update(ctx, loggedUser)
	if err != nil {
		return "", nil, err
	}

	return createdToken, &loggedUser, nil
}

// Logout implements AuthService.
func (a *authService) Logout(ctx context.Context) error {
	panic("unimplemented")
}

func (a *authService) CreateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(12 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(a.jwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (a *authService) VerifyUserSession(ctx context.Context, userID, clientIP string) error {
	fmt.Println("userID", userID)
	fmt.Println("clientIP", clientIP)

	user, err := a.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	fmt.Println(user)

	if user.LastLoggedIP != clientIP {
		return domain.NewUnauthorizedError("client ip is different from the logged one, please login again!")
	}

	return nil
}
