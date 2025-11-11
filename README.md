# edge-device-sample

## Repository Structure

```
.
├── cmd/                     # Application entrypoints
│   └── main.go              
├── docker-compose.yml       
├── etc/                     # API Contracts/fixture data
│   ├── devices.csv          # Seed device IDs loaded into in-memory storage
│   └── openapi.json         # API contract describing available endpoints
├── go.mod                   
├── go.sum                   
├── internal/                # Application code
│   ├── data/                # Storage layer
│   │   ├── in_memory_storage.go
│   │   └── types.go
│   ├── handler/             # HTTP handlers
│   │   ├── core_handler.go
│   │   ├── errors.go
│   │   └── types.go
├── Makefile                 # Makefile to 
├── results.txt              # Simulator outputs
├── sim/                     # Device simulator binary
│   └── device-simulator-mac-amd64
```