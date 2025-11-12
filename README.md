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
│   ├── helpers/            # Helper functions
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

All core helpers, for example calculating the averages can be found under `internal/helpers`.

## Packages

Any package for third party dependencies can be found in `internal/pkg`. For example in `logger.go` we have written the functions that will initialize the zap logger.

## Data 

For this task, we used a memory store to keep track of all of the data. However, in a real customer environment this can easily be swapped out to use a relational or noSQL database. 

## Sample CI/CD

Under `.github/workflows` you can find a sample github actions workflow that just builds the docker image using the Docker file found in the root of the repo. 

This workflow can be extended to push the image up to a remote registry like Artifact Storage on GCP. Additionally, on a successful push Slack can be notified, so the rest of the team knows that the image is ready to use.

## Testing

Tests can be found in `/tests/*`. For now they are just unit tests, verifying the functionality of our helpers. 

## Running this project

To Run this project, start off by spinning up your core backend service using docker compose. If you don't remember the exact commands to run, the `Makefile` makes it easy for you. 

To spin up your service: 

```
make build-and-run-docker
```

One the server is running, you need to run the simulator `tests/services/core_test.go`. 

```
make run-sim
```

After the simulator runs, it will spin out a `results.txt` in your root repo, which contents resembling: 

```
[device-simulator] 2025/11/11 16:20:23 finished loading stats requests
[device-simulator] 2025/11/11 16:20:23 finished loading heartbeat requests
[device-simulator] 2025/11/11 16:20:23 starting requests for device_id: 60-6b-44-84-dc-64 
[device-simulator] 2025/11/11 16:20:23 starting requests for device_id: b4-45-52-a2-f1-3c 
[device-simulator] 2025/11/11 16:20:23 starting requests for device_id: 26-9a-66-01-33-83 
[device-simulator] 2025/11/11 16:20:23 starting requests for device_id: 18-b8-87-e7-1f-06 
[device-simulator] 2025/11/11 16:20:23 starting requests for device_id: 38-4e-73-e0-33-59 
[device-simulator] 2025/11/11 16:20:23 executing 100 stats and 446 heartbeats for 26-9a-66-01-33-83
[device-simulator] 2025/11/11 16:20:23 executing 100 stats and 479 heartbeats for 38-4e-73-e0-33-59
[device-simulator] 2025/11/11 16:20:23 executing 100 stats and 480 heartbeats for b4-45-52-a2-f1-3c
[device-simulator] 2025/11/11 16:20:23 executing 100 stats and 474 heartbeats for 18-b8-87-e7-1f-06
[device-simulator] 2025/11/11 16:20:23 executing 100 stats and 479 heartbeats for 60-6b-44-84-dc-64
[device-simulator] 2025/11/11 16:20:27 ############ RESULTS #################
[device-simulator] 2025/11/11 16:20:27 
DeviceID: 60-6b-44-84-dc-64
	Uptime
		Expected: 99.79167
		Actual: 99.58420

	AvgUploadTime
		Expected: 3m7.893379134s
		Actual: 3m7.893379134s
[device-simulator] 2025/11/11 16:20:27 
DeviceID: b4-45-52-a2-f1-3c
	Uptime
		Expected: 100.00000
		Actual: 99.79210

	AvgUploadTime
		Expected: 3m19.085533836s
		Actual: 3m19.085533836s
[device-simulator] 2025/11/11 16:20:27 
DeviceID: 26-9a-66-01-33-83
	Uptime
		Expected: 92.91667
		Actual: 92.72349

	AvgUploadTime
		Expected: 3m21.858747766s
		Actual: 3m21.858747766s
[device-simulator] 2025/11/11 16:20:27 
DeviceID: 18-b8-87-e7-1f-06
	Uptime
		Expected: 98.75000
		Actual: 98.54470

	AvgUploadTime
		Expected: 3m17.331667813s
		Actual: 3m17.331667813s
[device-simulator] 2025/11/11 16:20:27 
DeviceID: 38-4e-73-e0-33-59
	Uptime
		Expected: 99.79167
		Actual: 99.58420

	AvgUploadTime
		Expected: 3m29.226522788s
		Actual: 3m29.226522788s
[device-simulator] 2025/11/11 16:20:27 all done!
```

To bring your services back down, run the following command:

```
make stop-docker
```

## Potential Future considerations
- Add an actual database. 
- Add more logic to the Github Actions workflow and ensure that tests are run before we builf the image.