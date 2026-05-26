import { SOURCE_LABELS, SOURCE_TYPES } from '../../utils/constants';
import './SourceBadge.css';

/**
 * Small pill badge showing the article source with a colored indicator dot.
 * @param {object} props
 * @param {string} props.source - Source type key ("hackernews" | "rss")
 */
function SourceBadge({ source }) {
  const label = SOURCE_LABELS[source] || source;
  const sourceClass = source === SOURCE_TYPES.HACKERNEWS
    ? 'source-badge--hn'
    : 'source-badge--rss';

  return (
    <span className={`source-badge ${sourceClass}`}>
      <span className="source-badge__dot" />
      <span className="source-badge__label">{label}</span>
    </span>
  );
}

export default SourceBadge;
