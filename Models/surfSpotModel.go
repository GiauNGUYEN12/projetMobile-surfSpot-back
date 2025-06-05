package models

type SurfSpot struct {
	Id              int    `json:"id"`
	Destination     string `json:"destination"`
	Address         string `json:"address"`
	Country         string `json:"country"`
	DifficultyLevel int    `json:"difficulty_level"`
	//SurfBreak       string   `json:"surfBreak"`
	Description   string `json:"description"`
	PhotoURL      string `json:"photo_url"` //il n'y a pas d'espace entre json: et "le mot clé"
	AddedByUserId int    `json:"added_by_user_id"`
	SurfBreak     string `json:"surf_breaks"`
}
