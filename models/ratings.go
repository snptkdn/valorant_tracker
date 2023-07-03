package models

import "time"

type Ratings struct {
  PlayerID string
  Date time.Time
  MatchID string
  Rating int
}
