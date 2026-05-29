import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import ArticleSummaryModal from './ArticleSummaryModal';
import * as useSummaryHook from '../../hooks/useSummary';

// Mock useSummary custom hook
vi.mock('../../hooks/useSummary', () => ({
  useSummary: vi.fn(),
}));

describe('ArticleSummaryModal', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it('renders nothing when open is false', () => {
    vi.spyOn(useSummaryHook, 'useSummary').mockReturnValue({
      data: null,
      isLoading: false,
      error: null,
      refetch: vi.fn(),
    });

    const { container } = render(
      <ArticleSummaryModal articleId="1" open={false} onClose={vi.fn()} />
    );

    expect(container.firstChild).toBeNull();
  });

  it('renders skeleton loader when isLoading is true', () => {
    vi.spyOn(useSummaryHook, 'useSummary').mockReturnValue({
      data: null,
      isLoading: true,
      error: null,
      refetch: vi.fn(),
    });

    render(
      <ArticleSummaryModal articleId="1" open={true} onClose={vi.fn()} />
    );

    // Skeleton loader has animate-pulse class or role indicating skeleton
    expect(screen.getByLabelText('Kapat')).toBeInTheDocument();
    // Verify that loader elements are present (e.g. headers, paragraphs, etc.)
    expect(screen.getByText('HaberBot Yapay Zeka Özeti')).toBeInTheDocument();
  });

  it('renders error block when error occurs', () => {
    const mockRefetch = vi.fn();
    vi.spyOn(useSummaryHook, 'useSummary').mockReturnValue({
      data: null,
      isLoading: false,
      error: new Error('Network error'),
      refetch: mockRefetch,
    });

    render(
      <ArticleSummaryModal articleId="1" open={true} onClose={vi.fn()} />
    );

    expect(screen.getByText('Özet Yüklenemedi')).toBeInTheDocument();
    expect(screen.getByText('Tekrar Dene')).toBeInTheDocument();

    // Clicking retry triggers refetch
    fireEvent.click(screen.getByText('Tekrar Dene'));
    expect(mockRefetch).toHaveBeenCalledTimes(1);
  });

  it('renders markdown summary text successfully', () => {
    vi.spyOn(useSummaryHook, 'useSummary').mockReturnValue({
      data: { summary: 'Bu bir **kalın** özet ve *liste*:\n- Öğe 1\n- Öğe 2', article_id: '1' },
      isLoading: false,
      error: null,
      refetch: vi.fn(),
    });

    render(
      <ArticleSummaryModal articleId="1" open={true} onClose={vi.fn()} />
    );

    expect(screen.getByText('HaberBot Yapay Zeka Özeti')).toBeInTheDocument();
    // Since ReactMarkdown renders bold as strong tags
    expect(screen.getByText('kalın')).toBeInTheDocument();
    expect(screen.getByText('Öğe 1')).toBeInTheDocument();
    expect(screen.getByText('Öğe 2')).toBeInTheDocument();
  });

  it('triggers onClose when close button is clicked', () => {
    vi.spyOn(useSummaryHook, 'useSummary').mockReturnValue({
      data: { summary: 'Özet', article_id: '1' },
      isLoading: false,
      error: null,
      refetch: vi.fn(),
    });

    const mockOnClose = vi.fn();
    render(
      <ArticleSummaryModal articleId="1" open={true} onClose={mockOnClose} />
    );

    const closeBtn = screen.getByLabelText('Kapat');
    fireEvent.click(closeBtn);

    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });
});
