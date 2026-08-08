import React, { useState } from 'react';
import UrlForm from './components/UrlForm';
import ShortenedUrlDisplay from './components/ShortenedUrlDisplay';
import { shortenURL } from './services/api';
import './styles.css';

function App() {
  const [longUrl, setLongUrl] = useState('');
  const [shortUrl, setShortUrl] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Basic client-side validation
    if (!longUrl.trim()) {
      setError('Please enter a URL');
      return;
    }

    try {
      setLoading(true);
      setError('');
      const response = await shortenURL(longUrl.trim());
      setShortUrl(response.shortUrl);
      // Redirect to the short URL (will be handled by server redirect)
      window.location.href = response.shortUrl;
    } catch (err) {
      setError(err.message || 'Failed to shorten URL');
      setShortUrl(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>URL Shortener</h1>
        <p>Enter a long URL to get a shortened version</p>
      </header>

      <main>
        <UrlForm
          value={longUrl}
          onChange={(e) => setLongUrl(e.target.value)}
          onSubmit={handleSubmit}
          loading={loading}
          error={error}
        />

        {shortUrl && (
          <ShortenedUrlDisplay
            shortUrl={shortUrl}
          />
        )}
      </main>
    </div>
  );
}

export default App;