package dto

type ServiceRegisterDTO struct {
	Name				string				`json:"name"`
	Port				int					`json:"port"`
	HealthCheckDTO		HealthCheckDTO		`json:"healthCheck"`
}
