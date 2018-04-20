package shuttletracker

// Stop is a place where vehicles frequently stop.
type Stop struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Enabled     bool    `json:"enabled"`
}

// StopService is an interface for interacting with Stops.
type StopService interface {
	Stops() ([]*Stop, error)
	CreateStop(stop *Stop) error
	DeleteStop(id int) error
}
