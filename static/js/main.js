// Main JavaScript file for the WhatsApp Gateway UI

// Store API key in localStorage after the user enters it once
function storeApiKey(apiKey) {
    if (apiKey) {
        localStorage.setItem('api_key', apiKey);
    }
}

// Get API key from localStorage or prompt the user
function getApiKey() {
    let apiKey = localStorage.getItem('api_key');
    
    if (!apiKey) {
        // Use default API key from .env file first
        apiKey = 'changeme'; // default from .env
        storeApiKey(apiKey);
        console.log('Using default API key: changeme');
    }
    
    return apiKey;
}

// Function to reset the API key (for troubleshooting)
function resetApiKey() {
    localStorage.removeItem('api_key');
    const newKey = prompt('Please enter your API key:', 'changeme');
    storeApiKey(newKey);
    return newKey;
}

// Helper function to format date/time
function formatDateTime(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString();
}

// Function to handle common AJAX errors
function handleAjaxError(xhr, element) {
    let errorMessage = 'An error occurred';
    
    if (xhr.responseJSON && xhr.responseJSON.error) {
        errorMessage = xhr.responseJSON.error;
    } else if (xhr.status === 401) {
        errorMessage = 'Unauthorized: Invalid API key';
        // Clear stored API key if it's invalid
        localStorage.removeItem('api_key');
    } else if (xhr.status === 404) {
        errorMessage = 'Resource not found';
    } else if (xhr.status === 500) {
        errorMessage = 'Server error';
    }
    
    if (element) {
        $(element).html(`<div class="alert alert-danger">${errorMessage}</div>`);
    } else {
        alert(errorMessage);
    }
}

// Initialize tooltips and popovers on document ready
$(document).ready(function() {
    console.log('WhatsApp Gateway UI initialized');
    
    // Add visual indicator that JavaScript is working
    $('body').append('<div id="js-debug" style="position: fixed; bottom: 10px; right: 10px; background: rgba(0,0,0,0.7); color: white; padding: 5px 10px; border-radius: 5px; z-index: 9999;">JS Running</div>');
    
    // Check if CSS is loaded
    if ($('.navbar').css('margin-bottom') === '20px') {
        $('#js-debug').append(' | CSS OK');
    } else {
        $('#js-debug').append(' | CSS Missing').css('background-color', 'rgba(255,0,0,0.7)');
        console.error('CSS styles not applied correctly');
    }
    
    // Initialize Bootstrap tooltips
    try {
        const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
        tooltipTriggerList.map(function(tooltipTriggerEl) {
            return new bootstrap.Tooltip(tooltipTriggerEl);
        });
        
        // Initialize Bootstrap popovers
        const popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'));
        popoverTriggerList.map(function(popoverTriggerEl) {
            return new bootstrap.Popover(popoverTriggerEl);
        });
    } catch (error) {
        console.error('Error initializing Bootstrap components:', error);
        $('#js-debug').append(' | Bootstrap Error').css('background-color', 'rgba(255,0,0,0.7)');
    }

    // Add test API request to verify API connectivity
    $.ajax({
        url: '/api/clients',
        method: 'GET',
        headers: {
            'X-API-Key': getApiKey()
        },
        success: function(response) {
            $('#js-debug').append(' | API OK');
            console.log('API test successful:', response);
        },
        error: function(xhr) {
            $('#js-debug').append(' | API Error').css('background-color', 'rgba(255,0,0,0.7)');
            console.error('API test failed:', xhr.status, xhr.responseText);
        }
    });
});
