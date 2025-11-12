# edge-device-sample

## Repository Structure

```
.
├── cmd/                     # Application entrypoints
│   └── main.go              
├── docker-compose.yml       
├── etc/                     # API Contracts/fixture data
│   ├── devices.csv          
│   └── openapi.json         
├── go.mod                   
├── go.sum                   
├── internal/                # Application code
│   ├── data/                # Storage layer
│   ├── handler/             # HTTP handlers
│   ├── http/
│   ├── pkg/
│   ├── services/
├── Makefile                 # Makefile to 
├── results.txt              # Simulator outputs
├── sim/                     # Device simulator binary
│   └── device-simulator-mac-amd64
```

## Entrypoint 

The entrypoint to the application can be found at `cmd/main.go`. This file initializes the packages required to actually run the api as well as seed the in memory store on startup with the device ids found in `etc/devices.csv`. 

## Routing

The core router is what outlines the API routes that can be interacted with by the simulator or even tools like cURL and/or Postman. 

The routes can be found in `internal/http/router.go`. The `http` directory holds all things related to api routing and middleware. 

