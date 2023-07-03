package dto

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/exp/slog"
)

type MMRResponseDto struct {
	Status int     `json:"status"`
	Data   MMRData `json:"data"`
}

type MMRData struct {
	CurrentData CurrentMMRData `json:"current_data"`
}

type CurrentMMRData struct {
	CurrentTierPatched string `json:"currenttierpatched"`
	RP                 int    `json:"ranking_in_tier"`
	EloRating          int    `json:"elo"`
	Imges              Images `json:"images"`
}

type Images struct {
	Large string `json:"large"`
}

type MMRHistoryResponseDto struct {
	Datas []MMRHistoryData `json:"data"`
  Name string `json:"name"`
}

type MMRHistoryData struct {
	CurrentTierPatched string  `json:"currenttierpatched"`
	MMRDiff            int     `json:"mmr_change_to_last_game"`
	EloRating          int     `json:"elo"`
	Date               ApiDate `json:"date_raw"`
	MatchID            string  `json:"match_id"`
}

type ApiDate struct {
	time.Time
}

type MatchHistoryResponseDto struct {
	Data []MatchHistoryData `json:"data"`
}

type MatchHistoryData struct {
	PlayersData AllPlayersData `json:"players"`
	MetaData    MetaData       `json:"metadata"`
}

type MetaData struct {
	Date    ApiDate `json:"game_start"`
	DateFmt string  `json:"game_start_patched"`
}

type AllPlayersData struct {
	AllPlayers []PlayersData `json:"all_players"`
}

type PlayersData struct {
	ID          string      `json:"puuid"`
	Name        string      `json:"name"`
	AbilityCast AbilityCast `json:"ability_casts"`
	Stats       Stats       `json:"stats"`
	Damage      uint        `json:"damage_made"`
}

type AbilityCast struct {
	C uint `json:"c_cast"`
	Q uint `json:"q_cast"`
	E uint `json:"e_cast"`
	X uint `json:"x_cast"`
}

type Stats struct {
	Score    uint `json:"score"`
	Kill     uint `json:"kills"`
	Death    uint `json:"deaths"`
	Assist   uint `json:"assists"`
	BodyShot uint `json:"bodyshots"`
	HeadShot uint `json:"headshots"`
	LegShot  uint `json:"legshots"`
}

func (t *ApiDate) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	data_string := string(data)
	data_int, err := strconv.Atoi(data_string)
	if err != nil {
		slog.Error("api time parse error", "err", err)
	}
	dtFromUnix := time.Unix(int64(data_int), 0)
	*t = ApiDate{dtFromUnix}
	return err
}

func (historyData *MMRHistoryResponseDto) TotalRPByDate(date time.Time) string {
	var totalRP int = 0
	for _, history := range historyData.Datas {
		if history.Date.Time.Before(date.Add(-24 * time.Hour)) {
			continue
		}

		totalRP += history.MMRDiff
	}

	if totalRP < 0 {
		return fmt.Sprintf("%d", totalRP)
	} else if totalRP > 0 {
		return fmt.Sprintf("+%d", totalRP)
	} else {
		return fmt.Sprint("±0")
	}
}

func (historyData *MMRHistoryResponseDto) WinLoseString(date time.Time) string {
	var win = 0
	var lose = 0
	for _, history := range historyData.Datas {
		if history.Date.Time.Before(date.Add(-24 * time.Hour)) {
			continue
		}

		if history.MMRDiff < 0 {
			lose += 1
		} else if history.MMRDiff > 0 {
			win += 1
		}
	}

	return fmt.Sprintf("%dW-%dL", win, lose)
}

func (historyData *MMRHistoryResponseDto) LatestEloRating(date time.Time, compared_elo int) string {
	var latest_date = time.Date(1999, 1, 1, 0, 0, 0, 0, time.Local)
	var latest_elo = 0
	for _, history := range historyData.Datas {
		if history.Date.Time.After(date) {
			continue
		}

		if history.Date.Time.After(latest_date) {
			latest_date = history.Date.Time
			latest_elo = history.EloRating
		}

	}

	diff_elo := compared_elo - latest_elo
	if diff_elo < 0 {
		return fmt.Sprintf("%d", diff_elo)
	} else if diff_elo > 0 {
		return fmt.Sprintf("+%d", diff_elo)
	} else {
		return fmt.Sprint("±0")
	}

}

func (matchHistoryData MatchHistoryResponseDto) AverageAbilityCast(date time.Time, name string) uint {
	var sum uint = 0
	matchCount := 0
	for _, match := range matchHistoryData.Data {
		if match.MetaData.Date.Time.Before(date.Add(-24 * time.Hour)) {
			continue
		}

		for _, player := range match.PlayersData.AllPlayers {
			if player.Name != name {
				continue
			}

			matchCount++
			sum += player.AbilityCast.C + player.AbilityCast.Q + player.AbilityCast.E + player.AbilityCast.X
		}
	}

	if matchCount > 0 {
		return sum / uint(matchCount)
	} else {
		return 0
	}
}

func (matchHistoryData MatchHistoryResponseDto) AverageKDA(date time.Time, name string) float32 {
	var sum_kill float32 = 0
	var sum_death float32 = 0
	var sum_assist float32 = 0
	for _, match := range matchHistoryData.Data {
		if match.MetaData.Date.Time.Before(date.Add(-24 * time.Hour)) {
			continue
		}
		for _, player := range match.PlayersData.AllPlayers {
			if player.Name != name {
				continue
			}

			sum_kill += float32(player.Stats.Kill)
			sum_death += float32(player.Stats.Death)
			sum_assist += float32(player.Stats.Assist)
		}
	}

	if sum_death > 0 {
		return ((sum_kill + sum_assist) / sum_death)
	} else {
		return 0
	}
}

func (matchHistoryData MatchHistoryResponseDto) HeadShotPercentage(date time.Time, name string) float32 {
	var sumHead float32 = 0
	var sumExcludeHead float32 = 0
	for _, match := range matchHistoryData.Data {
		if match.MetaData.Date.Time.Before(date.Add(-24 * time.Hour)) {
			continue
		}
		for _, player := range match.PlayersData.AllPlayers {
			if player.Name != name {
				continue
			}

			sumHead += float32(player.Stats.HeadShot)
			sumExcludeHead += float32(player.Stats.BodyShot + player.Stats.LegShot)
		}
	}

	if (sumHead + sumExcludeHead) > 0 {
		return (sumHead / (sumHead + sumExcludeHead)) * 100
	} else {
		return 0.0
	}
}

func (matchHistoryData MatchHistoryResponseDto) AverageDamage(date time.Time, name string) uint {
	var sum uint = 0
	matchCount := 0
	for _, match := range matchHistoryData.Data {
		if match.MetaData.Date.Time.Before(date.Add(-24 * time.Hour)) {
			continue
		}
		for _, player := range match.PlayersData.AllPlayers {
			if player.Name != name {
				continue
			}

			matchCount++
			sum += player.Damage
		}
	}

	if matchCount > 0 {
		return sum / uint(matchCount)
	} else {
		return 0
	}
}

func (matchHistoryData MatchHistoryResponseDto) AverageScore(date time.Time, name string) uint {
	var sum uint = 0
	matchCount := 0
	for _, match := range matchHistoryData.Data {
		if match.MetaData.Date.Time.Before(date.Add(-24 * time.Hour)) {
			continue
		}
		for _, player := range match.PlayersData.AllPlayers {
			if player.Name != name {
				continue
			}

			matchCount++
			sum += player.Stats.Score
		}
	}

	if matchCount > 0 {
		return sum / uint(matchCount)
	} else {
		return 0
	}
}
