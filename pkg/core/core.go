package core

import (
	"github.com/aioncore/ouroboros/pkg/config"
	"github.com/aioncore/ouroboros/pkg/core/middleware"
	"github.com/aioncore/ouroboros/pkg/service"
	"github.com/aioncore/ouroboros/pkg/service/crypto"
	"github.com/aioncore/ouroboros/pkg/service/log"
	httphandler "github.com/aioncore/ouroboros/pkg/service/server/http"
	"github.com/aioncore/ouroboros/pkg/service/server/rpc/jsonrpc"
	"net"
	"net/http"
	"strings"
)

type Core struct {
	service.BaseService
	serverListeners []net.Listener
	sw              middleware.Switch
	nodeKey         *crypto.NodeKey
}

// NewCore returns a new, ready to go, Core.
func NewCore() *Core {

	core := &Core{
		sw: middleware.NewSwitch(),
	}
	core.BaseService = *service.NewBaseService("core", core)
	nodeKey, err := crypto.LoadOrGenerateNodeKey(config.Core.NodeKeyPath)
	if err != nil {
		panic(err)
	}
	core.nodeKey = nodeKey
	return core
}

func (c *Core) OnStart() error {

	listeners, err := c.startServer()
	if err != nil {
		return err
	}
	c.serverListeners = listeners
	return nil
}

func (c *Core) startServer() ([]net.Listener, error) {
	listenAddresses := []string{"tcp://127.0.0.1:26657"}
	listeners := make([]net.Listener, len(listenAddresses))
	for i, listenAddr := range listenAddresses {
		mux := http.NewServeMux()
		coreRoutes := GenerateRoutes(c.sw)

		// HTTP endpoints
		for funcName, coreFunc := range coreRoutes {
			mux.HandleFunc("/"+funcName, httphandler.MakeHTTPHandler(coreFunc))
		}

		// JSON RPC endpoints
		mux.HandleFunc("/rpc", jsonrpc.MakeJSONRPCHandler(coreRoutes))

		parts := strings.SplitN(listenAddr, "://", 2)
		proto, addr := parts[0], parts[1]
		listener, err := net.Listen(proto, addr)
		if err != nil {
			return nil, err
		}

		go func() {
			s := &http.Server{
				Handler: mux,
			}
			err := s.Serve(listener)
			if err != nil {
				return
			}
		}()

		listeners[i] = listener
	}
	return listeners, nil
}

func (c *Core) OnStop() {
	for _, l := range c.serverListeners {
		log.Info("Closing rpc listener")
		if err := l.Close(); err != nil {
			log.Error("Error closing listener")
		}
	}
}
