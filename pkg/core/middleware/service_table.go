package middleware

import (
	"fmt"
	servicetypes "github.com/aioncore/ouroboros/pkg/service/types/services"
	dbm "github.com/tendermint/tm-db"
)

type ServiceTable interface {
	GetServices() ([]servicetypes.ServiceData, error)
	GetServiceByType(requiredType string) (servicetypes.ServiceData, error)
	UseService(serviceData servicetypes.ServiceData) error
}

type DefaultServiceTable struct {
	db       dbm.DB
	services map[string]servicetypes.ServiceData
}

func NewServiceTable() ServiceTable {
	serviceTable := &DefaultServiceTable{
		services: map[string]servicetypes.ServiceData{},
	}
	return serviceTable
}

func (st *DefaultServiceTable) GetServices() ([]servicetypes.ServiceData, error) {
	var services []servicetypes.ServiceData
	for _, service := range st.services {
		services = append(services, service)
	}
	return services, nil
}

func (st *DefaultServiceTable) GetServiceByType(requiredType string) (servicetypes.ServiceData, error) {
	for serviceType, service := range st.services {
		if serviceType == requiredType {
			return service, nil
		}
	}
	return nil, fmt.Errorf("type %s not found", requiredType)
}

func (st *DefaultServiceTable) UseService(serviceData servicetypes.ServiceData) error {
	st.services[serviceData.GetType()] = serviceData
	return nil
}
