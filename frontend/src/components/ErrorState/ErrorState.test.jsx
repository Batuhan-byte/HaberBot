import { render, screen, fireEvent } from '@testing-library/react';
import ErrorState from './ErrorState';
import { describe, it, expect, vi } from 'vitest';

describe('ErrorState', () => {
  it('renders default error message when none provided', () => {
    render(<ErrorState />);

    expect(screen.getByText('Bir şeyler ters gitti')).toBeInTheDocument();
    expect(screen.getByText('Veriler yüklenirken bir hata oluştu. Lütfen tekrar deneyin.')).toBeInTheDocument();
  });

  it('renders custom error message', () => {
    render(<ErrorState message="Özel hata mesajı" />);

    expect(screen.getByText('Özel hata mesajı')).toBeInTheDocument();
  });

  it('does not render retry button when onRetry is not provided', () => {
    render(<ErrorState message="Hata" />);

    expect(screen.queryByText('Tekrar Dene')).not.toBeInTheDocument();
  });

  it('renders retry button when onRetry is provided', () => {
    render(<ErrorState message="Hata" onRetry={() => {}} />);

    expect(screen.getByText('Tekrar Dene')).toBeInTheDocument();
  });

  it('calls onRetry when retry button is clicked', () => {
    const onRetry = vi.fn();

    render(<ErrorState message="Hata" onRetry={onRetry} />);

    fireEvent.click(screen.getByText('Tekrar Dene'));
    expect(onRetry).toHaveBeenCalledTimes(1);
  });
});
