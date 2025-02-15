document.addEventListener('DOMContentLoaded', () => {
    // Clear filters functionality
    const clearButton = document.querySelector('.clear-filters');
    const filterForm = document.getElementById('filter-form');

    if (clearButton && filterForm) {
        clearButton.addEventListener('click', (e) => {
            e.preventDefault();
            
            // Clear all checkboxes
            filterForm.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
                checkbox.checked = false;
            });

            // Reset range sliders to defaults
            document.querySelector('.creation-start').value = '1950';
            document.querySelector('.creation-end').value = '2024';
            document.querySelector('.album-start').value = '1950';
            document.querySelector('.album-end').value = '2024';
            
            // Update range value displays
            updateRangeValues();

            // Clear search input if exists
            const searchInput = filterForm.querySelector('.location-search-input');
            if (searchInput) {
                searchInput.value = '';
            }

            filterForm.submit();
        });
    }

    // Setup range slider functionality
    function setupRangeSlider(startClass, endClass, startValueId, endValueId) {
        const startSlider = document.querySelector(`.${startClass}`);
        const endSlider = document.querySelector(`.${endClass}`);
        const startValue = document.getElementById(startValueId);
        const endValue = document.getElementById(endValueId);

        if (startSlider && endSlider) {
            startSlider.addEventListener('input', () => {
                const start = parseInt(startSlider.value);
                const end = parseInt(endSlider.value);
                if (start > end) {
                    startSlider.value = end;
                }
                startValue.textContent = startSlider.value;
            });

            endSlider.addEventListener('input', () => {
                const start = parseInt(startSlider.value);
                const end = parseInt(endSlider.value);
                if (end < start) {
                    endSlider.value = start;
                }
                endValue.textContent = endSlider.value;
            });
        }
    }

    // Setup both range sliders
    setupRangeSlider(
        'creation-start', 
        'creation-end', 
        'creation-start-value', 
        'creation-end-value'
    );
    
    setupRangeSlider(
        'album-start', 
        'album-end', 
        'album-start-value', 
        'album-end-value'
    );

    // Update all range values displays
    function updateRangeValues() {
        const displays = {
            'creation-start-value': '.creation-start',
            'creation-end-value': '.creation-end',
            'album-start-value': '.album-start',
            'album-end-value': '.album-end'
        };

        Object.entries(displays).forEach(([displayId, sliderClass]) => {
            const display = document.getElementById(displayId);
            const slider = document.querySelector(sliderClass);
            if (display && slider) {
                display.textContent = slider.value;
            }
        });
    }

    // Initialize range values
    updateRangeValues();
});