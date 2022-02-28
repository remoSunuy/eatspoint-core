package consul

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
	"github.com/remoSunuy/eatspoint-core/dto"
	"github.com/remoSunuy/eatspoint-core/service"
)

func RegisterServiceWithConsul(dto *dto.ServiceRegisterDTO) {

  	consul, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		log.Fatalln(err)
	}

	registration := new(api.AgentServiceRegistration)
	registration.ID = dto.ID
	registration.Name = dto.Name
	address := service.Hostname()
	registration.Address = address
	registration.Port = dto.Port
	registration.Check = new(api.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/%s", address, registration.Port, dto.HealthCheckDTO.Service)
	registration.Check.Interval = dto.HealthCheckDTO.Interval
	registration.Check.Timeout = dto.HealthCheckDTO.Timeout
	consul.Agent().ServiceRegister(registration)
}

func LookupServiceWithConsul(serviceName string) (string, error) {
	config := api.DefaultConfig()
	consul, err := api.NewClient(config)
	if err != nil {
		return "", err
	}
	services, err := consul.Agent().Services()
	if err != nil {
		return "", err
	}
	srvc := services["product-service"]
	address := srvc.Address
	port := srvc.Port
	return fmt.Sprintf("http://%s:%v", address, port), nil
}