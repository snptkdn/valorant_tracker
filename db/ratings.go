package db

import (
	"valorant/models"
	"valorant/services"
)

func InsertMMRHistory(user models.User) {
	db := Db()
	histories := services.GetMMRHistory(user)
	for _, history := range histories.Datas {
		var count int64
		db.Model(&models.Ratings{}).Where("match_id = ? AND player_id = ?", history.MatchID, user.ValorantName).Count(&count)

		if count == 0 {
			db.Create(
				models.Ratings{
					MatchID:  history.MatchID,
					PlayerID: histories.Name,
					Date:     history.Date.Time,
					Rating:   history.EloRating,
				},
			)
		}
	}
}
