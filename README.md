# Groupie Tracker

A web application built with Go that visualizes music artist data, concert locations, and tour dates using the [Groupie Trackers API](https://groupietrackers.herokuapp.com/api). Concert locations are displayed on an interactive map powered by Leaflet and OpenStreetMap.

## Features

| Feature | Description |
|---------|-------------|
| Artist Directory | Browse all artists with detail pages |
| Search | Live suggestions across artists, members, locations, and dates |
| Filters | Filter by creation date, first album year, member count, and location |
| Concert Map | Interactive map with geocoded concert locations via Nominatim |

## Requirements

- Go 1.22 or later
- Internet access (API data, map tiles, geocoding)

## Getting Started

### Manual

```bash
go run main.go
```

The server starts at `http://localhost:8080`.

### Docker

```dockerfile
FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go build -o groupie .
EXPOSE 8080
CMD ["./groupie"]
```

```bash
docker build -t groupie-tracker .
docker run -p 8080:8080 groupie-tracker
```

## Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Home page with artist grid |
| GET | `/artist?id={id}` | Artist detail page |
| GET | `/search?q={query}` | Search results (HTML) or suggestions (JSON via XHR) |
| GET | `/filter` | Filtered artist results |
| GET | `/api/coordinates?id={id}` | Concert location coordinates (JSON) |

## Project Structure

```
main.go              Entry point and server configuration
handlers/            HTTP request handlers
models/              Data structures
store/               Data fetching, caching, and storage
utils/               Formatting and helper functions
templates/           HTML templates
static/              CSS and JavaScript assets
```
