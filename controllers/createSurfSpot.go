package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	models "surfSpot.com/surfSpot2-back/Models"
	"surfSpot.com/surfSpot2-back/initializers"
)

func CreateSurfSpot(c *gin.Context) {
	var spot models.SurfSpot
	var spotID int
	//Ajouter addedByUserId depuis JWT authentification

	cookie, cookieError := c.Request.Cookie("Authorization")
	if cookieError != nil {
		fmt.Println("cookie not found")
	}
	fmt.Println("CreateSurfSpot")
	fmt.Println(cookie)

	token, _ := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)

	userIdFloat, ok := claims["sub"].(float64)
	if !ok {
		fmt.Println("Invalid or missing user ID in token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}
	addedByUserId := int(userIdFloat)

	bindingError := c.Bind(&spot)
	if bindingError != nil {
		fmt.Println(bindingError)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failes to read body",
		})
		return
	}

	err := initializers.DB.QueryRow(context.Background(),
		`INSERT INTO surf_spots (destination, address, country, photo_url,
		 description, added_by_user_id, difficulty_level) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		spot.Destination,
		spot.Address,
		spot.Country,
		spot.PhotoURL,
		spot.Description,
		addedByUserId,
		spot.DifficultyLevel,
	).Scan(&spotID)

	if err != nil {
		fmt.Println("Erreur lors de l'insertion:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création du spot",
		})
		return
	}

	fmt.Println(spot.SurfBreak)
	breakName := spot.SurfBreak
	fmt.Println("breakName =", breakName)
	var breakID int
	err = initializers.DB.QueryRow(
		context.Background(),
		`SELECT id FROM surf_breaks WHERE name = $1`, breakName,
	).Scan(&breakID)

	fmt.Println("print breakID", breakID)
	if err != nil {
		fmt.Println("Erreur récupération break:", err)
		//continue // ou tu peux retourner une erreur
	}

	// Insère dans la table de liaison
	_, err = initializers.DB.Exec(
		context.Background(),
		`INSERT INTO spot_breaks (spot_id, break_id) VALUES ($1, $2)`,
		spotID, breakID,
	)

	if err != nil {
		fmt.Println("Erreur insertion liaison:", err)
	}

	c.JSON(http.StatusOK, spot)
}
