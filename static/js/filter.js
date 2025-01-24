document.addEventListener('DOMContentLoaded', () => {
    // Range sliders update
    ['creation', 'album'].forEach(type => {
        const min = document.getElementById(`${type}YearMin`);
        const max = document.getElementById(`${type}YearMax`);
        const display = document.getElementById(`${type}YearDisplay`);
        
        [min, max].forEach(el => {
            el.addEventListener('input', () => {
                const minVal = Math.min(parseInt(min.value), parseInt(max.value));
                const maxVal = Math.max(parseInt(min.value), parseInt(max.value));
                min.value = minVal;
                max.value = maxVal;
                display.textContent = `${minVal}-${maxVal}`;
            });
        });
    });

    // Locations toggle
    document.getElementById('locationsToggle')?.addEventListener('click', () => {
        const list = document.getElementById('locationsList');
        const icon = document.querySelector('.toggle-icon');
        list.style.display = list.style.display === 'none' ? 'block' : 'none';
        icon.textContent = list.style.display === 'none' ? '▼' : '▲';
    });

    // Apply filters
    document.getElementById('applyFilters')?.addEventListener('click', () => {
        const filters = {
            creation: {
                min: parseInt(document.getElementById('creationYearMin').value),
                max: parseInt(document.getElementById('creationYearMax').value)
            },
            album: {
                min: parseInt(document.getElementById('albumYearMin').value),
                max: parseInt(document.getElementById('albumYearMax').value)
            },
            members: [...document.querySelectorAll('.members-grid input:checked')].map(i => parseInt(i.value)),
            locations: [...document.querySelectorAll('#locationsList input:checked')].map(i => i.value)
        };

        fetch('/filter', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(filters)
        })
        .then(res => res.json())
        .then(artists => {
            const grid = document.querySelector('.artists-grid');
            grid.innerHTML = artists.length ? artists.map(artist => `
                <div class="artist-card">
                    <a href="/artist?id=${artist.id}" class="artist-link">
                        <div class="image-container">
                            <img src="${artist.image}" alt="${artist.name}" loading="lazy">
                        </div>
                        <div class="artist-info">
                            <h2>${artist.name}</h2>
                        </div>
                    </a>
                </div>
            `).join('') : '<div class="no-results">No artists found matching your filters</div>';
        })
        .catch(() => {
            document.querySelector('.artists-grid').innerHTML = 
                '<div class="error-message">Failed to load filtered results</div>';
        });
    });
});