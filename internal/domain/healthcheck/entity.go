package healthcheck

type HealthCheck struct {
	Status    string `json:"status,omitempty"`
	DBStatus  string `json:"db_status,omitempty"`
	Version   string `json:"version,omitempty"`
	Timestamp string `json:"timestamp"`
}
