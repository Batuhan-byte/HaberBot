import { render, screen } from '@testing-library/react';
import { Skeleton, SkeletonGrid } from './Skeleton';
import { describe, it, expect } from 'vitest';

describe('Skeleton', () => {
  it('renders a single skeleton card', () => {
    const { container } = render(<Skeleton />);

    expect(container.querySelector('.skeleton-card')).toBeInTheDocument();
  });

  it('renders pulse elements', () => {
    const { container } = render(<Skeleton />);

    const pulses = container.querySelectorAll('.skeleton-pulse');
    expect(pulses.length).toBeGreaterThan(0);
  });
});

describe('SkeletonGrid', () => {
  it('renders default count of 6 skeletons', () => {
    const { container } = render(<SkeletonGrid />);

    const cards = container.querySelectorAll('.skeleton-card');
    expect(cards).toHaveLength(6);
  });

  it('renders custom count of skeletons', () => {
    const { container } = render(<SkeletonGrid count={3} />);

    const cards = container.querySelectorAll('.skeleton-card');
    expect(cards).toHaveLength(3);
  });

  it('renders in article-grid wrapper', () => {
    const { container } = render(<SkeletonGrid count={4} />);

    expect(container.querySelector('.article-grid')).toBeInTheDocument();
  });
});
