package services

import (
	"bytes"
	"encoding/json"
	"fmt"
  "strings"
	"net/http"
	"time"

	"golang.org/x/exp/slog"

	"valorant/dto"
	"valorant/models"
)


func SendNotify(users []models.User) {
	body := dto.WebhookRequestDto{}
	body.Content = "Today's result!"
  for _, user := range(users) {
    body.Embeds = append(body.Embeds, EmbedFromUser(user))
  }

	sample_json, _ := json.Marshal(body)
	res, err := http.Post(WebhookURL, "application/json", bytes.NewBuffer(sample_json))
	if err != nil {
		slog.Error("can't post discord webhook", "err", err)
	}
	defer res.Body.Close()
}

func EmbedFromUser(user models.User) dto.Embed{
  mmr := GetMMR(user)
  mmrHistory := GetMMRHistory(user)
  matchHistory := GetMatchHistory(user)

  return dto.Embed{
    Title: user.ValorantName,
    Description: fmt.Sprintf("Discord is <@%s>", user.DiscordID),
    Thumbnail: dto.Thumbnail{
      URL: mmr.Data.CurrentData.Imges.Large,
    },
    Color: ColorByRank(mmr.Data.CurrentData.CurrentTierPatched),
    Fields: []dto.Field{
      {
        Name: "CurrentRank",
        Value: mmr.Data.CurrentData.CurrentTierPatched,
        Inline: true,
      },
      {
        Name: "RP",
        Value: fmt.Sprintf(
          "%d(%s)",
          mmr.Data.CurrentData.RP,
          mmrHistory.TotalRPByDate(time.Now()),
        ),
        Inline: true,
      },
      {
        Name: "EloRating",
        Value: fmt.Sprintf(
          "%d(%s)",
          mmr.Data.CurrentData.EloRating,
          mmrHistory.LatestEloRating(time.Now().Add(-24 * time.Hour), mmr.Data.CurrentData.EloRating),
        ),
        Inline: true,
      },
      {
        Name: "KDA(Ave)",
        Value: fmt.Sprintf(
          "%.2f",
          matchHistory.AverageKDA(time.Now(), user.ValorantName)),
        Inline: true,
      },
      {
        Name: "HS%(Ave)",
        Value: fmt.Sprintf(
          "%.2f%%",
          matchHistory.HeadShotPercentage(time.Now(), user.ValorantName)),
        Inline: true,
      },
      {
        Name: "UseAbillity(Ave)",
        Value: fmt.Sprint(matchHistory.AverageAbilityCast(time.Now(), user.ValorantName)),
        Inline: true,
      },
      {
        Name: "History",
        Value: mmrHistory.WinLoseString(time.Now()),
        Inline: true,
      },
    },
    
  }
}

func ColorByRank(rank string) int {
  if strings.Contains(rank, "Iron") {
    return 6908265
  } else if strings.Contains(rank, "Bronze") {
    return 12092939
  } else if strings.Contains(rank, "Silver") {
    return 16119285
  } else if strings.Contains(rank, "Gold") {
    return 16766720
  } else if strings.Contains(rank, "Platinum") {
    return 49151
  } else {
    return 0
  }
}
