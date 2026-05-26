import React, { useState, useEffect } from 'react';
import './ConsoleWidget.css';

const ConsoleWidget = () => {
  const [isOpen, setIsOpen] = useState(false);
  const [apiKey, setApiKey] = useState('');
  const [statusMsg, setStatusMsg] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  // Load API Key from localStorage if exists
  useEffect(() => {
    const savedKey = localStorage.getItem('haberbot_admin_key');
    if (savedKey) {
      setApiKey(savedKey);
    }
  }, []);

  const handleKeyChange = (e) => {
    const val = e.target.value;
    setApiKey(val);
    localStorage.setItem('haberbot_admin_key', val);
  };

  const triggerAction = async (actionType) => {
    if (!apiKey) {
      setStatusMsg('HATA: Lütfen önce Admin API Anahtarını girin.');
      return;
    }

    setIsLoading(true);
    setStatusMsg(`İŞLEM BAŞLADI: Pipeline ${actionType.toUpperCase()} tetikleniyor...`);

    try {
      const endpoint = actionType === 'fetch' ? '/api/v1/admin/fetch' : '/api/v1/admin/process';
      const apiUrl = `${import.meta.env.VITE_API_URL}${endpoint}`;
      
      const response = await fetch(apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Admin-API-Key': apiKey,
        },
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'İstek başarısız oldu.');
      }

      if (actionType === 'fetch') {
        setStatusMsg(`BAŞARILI: ${data.fetched} yeni haber kaynaklardan başarıyla çekildi!`);
      } else {
        setStatusMsg(`BAŞARILI: ${data.processed} adet haber Gemini AI ile başarıyla Türkçe'ye çevrildi!`);
      }
    } catch (err) {
      setStatusMsg(`HATA: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="console-widget">
      {/* Floating Action Trigger Button */}
      <button 
        className={`console-floating-btn ${isOpen ? 'active' : ''}`}
        onClick={() => setIsOpen(!isOpen)}
        title="Sistem Yönetim Konsolu"
      >
        <span className="pulse-indicator"></span>
        <span className="btn-icon">🦾</span>
        <span className="btn-text">KONSOL</span>
      </button>

      {/* Side Slide-over Panel */}
      <div className={`console-panel ${isOpen ? 'open' : ''}`}>
        <div className="console-panel__header">
          <h3>🦾 HaberBot Yönetim</h3>
          <button className="close-btn" onClick={() => setIsOpen(false)}>×</button>
        </div>

        <div className="console-panel__body">
          <div className="console-section">
            <h4 className="section-label">⚡ Sistem Durumu</h4>
            <div className="status-grid">
              <div className="status-item">
                <span>Veritabanı:</span>
                <span className="status-val active">Aktif (Neon)</span>
              </div>
              <div className="status-item">
                <span>Yapay Zeka:</span>
                <span className="status-val active">Gemini 2.0</span>
              </div>
              <div className="status-item">
                <span>Akışlar:</span>
                <span className="status-val active">HN & RSS</span>
              </div>
            </div>
          </div>

          <div className="console-section">
            <h4 className="section-label">🔑 Admin Yetkilendirme</h4>
            <p className="section-help">İşlemleri tetiklemek için .env dosyasındaki ADMIN_API_KEY değerini girin.</p>
            <input 
              type="password" 
              placeholder="Admin API Anahtarı..." 
              value={apiKey}
              onChange={handleKeyChange}
              className="console-input"
            />
          </div>

          <div className="console-section actions-section">
            <h4 className="section-label">⚙️ Boru Hatları (Pipelines)</h4>
            
            <button 
              className="console-btn fetch-btn"
              disabled={isLoading}
              onClick={() => triggerAction('fetch')}
            >
              ⚡ HABER ÇEKİMİ TETİKLE
            </button>

            <button 
              className="console-btn process-btn"
              disabled={isLoading}
              onClick={() => triggerAction('process')}
            >
              🤖 AI ÖZETLEME TETİKLE
            </button>
          </div>

          {statusMsg && (
            <div className={`console-log ${statusMsg.startsWith('HATA') ? 'error' : statusMsg.startsWith('BAŞARILI') ? 'success' : 'info'}`}>
              <div className="log-title">{statusMsg.startsWith('HATA') ? '🛑 HATA RAPORU' : statusMsg.startsWith('BAŞARILI') ? '✅ İŞLEM BAŞARILI' : '⏳ SİSTEM LOGU'}</div>
              <p className="log-content">{statusMsg}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ConsoleWidget;
