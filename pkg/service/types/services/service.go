package services

import "fmt"

type ServiceData interface {
	GetHash() string
	GetAddress() string
	GetType() string
	String() string
}

type ServiceHeader struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	SHA256     string `json:"sha256"`
	RPCAddress string `json:"rpc_address"`
}

func (s *ServiceHeader) GetHash() string {
	return s.SHA256
}

func (s *ServiceHeader) GetAddress() string {
	return s.RPCAddress
}

func (s *ServiceHeader) GetType() string {
	return s.Type
}

func (s *ServiceHeader) String() string {
	return fmt.Sprintf("%s service %s", s.Type, s.Name)
}
