package consul

import (
	"fmt"
	"log"
	"math/rand"

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
	registration.Check.DeregisterCriticalServiceAfter = dto.HealthCheckDTO.Ttl
	consul.Agent().ServiceRegister(registration)
}

func LookupServiceWithConsul(serviceName string) (string, error) {
	config := api.DefaultConfig()
	consul, err := api.NewClient(config)
	if err != nil {
		return "", err
	}

	services, err := consul.Agent().ServicesWithFilterOpts("", &api.QueryOptions{
		AllowStale: true,
	})

	if err != nil {
		return "", err
	}
	srvc := services[serviceName]
	address := srvc.Address
	port := srvc.Port
	return fmt.Sprintf("http://%s:%v", address, port), nil
}

func HealthConsulService(service, tag string) (string, error) {
	passingOnly := true
	config := api.DefaultConfig()
	consul, err := api.NewClient(config)
	if err != nil {
		return "", err
	}
	addrs, _, err := consul.Health().Service(service, tag, passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return "",fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return "",  err
	}

	randomIndex := rand.Intn(len(addrs))
	serviceEntry := addrs[randomIndex]

	return fmt.Sprintf("http://%s:%v", serviceEntry.Service.Address, serviceEntry.Service.Port), nil
}
