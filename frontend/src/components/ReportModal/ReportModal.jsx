import React, { useState } from 'react';
import { api } from '../../services/api';
import './ReportModal.css';

export default function ReportModal({ reportedUserId, reportedUsername, onClose }) {
  const [reason, setReason] = useState('spam');
  const [comment, setComment] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setErrorMessage('');
    setSuccessMessage('');

    try {
      await api.reportProfile(reportedUserId, reason, comment);
      setSuccessMessage('Şikayetiniz başarıyla iletildi. Moderatör ekibimiz en kısa sürede inceleyecektir.');
      setTimeout(() => {
        onClose();
      }, 3000);
    } catch (err) {
      setErrorMessage(err.message || 'Şikayet iletilemedi. Lütfen tekrar deneyin.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="hb-modal-overlay" role="dialog" aria-modal="true">
      <div className="hb-modal-card">
        <div className="hb-modal-header">
          <h2 className="modal-title">Profili Şikayet Et</h2>
          <button className="modal-close" onClick={onClose} aria-label="Kapat" disabled={isSubmitting}>
            <span className="material-symbols-outlined">close</span>
          </button>
        </div>

        <form onSubmit={handleSubmit} className="hb-modal-form">
          {successMessage ? (
            <div className="modal-success-box">
              <span className="material-symbols-outlined success-icon">check_circle</span>
              <p className="success-text">{successMessage}</p>
            </div>
          ) : (
            <>
              {errorMessage && (
                <div className="modal-error-box">
                  <span className="material-symbols-outlined text-[16px]">error</span>
                  <span>{errorMessage}</span>
                </div>
              )}

              <p className="report-target-tip">
                Şu kullanıcının profilini şikayet ediyorsunuz: <span className="target-username">@{reportedUsername}</span>
              </p>

              {/* Reason dropdown */}
              <div className="modal-section mb-4">
                <label className="modal-label">Şikayet Nedeni</label>
                <select 
                  value={reason} 
                  onChange={(e) => setReason(e.target.value)}
                  className="modal-select-input"
                >
                  <option value="spam">Spam / İstenmeyen İçerik</option>
                  <option value="offensive_content">Rahatsız Edici / Müstehcen Görsel</option>
                  <option value="harassment">Taciz / Kötü Davranış</option>
                  <option value="other">Diğer Nedenler</option>
                </select>
              </div>

              {/* Comments */}
              <div className="modal-section mb-4">
                <label className="modal-label">Açıklama (İsteğe Bağlı)</label>
                <textarea 
                  value={comment}
                  onChange={(e) => setComment(e.target.value)}
                  placeholder="Şikayetinizi detaylandırın..."
                  maxLength={500}
                  className="modal-textarea-input"
                  rows={3}
                />
              </div>

              {/* Footer controls */}
              <div className="hb-modal-footer">
                <button 
                  type="button" 
                  className="btn-modal-cancel" 
                  onClick={onClose}
                  disabled={isSubmitting}
                >
                  İptal
                </button>
                <button 
                  type="submit" 
                  className="btn-modal-submit" 
                  disabled={isSubmitting}
                >
                  {isSubmitting ? 'Gönderiliyor...' : 'Şikayet Et'}
                </button>
              </div>
            </>
          )}
        </form>
      </div>
    </div>
  );
}
