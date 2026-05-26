import './ErrorState.css';

/**
 * Error state component with retry button.
 * @param {object} props
 * @param {string} props.message - Error message to display
 * @param {function} props.onRetry - Callback for retry button
 */
function ErrorState({ message, onRetry }) {
  return (
    <div className="error-state">
      <div className="error-state__icon">
        <span className="error-state__emoji">⚠️</span>
        <div className="error-state__glow" />
      </div>
      <h3 className="error-state__title">Bir şeyler ters gitti</h3>
      <p className="error-state__message">
        {message || 'Veriler yüklenirken bir hata oluştu. Lütfen tekrar deneyin.'}
      </p>
      {onRetry && (
        <button className="error-state__retry" onClick={onRetry}>
          <RefreshIcon />
          Tekrar Dene
        </button>
      )}
    </div>
  );
}

function RefreshIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
      <path
        d="M13.65 2.35A8 8 0 1 0 16 8h-2a6 6 0 1 1-1.76-4.24L10 6h6V0l-2.35 2.35z"
        fill="currentColor"
        transform="scale(0.85) translate(1, 1)"
      />
    </svg>
  );
}

export default ErrorState;
