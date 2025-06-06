# Weather Simple API

## Overview

Weather Simple API is a Go project that provides a simple HTTP API to fetch and aggregate weather forecasts from multiple providers. It includes features like worker-based task processing, and support for multiple weather APIs.

## Requirements

- [Go 1.18+](https://golang.org/dl/)
- Internet connection (to access weather APIs)
- (Optional) [Git](https://git-scm.com/) for cloning the repository

## Features

- Fetch weather forecasts from multiple providers.
- Worker-based task processing for efficient handling of API requests.
- Environment variable support using [github.com/joho/godotenv](https://github.com/joho/godotenv).

## Setup

1. **Clone the repository:**

   ```sh
   git clone https://github.com/faradayff/weather-simple-api.git
   cd weather-simple-api
   ```

2. **Set up API keys:**
   - Create a `.env` file in the project root with your API key(s):

     ```text
     WEATHER_API_KEY=your_weatherapi_key
     ```

   - The project uses [github.com/joho/godotenv](https://github.com/joho/godotenv) to load environment variables.

3. **Install dependencies:**

   ```sh
   go mod tidy
   ```

## Running the Application

Start the HTTP server with:

```sh
go run cmd/server.go
```

The server will listen on `http://localhost:8080`.

## Running with Docker

This project is Docker-ready, allowing you to run the application in a containerized environment.

### Build the Docker Image

To build the Docker image, run the following command in the root of the project:

```sh
docker build -t weather-simple-api .
```

### Run the Docker Container

Once the image is built, you can run the container with:

```sh
docker run --env-file .env weather-simple-api
```

- The `--env-file .env` flag ensures that the container has access to the required environment variables.

### Access the API

After starting the container, the API will be available at:

```text
http://localhost:8080/weather?lat=40.7128&lon=-74.0060
```

## Usage

Make a GET request to `/weather` with `lat` and `lon` query parameters:

```text
GET http://localhost:8080/weather?lat=40.7128&lon=-74.0060
```

You will receive a JSON response with aggregated weather forecasts from the configured providers.

## Notes

- Make sure the `.env` file is correctly configured before running the container and is **not** committed to version control.
- The API may be subject to rate limits from external providers.
- For development, you can use tools like [curl](https://curl.se/) or [Postman](https://www.postman.com/) to test the endpoint.

## Testing

Run the tests with:

```sh
go test ./...
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is for educational purposes.
