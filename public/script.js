document.addEventListener("DOMContentLoaded", async function () {
    var map = L.map('map').setView([51.1657, 10.4515], 5); // Centered in Germany

    // Load OpenStreetMap tiles
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    try {
        let response = await fetch('/api/concerts');
        let concerts = await response.json();

        // Group concerts by coordinates (key: "lat,lon")
        let groupedConcerts = {};
        concerts.forEach(concert => {
            if (concert.lat && concert.lon) {
                let key = `${concert.lat},${concert.lon}`;
                if (groupedConcerts[key]) {
                    groupedConcerts[key].artists.push(concert.artist);
                } else {
                    groupedConcerts[key] = {
                        artists: [concert.artist],
                        location: concert.location,
                        lat: concert.lat,
                        lon: concert.lon
                    };
                }
            }
        });

        // Add a single marker per group
        for (let key in groupedConcerts) {
            let group = groupedConcerts[key];
            let popupContent = `<b>${group.artists.join(", ")}</b><br>${group.location}`;
            L.marker([group.lat, group.lon]).addTo(map)
                .bindPopup(popupContent);
        }
    } catch (error) {
        console.error("Error loading concert data:", error);
    }
});
