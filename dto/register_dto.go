package dto

type ServiceRegisterDTO struct {
	Name				string				`json:"name"`
	IP					string				`json:"ip"`
	Port				int					`json:"port"`
	HealthCheckDTO		HealthCheckDTO		`json:"healthCheck"`
}
