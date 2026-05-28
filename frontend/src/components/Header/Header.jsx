import React, { useState, useEffect, useRef } from 'react';
import { Link, NavLink, useNavigate, useLocation } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { api } from '../../services/api';
import { formatRelativeDate } from '../../utils/formatDate';
import { useDebounce } from '../../hooks/useDebounce';
import SearchDropdown from '../SearchDropdown/SearchDropdown';
import './Header.css';

const FALLBACK_IMAGE_ALT = "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?auto=format&fit=crop&w=600&q=80";

function stripHtmlTags(str) {
  if (!str) return '';
  return str.replace(/<\/?[^>]+(>|$)/g, "").trim();
}

export default function Header({ searchQuery = '', setSearchQuery, selectedTopic = 'Tümü' }) {
  const navigate = useNavigate();
  const location = useLocation();
  const headerRef = useRef(null);
  
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isLeftDrawerOpen, setIsLeftDrawerOpen] = useState(false);
  const [activeDropdownSlug, setActiveDropdownSlug] = useState(null);
  const [localSearch, setLocalSearch] = useState(searchQuery);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const [theme, setTheme] = useState('dark');
  const [showDropdown, setShowDropdown] = useState(false);

  useEffect(() => {
    const adminKey = localStorage.getItem('admin_api_key');
    setIsLoggedIn(!!adminKey);
  }, [location.pathname]);

  // Load and apply theme
  useEffect(() => {
    const savedTheme = localStorage.getItem('theme') || 'dark';
    setTheme(savedTheme);
    if (savedTheme === 'light') {
      document.body.classList.add('light-mode');
    } else {
      document.body.classList.remove('light-mode');
    }
  }, []);

  const toggleTheme = () => {
    const nextTheme = theme === 'dark' ? 'light' : 'dark';
    setTheme(nextTheme);
    localStorage.setItem('theme', nextTheme);
    if (nextTheme === 'light') {
      document.body.classList.add('light-mode');
    } else {
      document.body.classList.remove('light-mode');
    }
  };

  // Synchronize local search state with search query prop
  useEffect(() => {
    setLocalSearch(searchQuery);
  }, [searchQuery]);

  // Debounced search for dropdown suggestions (300ms delay)
  const debouncedSearch = useDebounce(localSearch, 300);

  // Fetch topics
  const { data: topicsData } = useQuery({
    queryKey: ['topics'],
    queryFn: api.getTopics
  });
  const topics = topicsData?.topics || [];

  // Fetch drawer articles on-demand based on activeDropdownSlug
  const { data: drawerArticlesData, isLoading: isDrawerLoading } = useQuery({
    queryKey: ['drawerArticles', activeDropdownSlug],
    queryFn: () => api.getTopicArticles(activeDropdownSlug, 1, 5),
    enabled: !!activeDropdownSlug,
    staleTime: 2 * 60 * 1000 // 2 minutes cache
  });

  const drawerArticles = drawerArticlesData?.articles || [];

  // Handle outside clicks to close the dropdown panel, mobile menu and search wrapper
  useEffect(() => {
    function handleClickOutside(event) {
      if (headerRef.current && !headerRef.current.contains(event.target)) {
        setActiveDropdownSlug(null);
        setIsMobileMenuOpen(false);
        setIsSearchOpen(false);
        setShowDropdown(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // Close menus on path changes
  useEffect(() => {
    setActiveDropdownSlug(null);
    setIsMobileMenuOpen(false);
    setIsLeftDrawerOpen(false);
  }, [location.pathname]);

  // Disable page scroll when side drawer is open
  useEffect(() => {
    if (isLeftDrawerOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isLeftDrawerOpen]);

  const handleTopicClick = (e, topic) => {
    // If desktop view: toggle mega drawer
    if (window.innerWidth >= 1024) {
      e.preventDefault(); // Prevent direct routing on first click
      if (activeDropdownSlug === topic.slug) {
        // Second click: close drawer and navigate to category page
        setActiveDropdownSlug(null);
        navigate(`/kategori/${topic.slug}`);
      } else {
        // First click: open mega drawer
        setActiveDropdownSlug(topic.slug);
      }
    } else {
      // Mobile view: go directly to category page
      setIsMobileMenuOpen(false);
    }
  };

  const handleSearchSubmit = (e) => {
    if (e.key === 'Enter' || e.type === 'click') {
      if (localSearch.trim()) {
        navigate(`/arama?q=${encodeURIComponent(localSearch)}`);
        setShowDropdown(false);
      }
      setIsMobileMenuOpen(false);
    }
  };

  const activeTopicName = topics.find(t => t.slug === activeDropdownSlug)?.name || 'Kategori';

  return (
    <>
      <header className="hb-header" ref={headerRef}>
        <nav className="hb-nav" aria-label="Ana Navigasyon">
          <div className="flex items-center gap-4 h-full">
            {/* Hamburger Menu Icon */}
            <button 
              className="hb-hamburger-btn" 
              onClick={() => setIsLeftDrawerOpen(true)}
              aria-label="Sol Hamburger Menüyü Aç"
            >
              <span className="material-symbols-outlined text-[22px]">menu</span>
            </button>
            
            <Link to="/" className="hb-logo" aria-label="HaberBot Logo" onClick={() => setActiveDropdownSlug(null)}>
              Haber<span>Bot</span>
            </Link>
            
            {/* ========== DESKTOP TOPICS NAVIGATION ========== */}
            <div className="hidden lg:flex items-center gap-1 h-full" role="menubar">
              {topics.slice(0, 8).map(topic => {
                const isDrawerActive = activeDropdownSlug === topic.slug;
                const isPageActive = selectedTopic === topic.name && !activeDropdownSlug;
                return (
                  <a
                    key={topic.id}
                    href={`/kategori/${topic.slug}`}
                    onClick={(e) => handleTopicClick(e, topic)}
                    className={`menu-item ${isDrawerActive ? 'drawer-active' : ''} ${isPageActive ? 'active' : ''}`}
                    role="menuitem"
                  >
                    {topic.name}
                  </a>
                );
              })}
            </div>
          </div>

          <div className="nav-controls">
            {/* Animasyonlu Arama Kutusu */}
            <div className={`hb-search-wrapper ${isSearchOpen ? 'is-open' : ''}`}>
              <button 
                className="hb-search-trigger-btn"
                onClick={() => {
                  setIsSearchOpen(!isSearchOpen);
                  setShowDropdown(false);
                }}
                aria-label="Arama Yap"
              >
                <span className="material-symbols-outlined">search</span>
              </button>
              <div className="hb-search-input-container">
                <input 
                  placeholder="Haberlerde ara..." 
                  type="text"
                  value={localSearch}
                  onChange={(e) => {
                    setLocalSearch(e.target.value);
                    setShowDropdown(true);
                  }}
                  onKeyDown={handleSearchSubmit}
                  onFocus={() => setShowDropdown(true)}
                  onBlur={() => setTimeout(() => setShowDropdown(false), 200)}
                />
                {localSearch && (
                  <button 
                    onClick={() => {
                      setLocalSearch('');
                      if (setSearchQuery) setSearchQuery('');
                      setShowDropdown(false);
                    }}
                    className="material-symbols-outlined text-gray-500 hover:text-white text-[14px] ml-1"
                  >
                    close
                  </button>
                )}
              </div>
              {showDropdown && debouncedSearch.length >= 2 && (
                <SearchDropdown 
                  query={debouncedSearch}
                  onSelect={(article) => {
                    navigate(`/haber/${article.id}`);
                    setLocalSearch('');
                    setShowDropdown(false);
                    setIsSearchOpen(false);
                  }}
                />
              )}
            </div>

            {/* Nav Actions Group */}
            <div className="hb-nav-actions">
              {/* Favorilerim Butonu */}
              <button className="hb-action-icon-btn" aria-label="Favorilerim" title="Favorilerim">
                <span className="material-symbols-outlined">bookmarks</span>
              </button>

              {/* Gazete Butonu */}
              <button className="hb-action-icon-btn" aria-label="Haberler" title="Haberler">
                <span className="material-symbols-outlined">newspaper</span>
              </button>

              {/* Üye Girişi / Profil Yönetimi */}
              {isLoggedIn ? (
                <div className="flex items-center gap-1">
                  <Link to="/admin" className="hb-action-icon-btn" aria-label="Admin" title="Yönetici Ayarları">
                    <span className="material-symbols-outlined">settings</span>
                  </Link>
                  <button 
                    onClick={() => {
                      localStorage.removeItem('admin_api_key');
                      setIsLoggedIn(false);
                      navigate('/');
                    }}
                    className="hb-action-icon-btn btn-logout" 
                    title="Çıkış Yap"
                  >
                    <span className="material-symbols-outlined">logout</span>
                  </button>
                  <Link to="/admin" className="hb-profile-avatar" title="Yönetici Paneli">
                    <img 
                      src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=facearea&facepad=2&w=256&h=256&q=80" 
                      alt="Profil" 
                    />
                  </Link>
                </div>
              ) : (
                <Link to="/login" className="hb-premium-login-btn">
                  <span className="material-symbols-outlined">person</span>
                  <span>ÜYE GİRİŞİ</span>
                </Link>
              )}

              {/* Koyu/Açık Mod Switcher Toggle */}
              <button 
                className={`hb-theme-toggle ${theme === 'light' ? 'is-light' : ''}`}
                onClick={toggleTheme}
                aria-label="Tema Değiştir"
                title={theme === 'dark' ? "Açık Moda Geç" : "Koyu Moda Geç"}
              >
                <div className="hb-theme-toggle-thumb">
                  <span className="material-symbols-outlined">
                    {theme === 'dark' ? 'dark_mode' : 'light_mode'}
                  </span>
                </div>
              </button>
            </div>
            
            {/* Mobile Hamburger Toggle */}
            <button 
              className={`mobile-menu-toggle lg:hidden ${isMobileMenuOpen ? 'is-active' : ''}`}
              onClick={() => {
                setIsMobileMenuOpen(!isMobileMenuOpen);
                setActiveDropdownSlug(null); // Close active mega drawer
              }}
              aria-label="Menüyü Aç/Kapat"
            >
              <span className="bar"></span>
              <span className="bar"></span>
              <span className="bar"></span>
            </button>
          </div>
        </nav>

        {/* ========== RED UNDERLINE ACCENT BORDER ========== */}
        <div className="hb-bottom-accent"></div>

        {/* ========== MOBILE MENU DRAWER ========== */}
        <div className={`mobile-nav-drawer lg:hidden ${isMobileMenuOpen ? 'is-open' : ''}`}>
          <div className="mobile-nav-links">
            <Link 
              to="/" 
              className={`mobile-nav-item ${selectedTopic === 'Tümü' ? 'active' : ''}`}
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Tümü
            </Link>
            {topics.map(topic => (
              <Link 
                key={topic.id}
                to={`/kategori/${topic.slug}`}
                className={`mobile-nav-item ${selectedTopic === topic.name ? 'active' : ''}`}
                onClick={() => setIsMobileMenuOpen(false)}
              >
                {topic.name}
              </Link>
            ))}
            <div className="mobile-nav-divider"></div>
            <Link to="/admin" className="mobile-nav-item flex items-center gap-2" onClick={() => setIsMobileMenuOpen(false)}>
              <span className="material-symbols-outlined text-[18px]">settings</span> Yönetici Ayarları
            </Link>
          </div>
        </div>

        {/* ========== DESKTOP SLIDING MEGA DROPDOWN PANEL ========== */}
        <div className={`hb-dropdown-panel ${activeDropdownSlug ? 'is-open' : ''}`}>
          <div className="dropdown-panel-container">
            <div className="dropdown-header">
              <div className="flex items-center gap-2">
                <span className="dropdown-glow-dot"></span>
                <h4 className="dropdown-title">
                  Son <span>{activeTopicName}</span> Gelişmeleri
                </h4>
              </div>
              <button 
                onClick={() => setActiveDropdownSlug(null)} 
                className="dropdown-close-btn flex items-center gap-1 hover:text-white"
              >
                <span className="text-[12px] font-bold">KAPAT</span>
                <span className="material-symbols-outlined text-[16px]">close</span>
              </button>
            </div>

            {isDrawerLoading ? (
              <div className="dropdown-loading">
                <div className="spinner"></div>
                <span>İçerikler yükleniyor...</span>
              </div>
            ) : drawerArticles.length === 0 ? (
              <div className="dropdown-empty">
                <span className="material-symbols-outlined text-[32px] text-gray-600 mb-2">article</span>
                <span>Bu kategoriye ait henüz haber bulunmamaktadır.</span>
              </div>
            ) : (
              <div className="dropdown-content-grid">
                {/* Left Side: Featured Mega Article Card */}
                {drawerArticles[0] && (() => {
                  const art = drawerArticles[0];
                  const displayTitle = art.turkish_title || art.title;
                  const cleanSummary = stripHtmlTags(art.turkish_summary || art.summary || art.original_content).slice(0, 110);
                  return (
                    <Link 
                      to={`/haber/${art.id}`} 
                      className="dropdown-featured-card group"
                      onClick={() => setActiveDropdownSlug(null)}
                    >
                      <div className="featured-card-img-wrap">
                        <img 
                          src={art.image_url || FALLBACK_IMAGE_ALT} 
                          alt={displayTitle} 
                        />
                        <div className="featured-card-overlay"></div>
                      </div>
                      <div className="featured-card-body">
                        <span className="featured-card-badge">{activeTopicName.toUpperCase()}</span>
                        <h3 className="featured-card-title line-clamp-2">{displayTitle}</h3>
                        <p className="featured-card-desc line-clamp-2">{cleanSummary}...</p>
                        <div className="featured-card-footer">
                          <span>Deniz Şen</span>
                          <span className="footer-dot">•</span>
                          <span>{formatRelativeDate(art.created_at || art.fetched_at)}</span>
                        </div>
                      </div>
                    </Link>
                  );
                })()}

                {/* Right Side: List of 4 other items */}
                <div className="dropdown-side-list">
                  {drawerArticles.slice(1, 5).map(art => {
                    const displayTitle = art.turkish_title || art.title;
                    return (
                      <Link 
                        key={art.id} 
                        to={`/haber/${art.id}`} 
                        className="dropdown-side-item group"
                        onClick={() => setActiveDropdownSlug(null)}
                      >
                        <div className="side-item-body">
                          <h4 className="side-item-title line-clamp-2">{displayTitle}</h4>
                          <span className="side-item-time">{formatRelativeDate(art.created_at || art.fetched_at)}</span>
                        </div>
                        <div className="side-item-img-wrap">
                          <img 
                            src={art.image_url || FALLBACK_IMAGE_ALT} 
                            alt={displayTitle} 
                          />
                        </div>
                      </Link>
                    );
                  })}
                </div>
              </div>
            )}

            {/* Bottom Bar Controls */}
            <div className="dropdown-footer-bar">
              <div className="flex items-center gap-2 text-gray-500 text-xs">
                <span className="material-symbols-outlined text-[14px] text-blue-500">auto_awesome</span>
                <span>HaberBot Yapay Zekâ Destekli Haber Portalı</span>
              </div>
              <Link 
                to={`/kategori/${activeDropdownSlug}`} 
                className="dropdown-view-all-btn group"
                onClick={() => setActiveDropdownSlug(null)}
              >
                TÜM {activeTopicName.toUpperCase()} GELEŞMELERİNİ GÖR 
                <span className="material-symbols-outlined transition-transform group-hover:translate-x-1">arrow_forward</span>
              </Link>
            </div>
          </div>
        </div>
      </header>

      {/* ========== SLIDING LEFT SIDE HAMBURGER DRAWER ========== */}
      <div className={`drawer-overlay ${isLeftDrawerOpen ? 'is-visible' : ''}`} onClick={() => setIsLeftDrawerOpen(false)}></div>
      
      <div className={`hb-left-drawer ${isLeftDrawerOpen ? 'is-open' : ''}`}>
        {/* User Account / Close Button Bar */}
        <div className="left-drawer-header">
          <Link to="/login" className="drawer-user-info" onClick={() => setIsLeftDrawerOpen(false)}>
            <div className="drawer-user-avatar">
              <span className="material-symbols-outlined text-[20px]">person</span>
            </div>
            <div className="drawer-user-text flex flex-col">
              <span className="drawer-username">Giriş Yap / Üye Ol</span>
            </div>
          </Link>
          <button 
            className="drawer-close-btn" 
            onClick={() => setIsLeftDrawerOpen(false)}
            aria-label="Menüyü Kapat"
          >
            <span className="material-symbols-outlined text-[22px]">close</span>
          </button>
        </div>

        {/* Navigation Core Menu */}
        <div className="left-drawer-menu">
          <Link to="/" className="drawer-menu-item" onClick={() => setIsLeftDrawerOpen(false)}>
            <span className="material-symbols-outlined">power_settings_new</span>
            <span>Anasayfa</span>
          </Link>
          <Link to="/" className="drawer-menu-item" onClick={() => setIsLeftDrawerOpen(false)}>
            <span className="material-symbols-outlined">newspaper</span>
            <span>Haber</span>
          </Link>
          <a href="#" className="drawer-menu-item" onClick={(e) => { e.preventDefault(); setIsLeftDrawerOpen(false); }}>
            <span className="material-symbols-outlined">play_circle</span>
            <span>Video</span>
          </a>
          <a href="#" className="drawer-menu-item" onClick={(e) => { e.preventDefault(); setIsLeftDrawerOpen(false); }}>
            <span className="material-symbols-outlined">bookmarks</span>
            <span>Favorilerim</span>
          </a>
          <div className="drawer-menu-item dropdown-toggle">
            <span className="material-symbols-outlined">bar_chart</span>
            <span>En Çok Okunanlar</span>
            <span className="material-symbols-outlined ml-auto text-[14px]">keyboard_arrow_down</span>
          </div>
          <div className="drawer-menu-item dropdown-toggle">
            <span className="material-symbols-outlined">trending_up</span>
            <span>En Çok Paylaşılanlar</span>
            <span className="material-symbols-outlined ml-auto text-[14px]">keyboard_arrow_down</span>
          </div>
          <div className="drawer-menu-item dropdown-toggle">
            <span className="material-symbols-outlined">smart_display</span>
            <span>En Çok İzlenenler</span>
            <span className="material-symbols-outlined ml-auto text-[14px]">keyboard_arrow_down</span>
          </div>
        </div>

        {/* Categories Section */}
        <div className="left-drawer-section">
          <h4 className="section-title">Kategoriler</h4>
          <div className="drawer-categories-list">
            <Link to="/" className="category-link" onClick={() => setIsLeftDrawerOpen(false)}>Tümü</Link>
            {topics.map(topic => (
              <Link 
                key={topic.id} 
                to={`/kategori/${topic.slug}`} 
                className="category-link"
                onClick={() => setIsLeftDrawerOpen(false)}
              >
                {topic.name}
              </Link>
            ))}
          </div>
        </div>

        {/* Social Media Link Icons */}
        <div className="left-drawer-socials">
          <a href="#" className="social-icon text-[#3b5998]" aria-label="Facebook">
            <span className="material-symbols-outlined text-[20px]">public</span>
          </a>
          <a href="#" className="social-icon text-[#ffffff]" aria-label="Twitter">
            <span className="material-symbols-outlined text-[20px]">terminal</span>
          </a>
          <a href="#" className="social-icon text-[#c13584]" aria-label="Instagram">
            <span className="material-symbols-outlined text-[20px]">photo_camera</span>
          </a>
          <a href="#" className="social-icon text-[#0088cc]" aria-label="Telegram">
            <span className="material-symbols-outlined text-[20px]">send</span>
          </a>
          <a href="#" className="social-icon text-[#ff0000]" aria-label="YouTube">
            <span className="material-symbols-outlined text-[20px]">smart_display</span>
          </a>
        </div>

        {/* Secondary Corporate Links */}
        <div className="left-drawer-corporate">
          <a href="#" className="corporate-link">Hakkımızda</a>
          <a href="#" className="corporate-link">Yazarlar</a>
          <a href="#" className="corporate-link">Ödüllerimiz</a>
          <a href="#" className="corporate-link">Künye</a>
          <a href="#" className="corporate-link">Gizlilik</a>
          <a href="#" className="corporate-link">İletişim</a>
        </div>

        {/* Mediazone Branding */}
        <div className="left-drawer-branding">
          <span>Bir <span>mediazone</span> markasıdır.</span>
        </div>

        {/* Partner Logos Grid at the bottom */}
        <div className="left-drawer-partners">
          <span className="partner-logo">mackolik</span>
          <span className="partner-logo match-direct">MATCH EN DIRECT</span>
          <span className="partner-logo">mynet</span>
          <span className="partner-logo">onedio</span>
          <span className="partner-logo">HaberBot</span>
          <span className="partner-logo">Yemek.com</span>
        </div>
      </div>
    </>
  );
}
