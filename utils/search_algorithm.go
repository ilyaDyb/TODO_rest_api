package utils

import (
	"log"
	"math"
	"strings"

	"github.com/ilyaDyb/go_rest_api/models"
)


func CalculateScore(user1, user2 models.User) float64 {
	const (
		interestsWeight = 0.2
		distanceWeight = 0.4
		cityWeight = 0.4
	)
	user1Hobbies := strings.Split(user1.Hobbies, ",")
	user2Hobbies := strings.Split(user2.Hobbies, ",")
	intersection := 0
	for _, hobby := range user1Hobbies {
		for _, hobby2 := range user2Hobbies {
			if hobby == hobby2 {
				intersection++
			}
		}
	}
	interestsScore := float64(intersection) / float64(len(user1Hobbies))
	distance := haversine(
		float64(user1.Lat), float64(user1.Lon), float64(user2.Lat), float64(user2.Lon),
	)
	distanceScore := math.Max(0, 1-distance/50)
	cityScore := 0.0
	if user1.City == user2.City {
		cityScore = 1.0
	}

	totalScore := interestsWeight*interestsScore + distanceWeight*distanceScore + cityWeight*cityScore
	log.Println(totalScore)
	return totalScore
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}