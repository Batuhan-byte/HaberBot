import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ArticleSummaryLauncher from './ArticleSummaryLauncher';

describe('ArticleSummaryLauncher', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.restoreAllMocks();
  });

  it('renders correctly with sparkle icon and text', () => {
    render(<ArticleSummaryLauncher onClick={vi.fn()} />);

    expect(screen.getByText('Yapay Zeka Özeti')).toBeInTheDocument();
    expect(screen.getByText('auto_awesome')).toBeInTheDocument();
  });

  it('triggers onClick callback when clicked regardless of auth state', () => {
    const mockOnClick = vi.fn();

    render(<ArticleSummaryLauncher onClick={mockOnClick} />);

    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(mockOnClick).toHaveBeenCalledTimes(1);
  });
});
