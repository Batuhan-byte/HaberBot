/**
 * Formats a date string into a Turkish relative time expression.
 * Examples: "2 saat önce", "3 gün önce", "az önce"
 * @param {string|Date} dateInput - ISO date string or Date object
 * @returns {string} Turkish relative time string
 */
export function formatRelativeDate(dateInput) {
  const date = dateInput instanceof Date ? dateInput : new Date(dateInput);
  const now = new Date();
  const diffMs = now - date;
  const diffSeconds = Math.floor(diffMs / 1000);
  const diffMinutes = Math.floor(diffSeconds / 60);
  const diffHours = Math.floor(diffMinutes / 60);
  const diffDays = Math.floor(diffHours / 24);
  const diffWeeks = Math.floor(diffDays / 7);
  const diffMonths = Math.floor(diffDays / 30);

  if (diffSeconds < 60) {
    return 'az önce';
  }

  if (diffMinutes < 60) {
    return `${diffMinutes} dakika önce`;
  }

  if (diffHours < 24) {
    return `${diffHours} saat önce`;
  }

  if (diffDays < 7) {
    return `${diffDays} gün önce`;
  }

  if (diffWeeks < 4) {
    return `${diffWeeks} hafta önce`;
  }

  if (diffMonths < 12) {
    return `${diffMonths} ay önce`;
  }

  const diffYears = Math.floor(diffMonths / 12);
  return `${diffYears} yıl önce`;
}

/**
 * Formats a date string into a Turkish formatted date.
 * Example: "25 Ocak 2026"
 * @param {string|Date} dateInput - ISO date string or Date object
 * @returns {string} Turkish formatted date
 */
export function formatTurkishDate(dateInput) {
  const date = dateInput instanceof Date ? dateInput : new Date(dateInput);
  const months = [
    'Ocak', 'Şubat', 'Mart', 'Nisan', 'Mayıs', 'Haziran',
    'Temmuz', 'Ağustos', 'Eylül', 'Ekim', 'Kasım', 'Aralık',
  ];

  const day = date.getDate();
  const month = months[date.getMonth()];
  const year = date.getFullYear();

  return `${day} ${month} ${year}`;
}
