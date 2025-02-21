document.addEventListener('DOMContentLoaded', async function() {
    // Initialize the map centered at 0,0 (the middle of the world) with zoom level 1
    const map = L.map('artist-map', {
        center: [20, 0],  // Slightly above equator for better perspective
        zoom: 2,
        minZoom: 2,      // Prevent zooming out too far
        maxBounds: [      // Restrict panning to reasonable bounds
            [-90, -180],  // Southwest corner
            [90, 180]     // Northeast corner
        ]
    });
    
    // Add OpenStreetMap tiles
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors',
        noWrap: true     // Prevents the map from repeating horizontally
    }).addTo(map);

    // Get artist ID from URL
    const urlParams = new URLSearchParams(window.location.search);
    const artistId = urlParams.get('id');

    // Create a marker group
    const markerGroup = L.featureGroup();

    try {
        // Fetch coordinates from our backend
        const response = await fetch(`/api/coordinates?id=${artistId}`);
        const coordinates = await response.json();

        // Add markers for each location
        coordinates.forEach(coord => {
            const marker = L.marker([coord.lat, coord.lon])
                .bindPopup(`<b>${coord.address}</b>`)
                .addTo(map);
            markerGroup.addLayer(marker);
        });

        // Only adjust bounds if we have markers
        if (markerGroup.getLayers().length > 0) {
            markerGroup.addTo(map);
            map.fitBounds(markerGroup.getBounds(), {
                padding: [50, 50],
                maxZoom: 5  // Don't zoom in too far when fitting bounds
            });
        }
    } catch (error) {
        console.error('Error fetching coordinates:', error);
    }
});