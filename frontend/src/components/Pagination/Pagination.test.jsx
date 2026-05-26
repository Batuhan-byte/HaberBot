import { render, screen, fireEvent } from '@testing-library/react';
import Pagination from './Pagination';
import { describe, it, expect, vi } from 'vitest';

describe('Pagination', () => {
  it('renders page buttons', () => {
    render(<Pagination currentPage={1} totalPages={5} onPageChange={() => {}} />);

    expect(screen.getByLabelText('Sayfa 1')).toBeInTheDocument();
    expect(screen.getByLabelText('Sayfa 5')).toBeInTheDocument();
  });

  it('returns null when totalPages is 1 or less', () => {
    const { container } = render(
      <Pagination currentPage={1} totalPages={1} onPageChange={() => {}} />
    );

    expect(container.innerHTML).toBe('');
  });

  it('disables previous button on first page', () => {
    render(<Pagination currentPage={1} totalPages={5} onPageChange={() => {}} />);

    expect(screen.getByLabelText('Önceki sayfa')).toBeDisabled();
  });

  it('disables next button on last page', () => {
    render(<Pagination currentPage={5} totalPages={5} onPageChange={() => {}} />);

    expect(screen.getByLabelText('Sonraki sayfa')).toBeDisabled();
  });

  it('calls onPageChange when clicking a page button', () => {
    const onPageChange = vi.fn();

    render(<Pagination currentPage={3} totalPages={5} onPageChange={onPageChange} />);

    fireEvent.click(screen.getByLabelText('Sayfa 5'));

    expect(onPageChange).toHaveBeenCalledWith(5);
  });

  it('calls onPageChange when clicking next', () => {
    const onPageChange = vi.fn();

    render(<Pagination currentPage={2} totalPages={5} onPageChange={onPageChange} />);

    fireEvent.click(screen.getByLabelText('Sonraki sayfa'));

    expect(onPageChange).toHaveBeenCalledWith(3);
  });

  it('calls onPageChange when clicking previous', () => {
    const onPageChange = vi.fn();

    render(<Pagination currentPage={3} totalPages={5} onPageChange={onPageChange} />);

    fireEvent.click(screen.getByLabelText('Önceki sayfa'));

    expect(onPageChange).toHaveBeenCalledWith(2);
  });

  it('highlights current page button', () => {
    render(<Pagination currentPage={3} totalPages={5} onPageChange={() => {}} />);

    const activeBtn = screen.getByLabelText('Sayfa 3');
    expect(activeBtn).toHaveAttribute('aria-current', 'page');
  });

  it('shows ellipsis for large page counts', () => {
    render(<Pagination currentPage={5} totalPages={20} onPageChange={() => {}} />);

    const ellipsis = document.querySelectorAll('.pagination__ellipsis');
    expect(ellipsis.length).toBeGreaterThanOrEqual(1);
  });
});
