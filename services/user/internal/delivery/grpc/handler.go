package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	bcrpyt "github.com/samarthasthan/21BRS1248_Backend/common/bycrpyt"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Register registers a new user
func (s *UserService) Register(ctx context.Context, in *proto_go.RegisterRequest) (*proto_go.RegisterResponse, error) {
	uuid := uuid.New().String()
	// Hash the password
	var err error
	in.Password, err = bcrpyt.HashPassword(in.Password)
	if err != nil {
		s.log.Errorf("Failed to hash password for user email: %s: %v", in.Email, err)
		return &proto_go.RegisterResponse{
			Success: false,
			Message: "Failed to register user",
		}, err
	}
	err = s.repo.RegisterUser(ctx, in, uuid)
	if err != nil {
		s.log.Errorf("Failed to register user with email: %s: %v", in.Email, err)
		return &proto_go.RegisterResponse{
			Success: false,
			Message: "Failed to register user",
		}, err
	}
	return &proto_go.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

// Login logs in a user
func (s *UserService) Login(ctx context.Context, in *proto_go.LoginRequest) (*proto_go.LoginResponse, error) {
	password, err := s.repo.GetPasswordByEmail(ctx, in)
	if err != nil {
		s.log.Errorf("Failed to login user with email: %s: %v", in.Email, err)
		return &proto_go.LoginResponse{
			SessionId: "",
			Success:   false,
			ExpiresAt: nil,
			Message:   "Failed to login user",
		}, err
	}
	if !bcrpyt.ValidatePassword(password, in.Password) {
		return &proto_go.LoginResponse{
			SessionId: "",
			Success:   false,
			ExpiresAt: nil,
			Message:   "Invalid email or password",
		}, err
	}

	// Create a JWT token
	token, time, err := s.createToken(in)
	if err != nil {
		s.log.Errorf("Failed to create token for user with email: %s: %v", in.Email, err)
		return &proto_go.LoginResponse{
			SessionId: "",
			Success:   false,
			ExpiresAt: nil,
			Message:   "Failed to login user",
		}, err
	}
	return &proto_go.LoginResponse{
		SessionId: token,
		Success:   true,
		ExpiresAt: &timestamppb.Timestamp{
			Seconds: time.Unix(),
		},
		Message: "User logged in successfully",
	}, nil

}

// CheckJWT checks the validity of a JWT
func (s *UserService) CheckJWT(ctx context.Context, in *proto_go.CheckJWTRequest) (*proto_go.CheckJWTResponse, error) {
	claims, err := s.VerifyToken(in.SessionId)
	if err != nil {
		return &proto_go.CheckJWTResponse{
			Valid: false,
		}, nil
	}
	if claims["expires_at"].(float64) < float64(time.Now().Unix()) {
		return &proto_go.CheckJWTResponse{
			Valid: false,
		}, nil
	}
	return &proto_go.CheckJWTResponse{
		Valid:     true,
		Email:     claims["email"].(string),
		ExpiresAt: claims["expires_at"].(*timestamppb.Timestamp),
	}, nil
}

// Helper function to create a JWT token
func (s *UserService) createToken(in *proto_go.LoginRequest) (string, *time.Time, error) {
	time := time.Now().Add(356 * time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email":      in.Email,
			"expires_at": time,
		})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", &time, err
	}
	return tokenString, &time, nil
}

// Validate the JWT token
func (s *UserService) VerifyToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
