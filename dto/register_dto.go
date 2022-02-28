package dto

type ServiceRegisterDTO struct {
	ID					string 				`json:"id"`
	Name				string				`json:"name"`
	Port				int					`json:"port"`
	HealthCheckDTO		HealthCheckDTO		`json:"healthCheck"`
}
