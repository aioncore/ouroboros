package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/aioncore/ouroboros/pkg/service/client/rpc/jsonrpc"
	jsonrpctypes "github.com/aioncore/ouroboros/pkg/service/server/rpc/jsonrpc/types"
	"github.com/aioncore/ouroboros/pkg/service/server/utils"
	"github.com/aioncore/ouroboros/pkg/service/types/services"
	"net/http"
	"strings"
)

type Switch interface {
	Services(ctx *jsonrpctypes.Context) ([]services.ServiceData, error)
	UseService(ctx *jsonrpctypes.Context, rawServiceData json.RawMessage) error
	CallService(ctx *jsonrpctypes.Context, serviceType string, r *jsonrpctypes.RPCRequest) (*http.Response, error)
}

type DefaultSwitch struct {
	handlers     map[string]*utils.APIFunc
	serviceTable ServiceTable
}

func NewSwitch() Switch {
	sw := &DefaultSwitch{
		serviceTable: NewServiceTable(),
	}
	return sw
}

func (sw *DefaultSwitch) Services(ctx *jsonrpctypes.Context) ([]services.ServiceData, error) {
	return sw.serviceTable.GetServices()
}

func (sw *DefaultSwitch) UseService(ctx *jsonrpctypes.Context, rawServiceData json.RawMessage) error {
	var serviceData services.ServiceData
	serviceType, err := sw.parseType(rawServiceData)
	if err != nil {
		return err
	}
	switch serviceType {
	case "p2p":
		serviceData = &services.P2PServiceData{}
		err = json.Unmarshal(rawServiceData, serviceData)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("type %s not found", serviceType)
	}
	return sw.serviceTable.UseService(serviceData)
}

func (sw *DefaultSwitch) CallService(ctx *jsonrpctypes.Context, serviceType string, r *jsonrpctypes.RPCRequest) (*http.Response, error) {
	serviceData, err := sw.serviceTable.GetServiceByType(serviceType)
	if err != nil {
		return nil, err
	}
	addr := strings.Replace(serviceData.GetAddress(), "tcp://", "http://", 1)
	rpcClient := jsonrpc.NewRPCClient(addr)
	return rpcClient.Post(r)
}

func (sw *DefaultSwitch) parseType(rawServiceData json.RawMessage) (string, error) {
	var dataMap map[string]json.RawMessage
	err := json.Unmarshal(rawServiceData, &dataMap)
	if err != nil {
		return "", fmt.Errorf("json format error")
	}
	rawHeader, ok := dataMap["header"]
	if !ok {
		return "", fmt.Errorf("json format error")
	}
	var headerMap map[string]json.RawMessage
	err = json.Unmarshal(rawHeader, &headerMap)
	if err != nil {
		return "", fmt.Errorf("json format error")
	}
	rawServiceType, ok := headerMap["type"]
	if !ok {
		return "", fmt.Errorf("json format error")
	}
	serviceType := ""
	err = json.Unmarshal(rawServiceType, &serviceType)
	if err != nil {
		return "", err
	}
	return serviceType, nil
}
