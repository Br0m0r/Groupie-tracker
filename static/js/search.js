// Wait for the DOM to be fully loaded before executing the script.
document.addEventListener('DOMContentLoaded', () => {
    // Select the search input element.
    const searchInput = document.querySelector('.search-input');
    // Select the container that will hold the list of suggestions.
    const suggestionsList = document.querySelector('.suggestions-list');
    // Select the overall search suggestions container (used to show/hide suggestions).
    const searchContainer = document.querySelector('.search-suggestions');

    // Add an event listener to the search input for 'input' events (fired on each keystroke).
    searchInput.addEventListener('input', (e) => {
        // Retrieve the current value of the search input and trim whitespace.
        const query = e.target.value.trim();
        
        // If the query is empty, clear the suggestions and hide the suggestions container.
        if (query === '') {
            suggestionsList.innerHTML = '';
            searchContainer.style.display = 'none';
            return;
        }

        // Send an AJAX request (using fetch) to the /search endpoint with the encoded query.
        // The header 'X-Requested-With': 'XMLHttpRequest' marks this as an AJAX request.
        fetch(`/search?q=${encodeURIComponent(query)}`, {
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        })
        // Parse the response as JSON.
        .then(response => response.json())
        // Process the JSON results.
        .then(results => {
            // If no results are returned, clear suggestions and hide the container.
            if (results.length === 0) {
                suggestionsList.innerHTML = '';
                searchContainer.style.display = 'none';
                return;
            }

            // Otherwise, display the suggestions container.
            searchContainer.style.display = 'block';
            // Render all the suggestions by mapping over the results array.
            // Each suggestion is inserted into the container as HTML.
            suggestionsList.innerHTML = results.map(suggestion => `
                <div class="suggestion-item">
                    <span class="suggestion-text">${suggestion.text}</span>
                    <span class="suggestion-type">${suggestion.type}</span>
                </div>
            `).join('');  // join('') converts the array of HTML strings into a single string.
        })
        // If an error occurs during the fetch or processing, log the error and hide the suggestions.
        .catch(error => {
            console.error('Failed to fetch suggestions:', error);
            searchContainer.style.display = 'none';
        });
    });

    // Add a click event listener to the document to close the suggestions panel
    // if the user clicks outside of the search input or suggestions container.
    document.addEventListener('click', (e) => {
        // If the clicked element is not within the suggestions container or the search input,
        // hide the suggestions container.
        if (!searchContainer.contains(e.target) && !searchInput.contains(e.target)) {
            searchContainer.style.display = 'none';
        }
    });
});
