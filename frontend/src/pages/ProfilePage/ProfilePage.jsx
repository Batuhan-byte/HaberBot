import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link, useSearchParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../../services/api';
import { getAvatarUrl } from '../../utils/constants';
import Header from '../../components/Header/Header';
import ProfileEditModal from '../../components/ProfileEditModal/ProfileEditModal';
import ReportModal from '../../components/ReportModal/ReportModal';
import './ProfilePage.css';

export default function ProfilePage() {
  const { username } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [isEditOpen, setIsEditOpen] = useState(false);
  const [isReportOpen, setIsReportOpen] = useState(false);
  const [searchParams, setSearchParams] = useSearchParams();

  useEffect(() => {
    if (searchParams.get('edit') === 'true') {
      setIsEditOpen(true);
      setSearchParams({}, { replace: true });
    }
  }, [searchParams, setSearchParams]);

  // Get current logged in user from localStorage
  const currentUserJson = localStorage.getItem('user');
  const currentUser = currentUserJson ? JSON.parse(currentUserJson) : null;
  const isOwnProfile = currentUser && currentUser.username === username;

  // Query to fetch user profile
  const { data: profile, isLoading, isError, error } = useQuery({
    queryKey: ['profile', username],
    queryFn: () => api.fetchUserProfile(username),
    retry: false,
    staleTime: 2 * 60 * 1000,
  });

  // Mutate favorites (Add/Remove)
  const toggleFavoriteMutation = useMutation({
    mutationFn: async () => {
      if (!profile) return;
      if (profile.is_favorited) {
        return api.removeFavorite(profile.id);
      } else {
        return api.addFavorite(profile.id);
      }
    },
    // Optimistic UI updates
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ['profile', username] });
      const previousProfile = queryClient.getQueryData(['profile', username]);
      
      queryClient.setQueryData(['profile', username], (old) => {
        if (!old) return old;
        return {
          ...old,
          is_favorited: !old.is_favorited,
        };
      });

      return { previousProfile };
    },
    onError: (err, variables, context) => {
      if (context?.previousProfile) {
        queryClient.setQueryData(['profile', username], context.previousProfile);
      }
      alert(err.message || 'İşlem gerçekleştirilemedi.');
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ['profile', username] });
    },
  });

  // Handle loading state with premium modern skeletons
  if (isLoading) {
    return (
      <div className="hb-app-wrapper">
        <Header />
        <div className="profile-container animate-pulse">
          <div className="profile-header-skeleton">
            <div className="avatar-skeleton"></div>
            <div className="meta-skeleton">
              <div className="h-6 w-32 bg-[#1a1a1a] rounded"></div>
              <div className="h-4 w-48 bg-[#1a1a1a] rounded mt-2"></div>
              <div className="h-4 w-16 bg-[#1a1a1a] rounded mt-1"></div>
            </div>
          </div>
          <div className="profile-grid mt-8">
            <div className="profile-sidebar-skeleton">
              <div className="h-20 bg-[#121212] rounded border border-[#222]"></div>
              <div className="h-40 bg-[#121212] rounded border border-[#222] mt-4"></div>
            </div>
            <div className="profile-content-skeleton">
              <div className="h-64 bg-[#121212] rounded border border-[#222]"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Handle not found/error state
  if (isError || !profile) {
    const errorMsg = error?.message || 'Profil bulunamadı.';
    return (
      <div className="hb-app-wrapper">
        <Header />
        <div className="profile-error-container">
          <span className="material-symbols-outlined error-icon">person_off</span>
          <h2>Profil Bulunamadı</h2>
          <p>{errorMsg}</p>
          <button className="back-btn" onClick={() => navigate('/')}>
            Ana Sayfaya Dön
          </button>
        </div>
      </div>
    );
  }

  // Default preset or image logic
  const getAvatarSource = () => {
    if (profile.avatar_url) {
      // If it is preset path, convert to server URL or serve directly
      return getAvatarUrl(profile.avatar_url);
    }
    // Default fallback preset 1
    return getAvatarUrl("/avatars/presets/preset-1.png");
  };

  // Format join date nicely
  const formatJoinDate = (dateStr) => {
    try {
      const d = new Date(dateStr);
      return d.toLocaleDateString('tr-TR', { year: 'numeric', month: 'long', day: 'numeric' });
    } catch (e) {
      return 'Bilinmiyor';
    }
  };

  return (
    <div className="hb-app-wrapper">
      <Header />
      
      <div className="profile-container">
        {/* Profile Card Header */}
        <div className="profile-hero-section">
          <div className="profile-header-wrap">
            <div className="profile-avatar-wrap">
              <img 
                src={getAvatarSource()} 
                alt={profile.username} 
                className="profile-main-avatar"
                onError={(e) => {
                  e.target.src = getAvatarUrl("/avatars/presets/preset-1.png");
                }}
              />
              <div className="avatar-subtle-glow"></div>
            </div>

            <div className="profile-info-block">
              <div className="profile-title-row">
                <h1 className="profile-username">@{profile.username}</h1>
                {profile.role === 'Admin' && (
                  <span className="profile-role-badge">Yönetici</span>
                )}
              </div>
              
              <p className="profile-joined">
                <span className="material-symbols-outlined">calendar_today</span>
                <span>{formatJoinDate(profile.join_date)} tarihinde katıldı</span>
              </p>

              {/* Bio block */}
              <div className="profile-bio-box">
                {profile.bio ? (
                  <p className="bio-text">“{profile.bio}”</p>
                ) : (
                  <p className="bio-empty-text">Henüz bir biyografi yazılmamış.</p>
                )}
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="profile-actions-bar">
            {isOwnProfile ? (
              <button 
                className="btn-premium-action btn-edit" 
                onClick={() => setIsEditOpen(true)}
              >
                <span className="material-symbols-outlined text-[18px]">edit</span>
                <span>Profili Düzenle</span>
              </button>
            ) : (
              <div className="flex gap-2">
                {/* Favorites button */}
                <button 
                  className={`btn-premium-action btn-favorite ${profile.is_favorited ? 'is-active' : ''}`}
                  onClick={() => toggleFavoriteMutation.mutate()}
                >
                  <span className="material-symbols-outlined text-[18px]">
                    {profile.is_favorited ? 'star' : 'star_border'}
                  </span>
                  <span>{profile.is_favorited ? 'Favorilerden Çıkar' : 'Favorilere Ekle'}</span>
                </button>

                {/* Report button */}
                <button 
                  className="btn-premium-action btn-report"
                  onClick={() => setIsReportOpen(true)}
                >
                  <span className="material-symbols-outlined text-[18px]">flag</span>
                  <span>Şikayet Et</span>
                </button>
              </div>
            )}
          </div>
        </div>

        {/* Profile Details Grid */}
        <div className="profile-stats-grid mt-6">
          {/* Stats Widget */}
          <div className="profile-stats-card">
            <div className="stat-item">
              <span className="stat-label">Onaylanan Haberler</span>
              <span className="stat-val">{profile.articles_published}</span>
            </div>
            <div className="stat-divider"></div>
            <div className="stat-item">
              <span className="stat-label">Rol Sınıfı</span>
              <span className="stat-val-text">{profile.role}</span>
            </div>
          </div>

          {/* System status details or activity panel */}
          <div className="profile-activity-card">
            <h3 className="card-subtitle">
              <span className="material-symbols-outlined text-[16px] text-blue-500">verified_user</span>
              <span>Platform Güvenliği</span>
            </h3>
            <p className="card-description">
              HaberBot topluluğunun saygın bir üyesidir. Paylaşımları otomatik filtreleme ve yapay zeka analizine tabi tutulur.
            </p>
          </div>
        </div>
      </div>

      {/* Profile Edit Modal */}
      {isEditOpen && (
        <ProfileEditModal 
          profile={profile} 
          onClose={() => setIsEditOpen(false)} 
          onSuccess={() => {
            setIsEditOpen(false);
            queryClient.invalidateQueries({ queryKey: ['profile', username] });
            // Invalidate current user cache in header
            window.dispatchEvent(new CustomEvent('user-auth-changed'));
          }}
        />
      )}

      {/* Profile Report Modal */}
      {isReportOpen && (
        <ReportModal 
          reportedUserId={profile.id}
          reportedUsername={profile.username}
          onClose={() => setIsReportOpen(false)}
        />
      )}
    </div>
  );
}
