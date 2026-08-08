import React from 'react';

const UrlForm = ({ value, onChange, onSubmit, loading, error }) => {
  return (
    <form onSubmit={onSubmit} className="url-form">
      <div className="form-group">
        <label htmlFor="longUrl">Enter your long URL:</label>
        <input
          type="url"
          id="longUrl"
          value={value}
          onChange={onChange}
          placeholder="https://example.com/very/long/url/that/needs/to/be/shortened"
          required
          disabled={loading}
          className={loading ? 'loading' : ''}
        />
      </div>

      {error && <div className="error-message">{error}</div>}

      <button
        type="submit"
        disabled={loading || !value.trim()}
        className={`submit-button ${loading ? 'loading' : ''}`}
      >
        {loading ? 'Shortening...' : 'Shorten URL'}
      </button>
    </form>
  );
};

export default UrlForm;