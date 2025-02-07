// static/js/filter.js

//clear filters button  !!!
document.addEventListener('DOMContentLoaded', () => {
    // Get the clear button and filter form
    const clearButton = document.querySelector('.clear-filters');
    const filterForm = document.getElementById('filter-form');

    if (clearButton && filterForm) {
        clearButton.addEventListener('click', (e) => {
            e.preventDefault(); // Prevent the default reset behavior

            // Clear all checkboxes
            filterForm.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
                checkbox.checked = false;
            });

            // Clear all number inputs and set to default values
            filterForm.querySelectorAll('input[type="number"]').forEach(input => {
                if (input.name.includes('creation')) {
                    input.value = input.name.includes('start') ? '1950' : '2024';
                } else if (input.name.includes('album')) {
                    input.value = input.name.includes('start') ? '1950' : '2024';
                }
            });

            // Clear search input if exists
            const searchInput = filterForm.querySelector('.location-search-input');
            if (searchInput) {
                searchInput.value = '';
            }

            // Submit the form to refresh with cleared filters
            filterForm.submit();
        });
    }
});
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

    // Range slider functionality
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
    setupRangeSlider('creation-start', 'creation-end', 'creation-start-value', 'creation-end-value');
    setupRangeSlider('album-start', 'album-end', 'album-start-value', 'album-end-value');

    // Initial update of range values
    function updateRangeValues() {
        document.getElementById('creation-start-value').textContent = document.querySelector('.creation-start').value;
        document.getElementById('creation-end-value').textContent = document.querySelector('.creation-end').value;
        document.getElementById('album-start-value').textContent = document.querySelector('.album-start').value;
        document.getElementById('album-end-value').textContent = document.querySelector('.album-end').value;
    }

    updateRangeValues();
});