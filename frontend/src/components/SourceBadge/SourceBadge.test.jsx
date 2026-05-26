import { render, screen } from '@testing-library/react';
import SourceBadge from './SourceBadge';
import { describe, it, expect } from 'vitest';

describe('SourceBadge', () => {
  it('renders HackerNews label for hackernews source', () => {
    render(<SourceBadge source="hackernews" />);

    expect(screen.getByText('HackerNews')).toBeInTheDocument();
  });

  it('renders RSS label for rss source', () => {
    render(<SourceBadge source="rss" />);

    expect(screen.getByText('RSS')).toBeInTheDocument();
  });

  it('renders source value as label for unknown source', () => {
    render(<SourceBadge source="custom-source" />);

    expect(screen.getByText('custom-source')).toBeInTheDocument();
  });

  it('applies HN class for hackernews', () => {
    const { container } = render(<SourceBadge source="hackernews" />);

    expect(container.querySelector('.source-badge--hn')).toBeInTheDocument();
  });

  it('applies RSS class for rss', () => {
    const { container } = render(<SourceBadge source="rss" />);

    expect(container.querySelector('.source-badge--rss')).toBeInTheDocument();
  });
});
