package storage

type record struct {
	ID         string `json:"pk"`
	SK         string `json:"sk"`
	RecordType string `json:"type"`
	Version    int    `json:"version"`
}
