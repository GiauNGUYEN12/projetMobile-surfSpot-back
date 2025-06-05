package main

import (
	"context"
	"fmt"

	models "surfSpot.com/surfSpot2-back/Models"
	"surfSpot.com/surfSpot2-back/controllers"
	"surfSpot.com/surfSpot2-back/initializers"
	"surfSpot.com/surfSpot2-back/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnVariables()
	initializers.ConnectToDb()
}

func main() {
	router := gin.Default()
	router.GET("/surfSpots", getSurfSpots)
	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/surfSpots", controllers.CreateSurfSpot)
	//router.GET("/surfSpots/:id", getSurfSpotById)
	//router.POST("/surfSpots", postSurfSpots)

	router.Run()
}

func getSurfSpots(c *gin.Context) {
	fmt.Println("Did receive request")
	rows, err := initializers.DB.Query(context.Background(), `SELECT 
		surf_spots.id AS spot_id, destination, address, country, photo_url, description, added_by_user_id, difficulty_level, surf_breaks.name AS break_name
		FROM spot_breaks
		INNER JOIN surf_spots ON spot_breaks.spot_id = surf_spots.id
		INNER JOIN surf_breaks ON spot_breaks.break_id = surf_breaks.id
	`)

	if err != nil {
		fmt.Printf("Unable to establish connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur DB"})
		return
	}
	defer rows.Close()

	var spots []models.SurfSpot

	for rows.Next() {
		fmt.Println("✅ did find row!")
		var spot models.SurfSpot
		err := rows.Scan(&spot.Id, &spot.Destination, &spot.Address, &spot.Country, &spot.PhotoURL, &spot.Description, &spot.AddedByUserId, &spot.DifficultyLevel, &spot.SurfBreak)
		if err != nil {
			fmt.Printf("Scan error: %v", err)
			continue
		}
		fmt.Printf("Did find destination  %s\n", spot.Destination)
		spots = append(spots, spot)
	}

	c.JSON(http.StatusOK, spots)
}
