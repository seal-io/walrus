package types

import "time"

// CatalogSync is the sync information of a catalog.
type CatalogSync struct {
	Total     int       `json:"total"`
	Succeeded int       `json:"succeeded"`
	Failed    int       `json:"failed"`
	Time      time.Time `json:"time"`
}
