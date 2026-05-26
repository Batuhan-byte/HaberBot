import { useState, useEffect } from 'react';

/**
 * useDebounce hook
 * 
 * Amaç (Purpose): Arama kutusuna her harf girildiğinde anında API isteği atmamak için kullanılır.
 * Gecikme süresi (Delay): 300ms idealdir. Kullanıcının yazmayı bitirmesi için makul bir süredir, 
 * çok kısa olursa fazla istek gider, çok uzun olursa arayüz yavaş hissettirir.
 * 
 * @param {any} value - Debounce edilecek değer (örn: arama metni)
 * @param {number} delay - Milisaniye cinsinden bekleme süresi
 * @returns {any} - Debounce edilmiş değer
 */
export function useDebounce(value, delay = 300) {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    // Değer değiştiğinde bir zamanlayıcı başlat
    const timer = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    // Eğer delay süresi dolmadan value tekrar değişirse (kullanıcı yazmaya devam ederse),
    // önceki zamanlayıcıyı iptal et
    return () => {
      clearTimeout(timer);
    };
  }, [value, delay]);

  return debouncedValue;
}
