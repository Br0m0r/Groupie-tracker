# Groupie Tracker Map Feature

## Overview
The map feature in Groupie Tracker visualizes concert locations for artists using OpenStreetMap and Leaflet.js. It provides an interactive way to explore where artists have performed or will perform.

## Features
- Interactive world map showing concert locations
- Clickable markers with location information
- Automatic zoom and centering based on concert locations
- Responsive design that adapts to different screen sizes
- Geocoding of location addresses to coordinates

## Technical Implementation

### Backend Components

#### Coordinates Handler (`coordinates.go`)
- Endpoint: `/api/coordinates`
- Converts location addresses to geographic coordinates using Nominatim API
- Implements rate limiting to respect API usage policies
- Returns JSON array of coordinates with location information

```go
type Coordinates struct {
    Lat     float64 `json:"lat"`
    Lon     float64 `json:"lon"`
    Address string  `json:"address"`
}
```

### Frontend Components

#### Map Initialization (`artist-map.js`)
```javascript
const map = L.map('artist-map', {
    center: [20, 0],
    zoom: 2,
    minZoom: 2,
    maxBounds: [
        [-90, -180],
        [90, 180]
    ]
});
```

#### Map Container (`artist.html`)
```html
<div id="artist-map"></div>
```

#### Styling (`artist.css`)
```css
#artist-map {
    width: 900px;
    height: 550px;
    border-radius: 1rem;
    overflow: hidden;
    box-shadow: 0 5px 15px rgba(0,0,0,0.1);
}
```

## Dependencies
- Leaflet.js: v1.9.3
- OpenStreetMap Tiles
- Nominatim Geocoding API

## Setup Instructions

1. Ensure all dependencies are included:
```html
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.3/dist/leaflet.css" />
<script src="https://unpkg.com/leaflet@1.9.3/dist/leaflet.js"></script>
```

2. Create a map container in your HTML:
```html
<div id="artist-map"></div>
```

3. Initialize the map in your JavaScript:
```javascript
document.addEventListener('DOMContentLoaded', async function() {
    const map = L.map('artist-map', {
        center: [20, 0],
        zoom: 2
    });
    
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors'
    }).addTo(map);
});
```

## API Usage

### Getting Coordinates
```bash
GET /api/coordinates?id={artistId}
```

Response format:
```json
[
    {
        "lat": 49.59380,
        "lon": 8.15052,
        "address": "Germany Mainz"
    }
]
```

## Error Handling
- Invalid artist IDs return a 400 Bad Request
- Artist not found returns a 404 Not Found
- Server errors return a 500 Internal Server Error

## Responsive Design
The map automatically adjusts its size based on screen width:
```css
@media (max-width: 768px) {
    #artist-map {
        height: 400px;
    }
}

@media (max-width: 480px) {
    #artist-map {
        height: 300px;
    }
}
```

## Best Practices
1. Always implement rate limiting for geocoding requests
2. Cache coordinates when possible to reduce API calls
3. Handle map bounds to prevent excessive scrolling
4. Provide fallback behavior when geocoding fails
5. Ensure proper error handling and user feedback

## Known Limitations
- Nominatim API has usage limits (1 request per second)
- Some addresses may not resolve to exact coordinates
- Map interactions may be limited on mobile devices

## Future Improvements
- Implement coordinate caching
- Add clustering for multiple nearby locations
- Enhance marker information windows
- Add route visualization between venues
- Implement custom map styles