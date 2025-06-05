package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	models "surfSpot.com/surfSpot2-back/Models"
	"surfSpot.com/surfSpot2-back/initializers"
)

func Signup(c *gin.Context) {
	//get email et password of req body
	var body struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
		Username string `json:"Username"`
	}

	bindingError := c.Bind(&body)
	if bindingError != nil {
		fmt.Println(bindingError)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failes to read body",
		})
		return
	}
	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failes to hash passwork",
		})
		return
	}
	//create the user
	user := models.User{Email: body.Email, Password: string(hash), Username: body.Username}
	_, err = initializers.DB.Exec(
		c,
		"INSERT INTO users (email, password, username) VALUES ($1, $2, $3)",
		user.Email, user.Password, user.Username,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	//Get the email et pass of req body
	var body struct {
		Email     string `json:"Email"`
		Passsword string `json:"Password"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failes to read body",
		})
		return
	}
	//look up request user
	var user models.User
	fmt.Println("did receive user")
	fmt.Println(body.Passsword)
	fmt.Println(body.Email)
	row := initializers.DB.QueryRow(
		c,
		"SELECT id, email, password, username FROM users WHERE email = $1",
		body.Email,
	)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Username)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found",
		})
		return
	}

	//Compare hashed password and password.
	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Passsword))

	if result != nil {
		fmt.Println(result)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid Email or password",
		})
		return
	}
	//Generate a jwt(json web token) token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failes to creat token",
		})
		return
	}
	//send the reponse
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Utilisateur non attaché"})
		return
	}

	c.JSON(http.StatusOK, user)
}
