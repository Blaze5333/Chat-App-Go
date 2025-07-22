package user

import (
	"chat-server/db"
	"chat-server/models"
	"chat-server/services"
	"chat-server/tokens"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection = db.UserData(db.Client, "users")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic("Error hashing password:", err)
		return "", err
	}

	return string(bytes), nil
}
func VerifyPassword(userpassword, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.UserRegisterReq
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err, "message": "Invalid request"})
			return
		}
		var existingUser models.User
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
		if err == nil {
			if existingUser.Verified {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already registered", "message": "Please login or use a different email"})
				return
			} else {
				otp := services.GenerateOTP()
				expiry := time.Now().Add(10 * time.Minute)
				hashedPassword, err := HashPassword(user.Password)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err, "message": "Error hashing password"})
					return
				}
				update := bson.M{
					"$set": bson.M{
						"otp":         otp,
						"otp_expires": expiry,
						"username":    user.Username,
						"verified":    false,
						"password":    hashedPassword,
					},
				}
				_, err = UserCollection.UpdateOne(ctx, bson.M{"email": user.Email}, update)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP", "message": err.Error()})
					return
				}
				err = services.SendEmail(user.Email, otp)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email", "message": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email. Please verify your email to complete registration."})
				return
			}
		}
		var userData models.User
		userData.Username = user.Username
		userData.ID = primitive.NewObjectID()
		userData.Email = user.Email
		userData.UserId = userData.ID.Hex()
		otp := services.GenerateOTP()
		userData.Otp = otp
		userData.OtpExpires = time.Now().Add(10 * time.Minute)
		userData.Verified = false
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err, "message": "Error hashing password"})
			return
		}
		userData.Password = hashedPassword
		_, err = UserCollection.InsertOne(ctx, userData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err, "message": "Error inserting user"})
			return
		}
		err = services.SendEmail(user.Email, otp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email", "message": err.Error()})
			return
		}
		log.Println("OTP sent to email:", user.Email)

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": gin.H{
			"id":       userData.UserId,
			"username": userData.Username,
			"email":    userData.Email,
		}})

	}
}
func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.UserLoginReq
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err, "message": "Invalid request"})
			return
		}
		var foundUser models.User
		err := UserCollection.FindOne(ctx, primitive.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password", "message": "Please check your credentials"})
			return
		}
		if !foundUser.Verified {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified", "message": "Please verify your email before logging in"})
			return
		}
		isValid, msg := VerifyPassword(user.Password, foundUser.Password)
		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg, "message": "Invalid password"})
			return
		}
		// Generate JWT token here (not implemented in this snippet)
		token, err := tokens.GenerateToken(foundUser.Email, foundUser.UserId, foundUser.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": gin.H{
			"id":       foundUser.UserId,
			"username": foundUser.Username,
			"email":    foundUser.Email,
			"token":    token,
			"image":    foundUser.Image,
		}})
	}
}
func GetUserByEmail() gin.HandlerFunc {
	fmt.Println("Email to search:")
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		email := c.Query("email")
		var user models.User
		err := UserCollection.FindOne(ctx, primitive.M{"email": email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "message": "Please check the email"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"email": user.Email, "username": user.Username, "id": user.UserId, "image": user.Image})
	}
}
func SocialLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		q := c.Request.URL.Query()
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}
func SocialLoginCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		q := c.Request.URL.Query()
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			log.Println("Error completing user auth:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete social login", "message": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var existingUser models.User
		err = UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
		if err == nil {
			if existingUser.GoogleLogin {
				token, err := tokens.GenerateToken(existingUser.Email, existingUser.UserId, existingUser.Username)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "message": err.Error()})
					return
				}

				userData := url.QueryEscape(fmt.Sprintf(`{"id":"%s","username":"%s","email":"%s","image":"%s"}`,
					existingUser.UserId, existingUser.Username, existingUser.Email, existingUser.Image))
				redirectURL := fmt.Sprintf("http://localhost:3000/auth/google/callback?token=%s&user=%s", token, userData)
				c.Redirect(http.StatusFound, redirectURL)
				return
			} else {
				update := bson.M{
					"$set": bson.M{
						"google_login": true,
					},
				}
				_, err = UserCollection.UpdateOne(ctx, bson.M{"email": user.Email}, update)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "message": err.Error()})
					return
				}
				token, err := tokens.GenerateToken(existingUser.Email, existingUser.UserId, existingUser.Username)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "message": err.Error()})
					return
				}

				userData := url.QueryEscape(fmt.Sprintf(`{"id":"%s","username":"%s","email":"%s","image":"%s"}`,
					existingUser.UserId, existingUser.Username, existingUser.Email, existingUser.Image))
				redirectURL := fmt.Sprintf("http://localhost:3000/auth/google/callback?token=%s&user=%s", token, userData)
				c.Redirect(http.StatusFound, redirectURL)
				return
			}
		} else {
			var newUser models.User
			newUser.ID = primitive.NewObjectID()
			newUser.Username = user.NickName
			newUser.Email = user.Email
			newUser.UserId = newUser.ID.Hex()
			newUser.Image = user.AvatarURL
			newUser.GoogleLogin = true
			newUser.Verified = true
			UserCollection.InsertOne(ctx, newUser)
			token, err := tokens.GenerateToken(newUser.Email, newUser.UserId, newUser.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "message": err.Error()})
				return
			}

			userData := url.QueryEscape(fmt.Sprintf(`{"id":"%s","username":"%s","email":"%s","image":"%s"}`,
				newUser.UserId, newUser.Username, newUser.Email, newUser.Image))
			redirectURL := fmt.Sprintf("http://localhost:3000/auth/google/callback?token=%s&user=%s", token, userData)
			c.Redirect(http.StatusFound, redirectURL)
			return

		}
	}
}

func VerifyOtp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var request struct {
			Email string `json:"email" binding:"required"`
			Otp   string `json:"otp" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
			return
		}
		var user models.User
		err := UserCollection.FindOne(ctx, bson.M{"email": request.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "message": "Please check the email"})
			return
		}
		if user.Otp != request.Otp || time.Now().After(user.OtpExpires) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP", "message": "Please try again"})
			return
		}
		user.Verified = true
		user.Otp = ""
		user.OtpExpires = time.Time{}
		_, err = UserCollection.UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": user})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
	}
}

func UploadHandler(c *gin.Context) {
	userId, _ := c.Get("user_id")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": "Please login to upload an image"})
		return
	}
	filter := bson.M{"user_id": userId.(string)}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// 2. Load AWS config (uses env vars or shared config files)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // replace with your actual bucket's region
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AWS config error"})
		return
	}

	svc := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(svc)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String("gochatappimages"),
		Key:         aws.String(userId.(string)),
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	update := bson.M{"$set": bson.M{"image": result.Location}}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Failed to update user image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "image_url": result.Location})
}
