document.addEventListener('DOMContentLoaded', () => {
    // Range input display updates
    const updateRangeDisplay = (input, displayId) => {
        const display = document.getElementById(displayId);
        if (display) {
            display.textContent = input.value;
        }
    };

    // Creation year inputs
    document.querySelectorAll('input[name="creation_min"], input[name="creation_max"]')
        .forEach(input => {
            input.addEventListener('input', (e) => {
                const min = document.querySelector('input[name="creation_min"]');
                const max = document.querySelector('input[name="creation_max"]');
                
                if (parseInt(min.value) > parseInt(max.value)) {
                    if (e.target === min) {
                        max.value = min.value;
                    } else {
                        min.value = max.value;
                    }
                }
            });
        });

    // Album year inputs
    document.querySelectorAll('input[name="album_min"], input[name="album_max"]')
        .forEach(input => {
            input.addEventListener('input', (e) => {
                const min = document.querySelector('input[name="album_min"]');
                const max = document.querySelector('input[name="album_max"]');
                
                if (parseInt(min.value) > parseInt(max.value)) {
                    if (e.target === min) {
                        max.value = min.value;
                    } else {
                        min.value = max.value;
                    }
                }
            });
        });

    // Toggle locations list
    const locationsHeader = document.querySelector('.locations-header');
    const locationsList = document.querySelector('.locations-list');
    
    if (locationsHeader && locationsList) {
        locationsHeader.addEventListener('click', () => {
            const isHidden = locationsList.style.display === 'none';
            locationsList.style.display = isHidden ? 'block' : 'none';
            locationsHeader.innerHTML = `Locations ${isHidden ? '▲' : '▼'}`;
        });
    }
});