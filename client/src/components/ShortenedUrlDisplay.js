import React from 'react';

const ShortenedUrlDisplay = ({ shortUrl }) => {
  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(shortUrl);
      // Show temporary success feedback
      const originalText = document.getElementById('copy-button').textContent;
      document.getElementById('copy-button').textContent = 'Copied!';
      setTimeout(() => {
        document.getElementById('copy-button').textContent = originalText;
      }, 2000);
    } catch (err) {
      alert('Failed to copy to clipboard');
    }
  };

  return (
    <div className="shortened-url-display">
      <div className="short-url">{shortUrl}</div>
      <button
        id="copy-button"
        onClick={copyToClipboard}
        className="copy-button"
      >
        Copy
      </button>
    </div>
  );
};

export default ShortenedUrlDisplay;