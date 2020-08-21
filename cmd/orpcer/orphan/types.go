package orphan

import "time"

type Cert struct {
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
	Age       int       `json:"age"`
	Data      string    `json:"cert"`
}

type DBAdapter interface {
	GetCerts() ([]Cert, error)
}
