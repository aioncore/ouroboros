# Services

## JSONRPC

### POST 127.0.0.1:26657/rpc

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "services"
}
```

### Response

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "status": "success",
    "reason": "success",
    "services": [
      {
      "type": "p2p",
      "name": "libp2p",
      "sha256": ""
      }
    ]
  }
}
```

# UseService

## JSONRPC

### POST 127.0.0.1:26657/rpc

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "use_service",
  "params": {
    "service_data": {
      "header": {
        "name": "libp2p",
        "type": "p2p",
        "sha256": "E3AD6ABCCD8B4F5E63B3D3C9380D084F6B28C962AA724798C7E73F1A6966766B",
        "rpc_address": "tcp://127.0.0.1:26658"
      },
      "body": {

      }
    }
  }
}
```

### Response

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "status": "success",
    "reason": "success"
  }
}
```

# CallService

## JSONRPC

### POST 127.0.0.1:26657/rpc

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "call_service",
  "params": {
    "type": "p2p",
    "request": {
      "jsonrpc": "2.0",
      "id": "2",
      "method": "send_transaction",
      "params": {
        "from": "a",
        "to": "b"
      } 
    }
  }
}
```

### Response

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "status": "success",
    "reason": "success"
  }
}
```
