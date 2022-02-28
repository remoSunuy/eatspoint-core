package dto

type ServiceRegisterDTO struct {
	Name				string
	IP					string
	Port				int
	HealthCheckDTO		HealthCheckDTO
}
