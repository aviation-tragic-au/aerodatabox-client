package model

import "time"

// FAALADDAircraftStatusContract represents the FAALADDAircraftStatus contract
type FAALADDAircraftStatusContract struct {
	ID            string     `json:"id"`
	IsBlocked     bool       `json:"isBlocked"`
	BlockedSince  *time.Time `json:"blockedSince"`
	LastBlockedOn *time.Time `json:"lastBlockedOn"`
}
