package cloud_pockets

import "database/sql"

type CloudPocket struct {
	PocketID   int     `json:"pocketID"`
	PocketName string  `json:"pocketName"`
	Category   string  `json:"category"`
	Currency   string  `json:"currency"`
	Balance    float64 `json:"balance"`
}

type handler struct {
	db *sql.DB
}

func New(db *sql.DB) *handler {
	return &handler{db}
}
