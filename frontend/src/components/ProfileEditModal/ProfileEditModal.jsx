import React, { useState, useEffect, useRef } from 'react';
import { api } from '../../services/api';
import { getAvatarUrl } from '../../utils/constants';
import './ProfileEditModal.css';

export default function ProfileEditModal({ profile, onClose, onSuccess }) {
  const [bio, setBio] = useState(profile.bio || '');
  const [selectedPreset, setSelectedPreset] = useState('');
  const [avatarFile, setAvatarFile] = useState(null);
  const [imagePreview, setImagePreview] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  
  const fileInputRef = useRef(null);

  // Default solid color preset names/URLs
  const presets = [
    { name: 'preset-1', color: '#2d2d2d', label: 'Mat Kömür' },
    { name: 'preset-2', color: '#6366f1', label: 'Gece Mavisi' },
    { name: 'preset-3', color: '#ef4444', label: 'Volkan Kırmızısı' },
    { name: 'preset-4', color: '#10b981', label: 'Zümrüt Yeşili' },
    { name: 'preset-5', color: '#f59e0b', label: 'Günışığı Kehribar' },
    { name: 'preset-6', color: '#8b5cf6', label: 'Uluslararası Mor' },
    { name: 'preset-7', color: '#ec4899', label: 'Neon Pembe' },
    { name: 'preset-8', color: '#14b8a6', label: 'Siber Turkuaz' },
  ];

  useEffect(() => {
    // Determine image preview based on initial avatar
    if (profile.avatar_url) {
      setImagePreview(getAvatarUrl(profile.avatar_url));
      
      // If active avatar is a preset, pre-select it
      const match = presets.find(p => profile.avatar_url.includes(p.name));
      if (match) {
        setSelectedPreset(match.name);
      }
    } else {
      setImagePreview(getAvatarUrl('/avatars/presets/preset-1.png'));
      setSelectedPreset('preset-1');
    }
  }, [profile]);

  // Handle Preset selection
  const handlePresetSelect = (presetName) => {
    setSelectedPreset(presetName);
    setAvatarFile(null); // Clear file upload
    setImagePreview(getAvatarUrl(`/avatars/presets/${presetName}.png`));
    setErrorMessage('');
  };

  // Handle custom file upload
  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    // Validate size (2MB)
    if (file.size > 2 * 1024 * 1024) {
      setErrorMessage('Resim boyutu en fazla 2 MB olmalıdır.');
      return;
    }

    // Validate MIME extension
    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
      setErrorMessage('Yalnızca JPG, PNG ve WebP dosyaları kabul edilir.');
      return;
    }

    setAvatarFile(file);
    setSelectedPreset(''); // Clear preset selection
    setErrorMessage('');

    // Generate object URL for preview
    const reader = new FileReader();
    reader.onloadend = () => {
      setImagePreview(reader.result);
    };
    reader.readAsDataURL(file);
  };

  // Handle Form Submission
  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSaving(true);
    setErrorMessage('');

    // Enforce 160 bio character limit
    if (bio.length > 160) {
      setErrorMessage('Biyografi 160 karakter sınırını aşamaz.');
      setIsSaving(false);
      return;
    }

    try {
      const formData = new FormData();
      formData.append('bio', bio);

      if (avatarFile) {
        formData.append('avatar', avatarFile);
      } else if (selectedPreset) {
        formData.append('preset_avatar', selectedPreset);
      }

      const updatedUser = await api.updateUserProfile(formData);
      
      // Save updated user to localstorage
      const currentStoredUser = localStorage.getItem('user');
      if (currentStoredUser) {
        const parsed = JSON.parse(currentStoredUser);
        localStorage.setItem('user', JSON.stringify({
          ...parsed,
          bio: updatedUser.bio,
          avatar_url: updatedUser.avatar_url,
        }));
      }

      onSuccess();
    } catch (err) {
      setErrorMessage(err.message || 'Profil güncellenirken hata oluştu.');
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="hb-modal-overlay" role="dialog" aria-modal="true">
      <div className="hb-modal-card">
        <div className="hb-modal-header">
          <h2 className="modal-title">Profili Düzenle</h2>
          <button className="modal-close" onClick={onClose} aria-label="Kapat">
            <span className="material-symbols-outlined">close</span>
          </button>
        </div>

        <form onSubmit={handleSubmit} className="hb-modal-form">
          {errorMessage && (
            <div className="modal-error-box">
              <span className="material-symbols-outlined text-[16px]">error</span>
              <span>{errorMessage}</span>
            </div>
          )}

          {/* Avatar Section */}
          <div className="modal-section">
            <label className="modal-label">Avatar Tercihi</label>
            <div className="avatar-edit-grid">
              {/* Preview */}
              <div className="avatar-preview-wrap">
                <img 
                  src={imagePreview} 
                  alt="Önizleme" 
                  className="modal-avatar-preview"
                  onError={(e) => {
                    e.target.src = getAvatarUrl('/avatars/presets/preset-1.png');
                  }}
                />
                <span className="preview-label">Görünüm</span>
              </div>

              {/* Upload & Files */}
              <div className="avatar-options-box">
                <button 
                  type="button" 
                  className="btn-upload-trigger"
                  onClick={() => fileInputRef.current?.click()}
                >
                  <span className="material-symbols-outlined">cloud_upload</span>
                  <span>Özel Resim Yükle</span>
                </button>
                <input 
                  type="file" 
                  ref={fileInputRef}
                  onChange={handleFileChange}
                  accept="image/jpeg,image/png,image/webp"
                  className="hidden-file-input"
                />
                <p className="upload-tip">JPG, PNG veya WebP. Max 2MB.</p>
              </div>
            </div>

            {/* Presets Grid */}
            <div className="presets-block">
              <span className="presets-title">Hazır Renk Paletleri</span>
              <div className="presets-grid">
                {presets.map((preset) => {
                  const isActive = selectedPreset === preset.name;
                  return (
                    <button
                      key={preset.name}
                      type="button"
                      onClick={() => handlePresetSelect(preset.name)}
                      className={`preset-item ${isActive ? 'is-active' : ''}`}
                      style={{ '--preset-color': preset.color }}
                      title={preset.label}
                    >
                      <span className="preset-color-dot"></span>
                      {isActive && (
                        <span className="material-symbols-outlined check-icon">check</span>
                      )}
                    </button>
                  );
                })}
              </div>
            </div>
          </div>

          {/* Bio Section */}
          <div className="modal-section mt-4">
            <div className="flex justify-between items-center mb-1">
              <label className="modal-label">Biyografi</label>
              <span className={`char-counter ${bio.length > 160 ? 'limit-exceeded' : ''}`}>
                {bio.length} / 160
              </span>
            </div>
            <textarea 
              value={bio}
              onChange={(e) => setBio(e.target.value)}
              placeholder="Bize kendinizden bahsedin..."
              maxLength={170}
              className="modal-bio-textarea"
              rows={3}
            />
          </div>

          {/* Footer controls */}
          <div className="hb-modal-footer">
            <button 
              type="button" 
              className="btn-modal-cancel" 
              onClick={onClose}
              disabled={isSaving}
            >
              İptal
            </button>
            <button 
              type="submit" 
              className="btn-modal-save" 
              disabled={isSaving || bio.length > 160}
            >
              {isSaving ? 'Kaydediliyor...' : 'Kaydet'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
