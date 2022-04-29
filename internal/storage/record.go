package storage

type record struct {
	ID         string `json:"PK"`
	SK         string `json:"SK"`
	RecordType string `json:"type"`
	Version    int    `json:"version"`
}

type audit struct {
	CreatedAt int64 `json:"CreatedAt"`
	UpdatedAt int64 `json:"UpdatedAt"`
}
