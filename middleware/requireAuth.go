package middleware

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	models "surfSpot.com/surfSpot2-back/Models"
	"surfSpot.com/surfSpot2-back/initializers"
)

func RequireAuth(c *gin.Context) {
	// Lire le token JWT stocké dans un cookie nommé Authorization
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
		return
	}

	// Parser et valider le token JWT
	//Verifier que le token est bien signé avec la bonne clé secrète.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
		return
	}

	//Recuperer les claims du token (contenu du token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token mal formé"})
		return
	}

	//Vérifier expiration
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expiré"})
		return
	}

	// Trouver l’utilisateur en base via l’ID contenu dans le `sub`
	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "ID utilisateur invalide dans le token"})
		return
	}
	userID := int(userIDFloat)

	var user models.User
	err = initializers.DB.QueryRow(
		context.Background(),
		"SELECT id, email, username FROM users WHERE id=$1",
		userID,
	).Scan(&user.ID, &user.Email, &user.Username)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	//Attacher l’utilisateur au contexte Gin
	c.Set("user", user)

	//Continuer la requête
	c.Next()
}
