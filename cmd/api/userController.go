package main

import (
	"net/http"
	"rest_app_in_gin/internal/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type registerUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}
// loginUser authenticates user
//
// @Summary Login user
// @Description Authenticates a user and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body loginUserRequest true "Login credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (app *application) loginUserHandler(c *gin.Context) {
	var auth loginUserRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": existingUser.Id,
		"expires_at": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(app.jwtSecret)) //jwtSecret прописаний в main як структура та в енв для самого додатку
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, loginResponse{Token: tokenString})
}
	
// createUser registers a new user
//
// @Summary Register user
// @Description Creates a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body registerUserRequest true "User registration data"
// @Success 201 {object} database.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (app *application) createUserHandler(c *gin.Context) {
	var req registerUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := database.User{
		Name:  req.Name,
		Email: req.Email,
		Password: string(hashedPassword),
	}

	if err := app.models.Users.Insert(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	//c.JSON(http.StatusCreated, gin.H{"user": user})
	c.JSON(http.StatusCreated, user)
}