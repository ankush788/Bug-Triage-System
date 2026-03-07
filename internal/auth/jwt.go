package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents JWT claims for our application
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	secretKey string
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{secretKey: secretKey}
}

// GenerateToken creates a new JWT token for a user
func (m *JWTManager) GenerateToken(userID int64, email string, expiresIn time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) 
    // it give json { header , claims , Method}
	// Header → { "typ": "JWT", "alg": "HS256" }  (HS256 --> signature algorithm)
	// Claims → payload (userID, email, exp, iat)
	// Method → interface name in which signature function is implemneted
	
	signedToken, err := token.SignedString([]byte(m.secretKey)) //base64(header) . base64(payload) . signature
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken parses and validates a JWT token

//JWT token
//    ↓
// Parse
//    ↓
// Verify signature with secret
// break token(header.payload.signature) in three part --> check  HMAC-SHA256(secretKey, header.payload) == token_signature
//    ↓
// Check expiration
//    ↓
// Return user claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
