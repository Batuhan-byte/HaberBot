import './Skeleton.css';

/**
 * Animated skeleton placeholder that matches the ArticleCard layout.
 * Shown during loading states.
 */
function Skeleton() {
  return (
    <div className="skeleton-card">
      <div className="skeleton-card__title skeleton-pulse" />
      <div className="skeleton-card__title skeleton-card__title--short skeleton-pulse" />
      <div className="skeleton-card__body">
        <div className="skeleton-card__line skeleton-pulse" />
        <div className="skeleton-card__line skeleton-pulse" />
        <div className="skeleton-card__line skeleton-card__line--short skeleton-pulse" />
      </div>
      <div className="skeleton-card__footer">
        <div className="skeleton-card__badge skeleton-pulse" />
        <div className="skeleton-card__meta skeleton-pulse" />
      </div>
    </div>
  );
}

/**
 * Grid of skeleton cards for loading states.
 * @param {object} props
 * @param {number} props.count - Number of skeleton cards to render
 */
function SkeletonGrid({ count = 6 }) {
  return (
    <div className="article-grid">
      {Array.from({ length: count }, (_, index) => (
        <Skeleton key={index} />
      ))}
    </div>
  );
}

export { Skeleton, SkeletonGrid };
export default Skeleton;
