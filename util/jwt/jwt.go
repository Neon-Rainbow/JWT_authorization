package jwt

import (
	"JWT_authorization/config"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenExpireDuration  = time.Minute * 20    // Access token expiration time
	refreshTokenExpireDuration = time.Hour * 24 * 15 // Refresh token expiration time
)

// MyClaims is a custom JWT claims structure that includes standard claims along with user-specific information
type MyClaims struct {
	Username  string `json:"username"`
	UserID    uint   `json:"user_id"`
	TokenType string `json:"token_type"`
	IsAdmin   bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateToken generates both access and refresh JWT tokens for a given user
// @title GenerateToken
// @description Generates JWT tokens
// @param username string The username of the user
// @param userId uint The user ID
// @return accessToken string The access token
// @return refreshToken string The refresh token
// @return err error information
func GenerateToken(username string, userId uint, isAdmin bool) (accessToken string, refreshToken string, err error) {
	// Retrieve the secret key from the configuration
	var jwtSecret = config.GetConfig().JWT.Secret
	var mySecret = []byte(jwtSecret)

	// Channels to handle success and error signals
	successChannel := make(chan bool, 2)
	errorChannel := make(chan error, 2)

	// Create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Goroutine to generate the access token
	go func() {
		// Define the access token claims
		accessTokenClaims := MyClaims{
			Username:  username,
			UserID:    userId,
			TokenType: "access_token",
			IsAdmin:   isAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpireDuration)), // Expiration time
				IssuedAt:  jwt.NewNumericDate(time.Now()),                                // Issued time
			},
		}
		// Sign the access token with the secret key
		accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString(mySecret)
		if err != nil {
			errorChannel <- err
			return
		}
		successChannel <- true
	}()

	// Goroutine to generate the refresh token
	go func() {
		// Define the refresh token claims
		refreshTokenClaims := MyClaims{
			Username:  username,
			UserID:    userId,
			TokenType: "refresh_token",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpireDuration)), // Expiration time
				IssuedAt:  jwt.NewNumericDate(time.Now()),                                 // Issued time
			},
		}
		// Sign the refresh token with the secret key
		refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString(mySecret)
		if err != nil {
			errorChannel <- err
			return
		}
		successChannel <- true
	}()

	// Loop to wait for both tokens to be generated or an error/timeout to occur
	for i := 0; i < 2; i++ {
		select {
		case <-successChannel:
			// Successfully generated a token
		case err = <-errorChannel:
			// An error occurred while generating a token
			return "", "", err
		case <-ctx.Done():
			// Context timeout
			return "", "", ctx.Err()
		}
	}

	// Return the generated access and refresh tokens
	return accessToken, refreshToken, nil
}

// ParseToken parses a given JWT string and returns the claims if the token is valid
// @title ParseToken
// @description Parses a JWT token
// @param tokenString string The JWT token string
// @return *MyClaims The claims in the JWT token
// @return error Error information
func ParseToken(tokenString string) (*MyClaims, error) {
	// Retrieve the secret key from the configuration
	var jwtSecret = config.GetConfig().JWT.Secret
	var mySecret = []byte(jwtSecret) // Custom secret key

	// Parse the token with the custom claims structure
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return mySecret, nil
		})
	if err != nil {
		return nil, err
	}

	// Validate the token and extract the claims
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
