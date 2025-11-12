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

All API handlers can be found under `internal/handler/core_handler.go`.

## Helpers

All core helpers, for example calculating the averages can be found under `internal/services`.

## Packages

Any package for third party dependencies can be found in `internal/pkg`. For example in `logger.go` we have written the functions that will initialize the zap logger.

## Sample CI/CD

Under `.github/workflows` you can find a sample github actions workflow that just builds the docker image using the Docker file found in the root of the repo. 

This workflow can be extended to push the image up to a remote registry like Artifact Storage on GCP. Additionally, on a successful push Slack can be notified, so the rest of the team knows that the image is ready to use.

