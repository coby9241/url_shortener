// API service for communicating with the URL shortener backend
// This assumes the backend is running on the same origin or properly configured CORS

const API_BASE_URL = process.env.REACT_APP_API_URL || '';

export const shortenURL = async (longUrl) => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/shorten`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ 'url': longUrl }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to shorten URL');
    }

    const data = await response.json();
    return data;
  } catch (error) {
    throw error;
  }
};