import axios from 'axios';

const apiClient = axios.create({
    baseURL: import.meta.env.VUE_APP_API_BASE_URL || '/api', // Use environment variables for the base URL
    timeout: 10000, // Timeout to prevent hanging requests
    headers: {
        Accept: 'application/json', // Ensure we accept JSON responses
    },
});

// Request Interceptor
apiClient.interceptors.request.use(
    (config) => {
        // Add Content-Type only for non-GET requests
        if (config.method !== 'get') {
            config.headers['Content-Type'] = 'application/json';
        } else {
            delete config.headers['Content-Type'];
        }

        // Attach authentication token if available
        const token = localStorage.getItem('authToken'); // Example: Use a secure storage mechanism
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }

        return config;
    },
    (error) => {
        console.error('[Request Error]', error);
        return Promise.reject(error);
    }
);

// Response Interceptor
apiClient.interceptors.response.use(
    (response) => {
        // Successful response handling
        return response;
    },
    (error) => {
        if (error.response) {
            console.error('[API Error]', {
                status: error.response.status,
                data: error.response.data,
            });

            // Handle specific HTTP status codes
            switch (error.response.status) {
                case 401:
                    alert('Session expired. Please log in again.');
                    // Optionally, clear token and redirect to login
                    localStorage.removeItem('authToken');
                    window.location.href = '/login'; // Replace with your login route
                    break;
                case 403:
                    alert('You do not have permission to perform this action.');
                    break;
                case 404:
                    alert('The requested resource was not found.');
                    break;
                case 500:
                    alert('An error occurred on the server. Please try again later.');
                    break;
                default:
                    alert(`Error: ${error.response.statusText}`);
            }
        } else if (error.request) {
            console.error('[No Response]', error.request);
            alert('Unable to connect to the server. Please check your internet connection.');
        } else {
            console.error('[Unexpected Error]', error.message);
            alert('An unexpected error occurred. Please try again.');
        }

        return Promise.reject(error); // Ensure errors can be caught in the component
    }
);

// Example API methods
const api = {
    createExpense(data) {
        return apiClient.post('/expenses', data);
    },
    getExpenses(page, size) {
        return apiClient.get(`/expenses?page=${page}&size=${size}`);
    },
    getCategories() {
        return apiClient.get('/categories');
    },
    addCategory(data) {
        return apiClient.post('/categories', data);
    },
    updateCategory(id, data) {
        return apiClient.put(`/categories/${id}`, data);
    },
    deleteCategory(id) {
        return apiClient.delete(`/categories/${id}`);
    },
};

export default api;