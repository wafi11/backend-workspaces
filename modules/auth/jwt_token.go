package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/wafi11/backend-workspaces/pkg/config"
)

type JwtTokenRequest struct {
	UserId int `json:"userId"`
}
type TokenClaims struct {
	UserId int `json:"userId"`
	jwt.StandardClaims
}

func GenerateTokenPair(userID int, cfg config.Config) (accessToken, refreshToken string, err error) {
	req := JwtTokenRequest{
		UserId: userID,
	}

	// Access token: 15 menit
	accessToken, err = GenereteToken(req, cfg.Duration.AccessToken, cfg)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Refresh token: 7 hari
	refreshToken, err = GenereteToken(req, cfg.Duration.AccessToken, cfg)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString string, cfg config.Config) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token string is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC (HS256).
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.SecretKey.JwtSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func GenereteToken(req JwtTokenRequest, duration int, cfg config.Config) (string, error) {
	secretKey := cfg.SecretKey.JwtSecretKey
	if secretKey == "" {
		return "", errors.New("jwt secret not found")
	}

	now := time.Now()
	expirationTime := now.Add(time.Duration(duration))

	claims := TokenClaims{
		UserId: req.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
