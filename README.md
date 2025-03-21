# go-serverless platform

This is a serverless platform built using Go and Docker. It allows users to upload a ZIP file containing either Python or Go code, which is then compiled and executed in a Docker container.

## Endpoints

### 1. Submit Code
- **Endpoint**: `POST /api/submit`
- **Description**: Accepts a ZIP file containing the user's code.
- **Request**:
  - **Parameter**: `code` (file)
  - **Example**: 
    ```bash
    curl -X POST -F "code=@$(zip_file)" http://0.0.0.0:8080/api/submit
    ```

### 2. Execute Code
- **Endpoint**: `GET /api/execute`
- **Description**: Executes the previously submitted code.
- **Query Parameter**:
  - `functionID`: The ID of the function to be executed.
- **Example**:
  ```bash
  curl http://0.0.0.0:8080/api/execute\?functionID\=$(functionID)

## Architecture

The platform consists of the following components:

1. **Go Server**: The main application that handles the incoming requests, extracts the code from the ZIP file, and compiles/executes it in a Docker container.
2. **Docker**: Used to create and manage the runtime environment for the user's code.
3. **YAML Configuration**: Defines the Docker image and build instructions as per programming language.

## Ensure Installed

- **Docker**: [https://www.docker.com/get-started](https://www.docker.com/get-started)
- **Go**: [https://golang.org/dl/](https://golang.org/dl/)

## Server Deployment

1. Clone the repository:
   ```bash
   git clone https://github.com/hardikkum444/go-serverless.git
   cd go-serverless

2. Build the binary:
   ```bash
    go build main.go

3. Run the application:
   ```bash
    chmod +x main
    ./main

4. Submit a ZIP file:
   ```bash
    make submit zip_file=<zip_file_name>

4. Execute the submitted code:
   ```bash
    make execute functionID=<functionID>

Make sure to replace <IP> in the Makefile with the correct domain or IP address of your server.

Make sure to replace <IP> in the Makefile with the correct domain or IP address of your server.

## Grafana and Prometheus Dashboard Setup

Below are the images of the Grafana and Prometheus dashboard setup for the go-serverless platform:

### Grafana Dashboard
![Grafana Dashboard](assets/go-serverless-graphana.png)

### Prometheus Dashboard
![Prometheus Dashboard](assets/go-serverless-prometheus.png)
