document.addEventListener('DOMContentLoaded', function() {
    // Toggle the visibility of the filters container when the Filters button is clicked
    const filtersBtn = document.getElementById('filters-btn');
    const filtersContainer = document.getElementById('filters-container');
  
    if (filtersBtn && filtersContainer) {
      filtersBtn.addEventListener('click', function() {
        if (filtersContainer.style.display === 'none' || filtersContainer.style.display === '') {
          filtersContainer.style.display = 'block';
        } else {
          filtersContainer.style.display = 'none';
        }
      });
    }
  
    // Function to fetch filtered results using AJAX
    function fetchFilteredResults() {
      // Gather the filter values from the form inputs
      const creationMin = document.getElementById('creation_min').value;
      const creationMax = document.getElementById('creation_max').value;
      const albumMin = document.getElementById('album_min').value;
      const albumMax = document.getElementById('album_max').value;
      const membersMin = document.getElementById('members_min').value;
      const membersMax = document.getElementById('members_max').value;
  
      // For checkboxes (concert locations), collect all checked values.
      const locationCheckboxes = document.querySelectorAll('input[name="locations"]:checked');
      let locations = [];
      locationCheckboxes.forEach(function(checkbox) {
        locations.push(checkbox.value);
      });
  
      // Build the query string with the provided filter values
      let params = new URLSearchParams();
      if (creationMin) params.append('creation_min', creationMin);
      if (creationMax) params.append('creation_max', creationMax);
      if (albumMin) params.append('album_min', albumMin);
      if (albumMax) params.append('album_max', albumMax);
      if (membersMin) params.append('members_min', membersMin);
      if (membersMax) params.append('members_max', membersMax);
      locations.forEach(function(loc) {
        params.append('locations', loc);
      });
  
      // Send AJAX (fetch) request to the filters endpoint
      fetch('/filters?' + params.toString(), {
        headers: {
          'X-Requested-With': 'XMLHttpRequest'
        }
      })
        .then(function(response) {
          return response.json();
        })
        .then(function(data) {
          // Update the results container with the filtered results
          const resultsContainer = document.getElementById('filters-results');
          resultsContainer.innerHTML = '';
  
          if (data.length === 0) {
            resultsContainer.innerHTML = '<p>No results found.</p>';
          } else {
            data.forEach(function(result) {
              // Each result is formatted as a clickable artist card
              const div = document.createElement('div');
              div.classList.add('filter-result-item');
              div.innerHTML = `
                <a href="/artist?id=${result.ArtistId}">
                  <h3>${result.ArtistName}</h3>
                  <p>${result.Description}</p>
                </a>
              `;
              resultsContainer.appendChild(div);
            });
          }
        })
        .catch(function(error) {
          console.error('Error fetching filter results:', error);
        });
    }
  
    // Attach change event listeners to all filter inputs within the filters container.
    // This will trigger an AJAX request every time a filter value changes.
    const filterInputs = document.querySelectorAll('#filters-container input');
    filterInputs.forEach(function(input) {
      input.addEventListener('change', fetchFilteredResults);
    });
  });
  