package tokens

import (
	"chat-server/db"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type SignedDetails struct {
	UserId   string
	Email    string
	Username string
	jwt.StandardClaims
}

var UserData = db.UserData(db.Client, "users")
var SECRET_KEY string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	SECRET_KEY = os.Getenv("SECRET_KEY") // Replace with your actual secret key
}
func GenerateToken(email, userId, username string) (string, error) {
	claims := &SignedDetails{
		UserId:   userId,
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24 * 10).Unix(),
			Issuer:    "chat-server",
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println("Error generating token:", err)
		return "", err
	}
	return token, nil
}
func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		log.Println("Error parsing token:", err)
		return nil, err
	}

	if !token.Valid {
		log.Println("Invalid token")
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
