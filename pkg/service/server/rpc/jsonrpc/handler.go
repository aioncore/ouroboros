package jsonrpc

import (
	"encoding/json"
	"fmt"
	"github.com/aioncore/ouroboros/pkg/service/log"
	"github.com/aioncore/ouroboros/pkg/service/server/rpc/jsonrpc/types"
	"github.com/aioncore/ouroboros/pkg/service/server/utils"
	"io"
	"net/http"
	"reflect"
)

func MakeJSONRPCHandler(funcMap map[string]*utils.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			res := types.RPCInvalidRequestError(nil,
				fmt.Errorf("error reading request body: %w", err),
			)
			if wErr := WriteRPCResponseHTTPError(w, http.StatusBadRequest, res); wErr != nil {
				log.Error("failed to write response")
			}
			return
		}

		// first try to unmarshal the incoming request as an array of RPC requests
		var (
			requests  []types.RPCRequest
			responses []types.RPCResponse
		)
		if err := json.Unmarshal(b, &requests); err != nil {
			// next, try to unmarshal as a single request
			var request types.RPCRequest
			if err := json.Unmarshal(b, &request); err != nil {
				res := types.RPCParseError(fmt.Errorf("error unmarshaling request: %w", err))
				if wErr := WriteRPCResponseHTTPError(w, http.StatusInternalServerError, res); wErr != nil {
					log.Error("failed to write response")
				}
				return
			}
			requests = []types.RPCRequest{request}
		}

		// Set the default response cache to true unless
		// 1. Any RPC request error.
		// 2. Any RPC request doesn't allow to be cached.
		// 3. Any RPC request has the height argument and the value is 0 (the default).
		for _, request := range requests {

			// A Notification is a Request object without an "id" member.
			// The Server MUST NOT reply to a Notification, including those that are within a batch request.
			if request.ID == nil {
				log.Debug(
					"HTTPJSONRPC received a notification, skipping... (please send a non-empty ID if you want to call a method)")
				continue
			}
			coreFunc, ok := funcMap[request.Method]
			if !ok {
				responses = append(responses, types.RPCMethodNotFoundError(request.ID))
				continue
			}
			ctx := &types.Context{JSONReq: &request, HTTPReq: r}
			args := []reflect.Value{reflect.ValueOf(ctx)}
			if len(request.Params) > 0 {
				fnArgs, err := jsonParamsToArgs(coreFunc, request.Params)
				if err != nil {
					responses = append(
						responses,
						types.RPCInvalidParamsError(request.ID, fmt.Errorf("error converting json params to arguments: %w", err)),
					)
					continue
				}
				args = append(args, fnArgs...)
			}

			returns := coreFunc.F.Call(args)

			if returns[len(returns)-1].Interface() != nil {
				responses = append(responses, types.RPCInternalError(request.ID, fmt.Errorf("%v", returns[len(returns)-1].Interface())))
				continue
			}
			result := reflect.New(returns[0].Type())
			result.Elem().Set(returns[0])
			responses = append(responses, types.NewRPCSuccessResponse(request.ID, result.Interface()))
		}

		if len(responses) > 0 {
			err = WriteRPCResponseHTTP(w, []httpHeader{}, responses...)
			if err != nil {
				log.Error("failed to write responses")
			}
		}
	}
}

// WriteRPCResponseHTTPError marshals res as JSON (with indent) and writes it
// to w.
//
// source: https://www.jsonrpc.org/historical/json-rpc-over-http.html
func WriteRPCResponseHTTPError(
	w http.ResponseWriter,
	httpCode int,
	res types.RPCResponse,
) error {
	if res.Error == nil {
		panic("tried to write http error response without RPC error")
	}

	jsonBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_, err = w.Write(jsonBytes)
	return err
}

// raw is unparsed json (from json.RawMessage) encoding either a map or an
// array.
//
// Example:
//
//	rpcFunc.args = [rpctypes.Context string]
//	rpcFunc.argNames = ["arg"]
func jsonParamsToArgs(rpcFunc *utils.APIFunc, raw []byte) ([]reflect.Value, error) {
	const argsOffset = 1

	// TODO: Make more efficient, perhaps by checking the first character for '{' or '['?
	// First, try to get the map.
	var m map[string]json.RawMessage
	err := json.Unmarshal(raw, &m)
	if err == nil {
		return mapParamsToArgs(rpcFunc, m, argsOffset)
	}

	// Otherwise, try an array.
	var a []json.RawMessage
	err = json.Unmarshal(raw, &a)
	if err == nil {
		return arrayParamsToArgs(rpcFunc, a, argsOffset)
	}

	// Otherwise, bad format, we cannot parse
	return nil, fmt.Errorf("unknown type for JSON params: %v. Expected map or array", err)
}

func mapParamsToArgs(
	rpcFunc *utils.APIFunc,
	params map[string]json.RawMessage,
	argsOffset int,
) ([]reflect.Value, error) {
	values := make([]reflect.Value, len(rpcFunc.ArgNames))
	for i, argName := range rpcFunc.ArgNames {
		argType := rpcFunc.Args[i+argsOffset]

		if p, ok := params[argName]; ok && p != nil && len(p) > 0 {
			val := reflect.New(argType)
			err := json.Unmarshal(p, val.Interface())
			if err != nil {
				return nil, err
			}
			values[i] = val.Elem()
		} else { // use default for that type
			values[i] = reflect.Zero(argType)
		}
	}

	return values, nil
}

func arrayParamsToArgs(
	rpcFunc *utils.APIFunc,
	params []json.RawMessage,
	argsOffset int,
) ([]reflect.Value, error) {
	if len(rpcFunc.ArgNames) != len(params) {
		return nil, fmt.Errorf("expected %v parameters (%v), got %v (%v)",
			len(rpcFunc.ArgNames), rpcFunc.ArgNames, len(params), params)
	}

	values := make([]reflect.Value, len(params))
	for i, p := range params {
		argType := rpcFunc.Args[i+argsOffset]
		val := reflect.New(argType)
		err := json.Unmarshal(p, val.Interface())
		if err != nil {
			return nil, err
		}
		values[i] = val.Elem()
	}
	return values, nil
}

type httpHeader struct {
	name  string
	value string
}

func WriteRPCResponseHTTP(w http.ResponseWriter, headers []httpHeader, res ...types.RPCResponse) error {
	var v interface{}
	if len(res) == 1 {
		v = res[0]
	} else {
		v = res
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	w.Header().Set("Content-Type", "application/json")
	for _, header := range headers {
		w.Header().Set(header.name, header.value)
	}
	w.WriteHeader(200)
	_, err = w.Write(jsonBytes)
	return err
}
