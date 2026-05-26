import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ArticleCard from './ArticleCard';
import { describe, it, expect } from 'vitest';

const baseArticle = {
  id: '1',
  title: 'Original Title',
  title_tr: 'Türkçe Başlık',
  summary: 'Original summary',
  summary_tr: 'Türkçe özet',
  source: 'hackernews',
  original_url: 'https://example.com',
  created_at: new Date().toISOString(),
  image_url: 'https://example.com/img.jpg',
};

function renderWithRouter(ui) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe('ArticleCard', () => {
  it('renders translated title and summary', () => {
    renderWithRouter(<ArticleCard article={baseArticle} index={0} />);

    expect(screen.getByText('Türkçe Başlık')).toBeInTheDocument();
    expect(screen.getByText('Türkçe özet')).toBeInTheDocument();
  });

  it('renders source badge and date', () => {
    renderWithRouter(<ArticleCard article={baseArticle} index={0} />);

    expect(screen.getByText('hackernews')).toBeInTheDocument();
  });

  it('links to article detail page', () => {
    renderWithRouter(<ArticleCard article={baseArticle} index={0} />);

    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/haber/1');
  });

  it('renders HN tags for hackernews source', () => {
    renderWithRouter(<ArticleCard article={baseArticle} index={0} />);

    expect(screen.getByText('#Tech')).toBeInTheDocument();
    expect(screen.getByText('#HN')).toBeInTheDocument();
  });

  it('renders Reddit tags for reddit source', () => {
    renderWithRouter(
      <ArticleCard article={{ ...baseArticle, source: 'reddit' }} index={0} />
    );

    expect(screen.getByText('#Reddit')).toBeInTheDocument();
    expect(screen.getByText('#Dev')).toBeInTheDocument();
  });

  it('renders default tags for unknown source', () => {
    renderWithRouter(
      <ArticleCard article={{ ...baseArticle, source: 'rss' }} index={0} />
    );

    expect(screen.getByText('#Haber')).toBeInTheDocument();
    expect(screen.getByText('#Gündem')).toBeInTheDocument();
  });

  it('falls back to original title when no translation', () => {
    renderWithRouter(
      <ArticleCard
        article={{ ...baseArticle, title_tr: undefined }}
        index={0}
      />
    );

    expect(screen.getByText('Original Title')).toBeInTheDocument();
  });

  it('renders wide card variant for index 3', () => {
    const { container } = renderWithRouter(
      <ArticleCard article={baseArticle} index={3} />
    );

    expect(container.querySelector('.article-card--wide')).toBeInTheDocument();
    expect(container.querySelector('.article-card__image-wrap')).toBeInTheDocument();
  });

  it('renders standard card for non-3 index', () => {
    const { container } = renderWithRouter(
      <ArticleCard article={baseArticle} index={0} />
    );

    expect(container.querySelector('.article-card--wide')).not.toBeInTheDocument();
  });
});
