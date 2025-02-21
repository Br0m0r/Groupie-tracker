document.addEventListener('DOMContentLoaded', async function() {
    console.log('Map initialization started');
    // Initialize the map
    const map = L.map('artist-map').setView([0, 0], 2);
    
    // Add OpenStreetMap tiles
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors'
    }).addTo(map);

    // Get artist ID from URL
    const urlParams = new URLSearchParams(window.location.search);
    const artistId = urlParams.get('id');

    // Create a marker group
    const markerGroup = L.featureGroup();

    try {
        // Fetch coordinates from our backend
        console.log('Fetching coordinates for artist:', artistId);
        const response = await fetch(`/api/coordinates?id=${artistId}`);
        const coordinates = await response.json();
        console.log('Received coordinates:', coordinates);

        // Add markers for each location
        coordinates.forEach(coord => {
            const marker = L.marker([coord.lat, coord.lon])
                .bindPopup(`<b>${coord.address}</b>`)
                .addTo(map);
            markerGroup.addLayer(marker);
        });

        // Fit map to show all markers
        if (markerGroup.getLayers().length > 0) {
            markerGroup.addTo(map);
            map.fitBounds(markerGroup.getBounds(), {
                padding: [50, 50]
            });
        }
    } catch (error) {
        console.error('Error fetching coordinates:', error);
    }
});