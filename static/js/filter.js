document.addEventListener('DOMContentLoaded', () => {
    // Clear filters functionality - simply redirects to /filter
    const clearButton = document.querySelector('.clear-filters');
    if (clearButton) {
        clearButton.addEventListener('click', (e) => {
            e.preventDefault();
            window.location.href = '/filter';
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
});