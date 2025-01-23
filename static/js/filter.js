document.addEventListener('DOMContentLoaded', () => {
    // Initialize range sliders
    const updateRange = (minEl, maxEl, displayEl) => {
        const minVal = parseInt(minEl.value);
        const maxVal = parseInt(maxEl.value);
        
        if (minVal > maxVal) {
            if (minEl === document.activeElement) {
                maxEl.value = minVal;
            } else {
                minEl.value = maxVal;
            }
        }
        
        displayEl.textContent = `${minEl.value}-${maxEl.value}`;
    };

    ['creation', 'album'].forEach(type => {
        const min = document.getElementById(`${type}YearMin`);
        const max = document.getElementById(`${type}YearMax`);
        const display = document.getElementById(`${type}YearDisplay`);
        
        [min, max].forEach(el => {
            el.addEventListener('input', () => updateRange(min, max, display));
        });
    });

    // Locations toggle
    const locationsToggle = document.getElementById('locationsToggle');
    const locationsList = document.getElementById('locationsList');
    const toggleIcon = locationsToggle.querySelector('.toggle-icon');
    
    locationsToggle.addEventListener('click', () => {
        const isHidden = locationsList.style.display === 'none';
        locationsList.style.display = isHidden ? 'block' : 'none';
        toggleIcon.style.transform = isHidden ? 'rotate(180deg)' : 'rotate(0)';
        toggleIcon.textContent = isHidden ? '▲' : '▼';
    });

    // Sample locations - replace with data from backend
    const locations = ["New York", "London", "Paris", "Tokyo", "Berlin"];
    locationsList.innerHTML = locations.map(loc => 
        `<div class="location-item">
            <label>
                <input type="checkbox" value="${loc}">
                ${loc}
            </label>
        </div>`
    ).join('');
    
    // Apply filters
    document.getElementById('applyFilters').addEventListener('click', () => {
        const filters = {
            creation: {
                min: parseInt(document.getElementById('creationYearMin').value),
                max: parseInt(document.getElementById('creationYearMax').value)
            },
            album: {
                min: parseInt(document.getElementById('albumYearMin').value),
                max: parseInt(document.getElementById('albumYearMax').value)
            },
            members: [...document.querySelectorAll('.members-grid input:checked')]
                .map(input => parseInt(input.value)),
            locations: [...document.querySelectorAll('.locations-list input:checked')]
                .map(input => input.value)
        };

        fetch('/filter', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(filters)
        })
        .then(response => response.json())
        .then(data => {
            const artistsGrid = document.querySelector('.artists-grid');
            artistsGrid.innerHTML = data.map(artist => `
                <div class="artist-card">
                    <a href="/artist?id=${artist.ID}" class="artist-link">
                        <div class="image-container">
                            <img src="${artist.Image}" alt="${artist.Name}" loading="lazy">
                        </div>
                        <div class="artist-info">
                            <h2>${artist.Name}</h2>
                        </div>
                    </a>
                </div>
            `).join('');
        })
        .catch(error => console.error('Error:', error));
    });
});