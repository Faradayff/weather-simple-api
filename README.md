# Weather Simple API

## Overview

Weather Simple API is a Go project that provides a simple HTTP API to fetch and aggregate weather forecasts from multiple providers.

## Requirements

- [Go 1.18+](https://golang.org/dl/)
- Internet connection (to access weather APIs)
- (Optional) [Git](https://git-scm.com/) for cloning the repository

## Setup

1. **Clone the repository:**

   ``` sh
   git clone https://github.com/faradayff/weather-simple-api.git
   cd weather-simple-api
   ```

2. **Set up API keys:**
   - Create a `.env` file in the project root with your API key(s):

     ``` text
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

## Usage

Make a GET request to `/weather` with `lat` and `lon` query parameters:

``` text
GET http://localhost:8080/weather?lat=40.7128&lon=-74.0060
```

You will receive a JSON response with aggregated weather forecasts from the configured providers.

## Notes

- Ensure your `.env` file is **not** committed to version control.
- The API may be subject to rate limits from external providers.
- For development, you can use tools like [curl](https://curl.se/) or [Postman](https://www.postman.com/) to test the endpoint.

## License

This project is for educational purposes.
