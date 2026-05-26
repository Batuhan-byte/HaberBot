import './Pagination.css';

/**
 * Modern pill-style pagination with page numbers.
 * @param {object} props
 * @param {number} props.currentPage - Active page number
 * @param {number} props.totalPages - Total number of pages
 * @param {function} props.onPageChange - Callback when page is selected
 */
function Pagination({ currentPage, totalPages, onPageChange }) {
  if (totalPages <= 1) return null;

  const pages = buildPageNumbers(currentPage, totalPages);

  return (
    <nav className="pagination" aria-label="Sayfa navigasyonu">
      <button
        className="pagination__btn pagination__btn--nav"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage <= 1}
        aria-label="Önceki sayfa"
      >
        <ChevronLeftIcon />
        <span className="pagination__btn-text">Önceki</span>
      </button>

      <div className="pagination__pages">
        {pages.map((page, index) =>
          page === '...' ? (
            <span key={`ellipsis-${index}`} className="pagination__ellipsis">
              ...
            </span>
          ) : (
            <button
              key={page}
              className={`pagination__btn pagination__btn--page ${
                page === currentPage ? 'pagination__btn--active' : ''
              }`}
              onClick={() => onPageChange(page)}
              aria-label={`Sayfa ${page}`}
              aria-current={page === currentPage ? 'page' : undefined}
            >
              {page}
            </button>
          )
        )}
      </div>

      <button
        className="pagination__btn pagination__btn--nav"
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage >= totalPages}
        aria-label="Sonraki sayfa"
      >
        <span className="pagination__btn-text">Sonraki</span>
        <ChevronRightIcon />
      </button>
    </nav>
  );
}

/**
 * Builds an array of page numbers with ellipsis for large page counts.
 * @param {number} current - Current page
 * @param {number} total - Total pages
 * @returns {Array} Array of page numbers and '...' strings
 */
function buildPageNumbers(current, total) {
  const delta = 1;
  const pages = [];

  for (let i = 1; i <= total; i++) {
    if (
      i === 1 ||
      i === total ||
      (i >= current - delta && i <= current + delta)
    ) {
      pages.push(i);
    } else if (pages[pages.length - 1] !== '...') {
      pages.push('...');
    }
  }

  return pages;
}

function ChevronLeftIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
      <path
        d="M10 12L6 8L10 4"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

function ChevronRightIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
      <path
        d="M6 12L10 8L6 4"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

export default Pagination;
