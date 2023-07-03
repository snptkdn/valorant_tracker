package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/exp/slog"

	"valorant/dto"
	"valorant/models"
)

const BASE = "https://api.henrikdev.xyz/valorant"

func GetMMR(user models.User) dto.MMRResponseDto {
	resp, err := http.Get(fmt.Sprintf("%s/v2/mmr/ap/%s/%s", BASE, user.ValorantName, user.ValorantTagName))
	if err != nil {
		slog.Error("Http Request failed!", "err", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("body read failed!", "err", err)
	}

  
  dto := dto.MMRResponseDto{}
	if err := json.Unmarshal(body, &dto); err != nil {
		slog.Error("JSON Unmarshal error", "err", err)
	}

  return dto
}

func GetMMRHistory(user models.User) dto.MMRHistoryResponseDto{
	resp, err := http.Get(fmt.Sprintf("%s/v1/mmr-history/ap/%s/%s", BASE, user.ValorantName, user.ValorantTagName))
	if err != nil {
		slog.Error("Http Request failed!", "err", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("body read failed!", "err", err)
	}

  
  dto := dto.MMRHistoryResponseDto{}
	if err := json.Unmarshal(body, &dto); err != nil {
		slog.Error("JSON Unmarshal error", "err", err)
	}

  return dto
}

func GetMatchHistory(user models.User) dto.MatchHistoryResponseDto {
	resp, err := http.Get(fmt.Sprintf("%s/v3/matches/ap/%s/%s", BASE, user.ValorantName, user.ValorantTagName))
	if err != nil {
		slog.Error("Http Request failed!", "err", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("body read failed!", "err", err)
	}

  
  dto := dto.MatchHistoryResponseDto{}
	if err := json.Unmarshal(body, &dto); err != nil {
		slog.Error("JSON Unmarshal error", "err", err)
	}

  return dto
}

