package model

type HealthCheck struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}
