import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../../services/api';

export default function AuthModal() {
  const [isOpen, setIsOpen] = useState(false);
  const [activeTab, setActiveTab] = useState('login'); // 'login' | 'register'
  
  // Login States
  const [loginUsername, setLoginUsername] = useState('');
  const [loginPassword, setLoginPassword] = useState('');
  const [showLoginPassword, setShowLoginPassword] = useState(false);
  const [loginError, setLoginError] = useState('');
  const [loginLoading, setLoginLoading] = useState(false);

  // Register States
  const [registerUsername, setRegisterUsername] = useState('');
  const [registerEmail, setRegisterEmail] = useState('');
  const [registerPassword, setRegisterPassword] = useState('');
  const [showRegisterPassword, setShowRegisterPassword] = useState(false);
  const [registerConfirmPassword, setRegisterConfirmPassword] = useState('');
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [registerError, setRegisterError] = useState('');
  const [registerSuccess, setRegisterSuccess] = useState(false);
  const [registerLoading, setRegisterLoading] = useState(false);

  const navigate = useNavigate();

  // Listen to global open event
  useEffect(() => {
    const handleOpen = (e) => {
      setIsOpen(true);
      if (e.detail && e.detail.tab) {
        setActiveTab(e.detail.tab);
      }
      // Reset errors & fields when opening
      setLoginError('');
      setRegisterError('');
      setRegisterSuccess(false);
    };

    window.addEventListener('open-auth-modal', handleOpen);
    return () => window.removeEventListener('open-auth-modal', handleOpen);
  }, []);

  // Handle escape key and overflow on body
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    
    const handleKeyDown = (e) => {
      if (e.key === 'Escape') setIsOpen(false);
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => {
      document.body.style.overflow = '';
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [isOpen]);

  const handleLoginSubmit = async (e) => {
    e.preventDefault();
    setLoginError('');
    setLoginLoading(true);

    try {
      const data = await api.login(loginUsername, loginPassword);
      setIsOpen(false);
      
      // Dispatch custom event to notify components that user has logged in
      window.dispatchEvent(new CustomEvent('user-auth-changed'));
      
      if (data.user.role === 'Admin') {
        navigate('/admin');
      } else {
        // Reload or stay on the current page with logged-in state
        window.location.reload();
      }
    } catch (err) {
      setLoginError(err.message || 'Giriş başarısız. Lütfen bilgilerinizi kontrol edin.');
    } finally {
      setLoginLoading(false);
    }
  };

  const handleRegisterSubmit = async (e) => {
    e.preventDefault();
    setRegisterError('');

    // 1. Username length validation
    if (registerUsername.length < 3) {
      setRegisterError('Kullanıcı adı en az 3 karakter uzunluğunda olmalıdır.');
      return;
    }

    // 2. Email format validation (Regex pattern)
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(registerEmail)) {
      setRegisterError('Lütfen geçerli bir e-posta adresi girin (Örn: ad@ornek.com).');
      return;
    }

    // 3. Password length validation
    if (registerPassword.length < 8) {
      setRegisterError('Şifre en az 8 karakter uzunluğunda olmalıdır.');
      return;
    }

    // 4. Password match validation
    if (registerPassword !== registerConfirmPassword) {
      setRegisterError('Şifreler uyuşmuyor.');
      return;
    }

    setRegisterLoading(true);
    try {
      await api.register(registerUsername, registerEmail, registerPassword);
      setRegisterSuccess(true);
      setTimeout(() => {
        // Switch to login tab and prefill username
        setActiveTab('login');
        setLoginUsername(registerUsername);
        setRegisterUsername('');
        setRegisterEmail('');
        setRegisterPassword('');
        setRegisterConfirmPassword('');
        setRegisterSuccess(false);
      }, 1500);
    } catch (err) {
      setRegisterError(err.message || 'Kayıt başarısız. Lütfen başka bir kullanıcı adı deneyin.');
    } finally {
      setRegisterLoading(false);
    }
  };

  // Helper to calculate password strength dynamically
  const getPasswordStrength = (pass) => {
    if (!pass) return { label: '', colorClass: '', widthClass: 'w-0', textClass: 'text-gray-500' };
    if (pass.length < 4) return { label: 'Çok Kısa', colorClass: 'bg-red-600 shadow-[0_0_8px_rgba(220,38,38,0.5)]', widthClass: 'w-1/4', textClass: 'text-red-500' };
    if (pass.length < 8) return { label: 'Zayıf', colorClass: 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.5)]', widthClass: 'w-2/4', textClass: 'text-red-400' };
    
    // Check if it has numbers, letters and symbols
    const hasLetters = /[a-zA-ZğüşıöçĞÜŞİÖÇ]/.test(pass);
    const hasNumbers = /[0-9]/.test(pass);
    const hasSpecial = /[^a-zA-Z0-9ğüşıöçĞÜŞİÖÇ]/.test(pass);
    
    if (hasLetters && hasNumbers && hasSpecial) {
      return { label: 'Çok Güçlü', colorClass: 'bg-cyan-400 shadow-[0_0_10px_rgba(34,211,238,0.6)]', widthClass: 'w-full', textClass: 'text-cyan-400' };
    }
    if (hasLetters && (hasNumbers || hasSpecial)) {
      return { label: 'Güçlü', colorClass: 'bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.6)]', widthClass: 'w-3/4', textClass: 'text-blue-400' };
    }
    
    return { label: 'Orta', colorClass: 'bg-orange-500 shadow-[0_0_8px_rgba(249,115,22,0.5)]', widthClass: 'w-2/4', textClass: 'text-orange-400' };
  };

  const strength = getPasswordStrength(registerPassword);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[9999] flex items-center justify-center p-4">
      {/* Dark Ambient Glassmorphic Backdrop */}
      <div 
        className="absolute inset-0 bg-black/85 backdrop-blur-md transition-opacity duration-300"
        onClick={() => setIsOpen(false)}
      ></div>

      {/* Premium Glassmorphism Auth Modal Container */}
      <div className="relative bg-[#050505]/95 border border-blue-500/20 max-w-[440px] w-full rounded-2xl shadow-2xl overflow-hidden z-[10000] p-8 sm:p-10 animate-fade-in transition-all duration-300 hover:border-blue-500/30">
        
        {/* Subtle top border illumination */}
        <div className="absolute top-0 inset-x-0 h-[1px] bg-gradient-to-r from-transparent via-blue-500/40 to-transparent"></div>
        
        {/* Premium Ambient Background Glowing Blobs inside modal */}
        <div className="absolute -top-12 -left-12 w-48 h-48 bg-blue-600/10 rounded-full blur-[80px] pointer-events-none"></div>
        <div className="absolute -bottom-16 -right-16 w-56 h-56 bg-indigo-600/8 rounded-full blur-[90px] pointer-events-none"></div>

        {/* Close Button */}
        <button 
          onClick={() => setIsOpen(false)}
          className="absolute top-5 right-5 text-gray-500 hover:text-white transition-colors duration-200 z-50 flex items-center justify-center p-1.5 bg-white/5 hover:bg-white/10 rounded-full"
          aria-label="Kapat"
        >
          <span className="material-symbols-outlined text-[18px]">close</span>
        </button>

        {/* Tab Headers */}
        <div className="flex border-b border-white/5 mb-8 relative z-10">
          <button
            onClick={() => setActiveTab('login')}
            className={`flex-1 pb-4 text-sm font-bold tracking-wider uppercase transition-all duration-300 relative font-sans ${activeTab === 'login' ? 'text-blue-400' : 'text-gray-500 hover:text-gray-300'}`}
          >
            Giriş Yap
            {activeTab === 'login' && (
              <span className="absolute bottom-0 inset-x-0 h-[2px] bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.6)] animate-fade-in"></span>
            )}
          </button>
          <button
            onClick={() => setActiveTab('register')}
            className={`flex-1 pb-4 text-sm font-bold tracking-wider uppercase transition-all duration-300 relative font-sans ${activeTab === 'register' ? 'text-blue-400' : 'text-gray-500 hover:text-gray-300'}`}
          >
            Kayıt Ol
            {activeTab === 'register' && (
              <span className="absolute bottom-0 inset-x-0 h-[2px] bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.6)] animate-fade-in"></span>
            )}
          </button>
        </div>

        {/* Content Body */}
        <div className="relative z-10 font-sans">
          {activeTab === 'login' ? (
            /* ==================== LOGIN FORM ==================== */
            <form onSubmit={handleLoginSubmit} className="space-y-6">
              <div className="text-center mb-2">
                <h2 className="text-2xl font-black text-white tracking-tight" style={{ fontFamily: "'Outfit', sans-serif" }}>
                  Haber<span className="text-blue-500">Bot</span>'a Giriş Yapın
                </h2>
                <p className="text-gray-500 text-xs mt-1 font-sans">
                  Tartışmalara katılmak için oturum açın
                </p>
              </div>

              {loginError && (
                <div className="bg-red-500/10 border border-red-500/20 text-red-400 px-4 py-3 rounded-xl text-xs flex items-center space-x-2 animate-shake">
                  <span className="font-bold text-red-500">⚠️</span>
                  <span>{loginError}</span>
                </div>
              )}

              <div className="space-y-2">
                <label className="text-[10px] font-bold tracking-widest text-gray-400 uppercase" htmlFor="login-username">
                  Kullanıcı Adı
                </label>
                <div className="relative group/input">
                  <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 text-[18px] group-focus-within/input:text-blue-400 transition-colors duration-200">
                    person
                  </span>
                  <input
                    id="login-username"
                    type="text"
                    required
                    value={loginUsername}
                    onChange={(e) => setLoginUsername(e.target.value)}
                    placeholder="Kullanıcı adınız"
                    className="w-full bg-[#0a0a0a]/90 border border-blue-500/15 text-white placeholder-gray-600 pl-11 pr-4 py-3 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all duration-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.6)] focus:shadow-[0_0_15px_rgba(37,99,235,0.15)] text-sm font-sans"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <label className="text-[10px] font-bold tracking-widest text-gray-400 uppercase" htmlFor="login-password">
                    Şifre
                  </label>
                  <a href="#" onClick={(e) => e.preventDefault()} className="text-[10px] font-bold text-blue-400 hover:text-blue-300 transition-colors uppercase tracking-wider">
                    Şifremi Unuttum?
                  </a>
                </div>
                <div className="relative group/input">
                  <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 text-[18px] group-focus-within/input:text-blue-400 transition-colors duration-200">
                    lock
                  </span>
                  <input
                    id="login-password"
                    type={showLoginPassword ? "text" : "password"}
                    required
                    value={loginPassword}
                    onChange={(e) => setLoginPassword(e.target.value)}
                    placeholder="••••••••"
                    className="w-full bg-[#0a0a0a]/90 border border-blue-500/15 text-white placeholder-gray-600 pl-11 pr-12 py-3 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all duration-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.6)] focus:shadow-[0_0_15px_rgba(37,99,235,0.15)] text-sm font-sans"
                  />
                  <button
                    type="button"
                    onClick={() => setShowLoginPassword(!showLoginPassword)}
                    className="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-gray-500 hover:text-white cursor-pointer select-none text-[18px] transition-colors duration-200"
                  >
                    {showLoginPassword ? 'visibility_off' : 'visibility'}
                  </button>
                </div>
              </div>

              <button
                type="submit"
                disabled={loginLoading}
                className={`w-full bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-500 hover:to-blue-600 text-white font-sans py-3.5 rounded-xl font-bold uppercase tracking-wider text-xs hover:scale-[1.01] hover:shadow-[0_0_20px_rgba(37,99,235,0.35)] active:scale-[0.99] transition-all duration-300 flex items-center justify-center space-x-2 ${loginLoading ? 'opacity-70 cursor-not-allowed' : ''}`}
              >
                {loginLoading ? (
                  <>
                    <svg className="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Giriş Yapılıyor...</span>
                  </>
                ) : (
                  <span>Giriş Yap</span>
                )}
              </button>

              {/* Decorative Social Login */}
              <div className="relative flex items-center justify-center my-4">
                <div className="absolute inset-x-0 h-[1px] bg-white/5"></div>
                <span className="relative bg-[#050505] px-3 text-[10px] uppercase tracking-widest text-gray-600 font-bold">veya</span>
              </div>

              <button
                type="button"
                onClick={() => alert('HaberBot Google Auth yakında aktif olacaktır!')}
                className="w-full bg-white/5 border border-white/10 hover:bg-white/10 text-white font-sans py-3 rounded-xl font-bold text-xs transition-all duration-300 flex items-center justify-center space-x-3 active:scale-[0.99]"
              >
                <img src="https://docs.kodular.io/guides/component-examples/google-sign-in/google.png" alt="Google" className="h-4 w-4 object-contain brightness-95" />
                <span>Google ile Giriş Yap</span>
              </button>

              <div className="text-center pt-2">
                <p className="text-xs text-gray-500 font-sans">
                  Hesabınız yok mu?{' '}
                  <button
                    type="button"
                    onClick={() => setActiveTab('register')}
                    className="text-blue-400 hover:text-blue-300 font-bold transition-colors duration-200"
                  >
                    Kayıt Olun
                  </button>
                </p>
              </div>
            </form>
          ) : (
            /* ==================== REGISTER FORM ==================== */
            <form onSubmit={handleRegisterSubmit} className="space-y-5">
              <div className="text-center mb-2">
                <h2 className="text-2xl font-black text-white tracking-tight" style={{ fontFamily: "'Outfit', sans-serif" }}>
                  Aramıza Katılın
                </h2>
                <p className="text-gray-500 text-xs mt-1 font-sans">
                  Teknoloji haberlerini tartışmak için kayıt olun
                </p>
              </div>

              {registerError && (
                <div className="bg-red-500/10 border border-red-500/20 text-red-400 px-4 py-3 rounded-xl text-xs flex items-center space-x-2 animate-shake">
                  <span className="font-bold text-red-500">⚠️</span>
                  <span>{registerError}</span>
                </div>
              )}

              {registerSuccess && (
                <div className="bg-green-500/10 border border-green-500/20 text-green-400 px-4 py-3 rounded-xl text-xs flex items-center space-x-2 animate-pulse">
                  <span className="font-bold text-green-500">✓</span>
                  <span>Kayıt başarılı! Giriş alanına aktarılıyorsunuz...</span>
                </div>
              )}

              <div className="space-y-2">
                <label className="text-[10px] font-bold tracking-widest text-gray-400 uppercase" htmlFor="reg-username">
                  Kullanıcı Adı
                </label>
                <div className="relative group/input">
                  <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 text-[18px] group-focus-within/input:text-blue-400 transition-colors duration-200">
                    person
                  </span>
                  <input
                    id="reg-username"
                    type="text"
                    required
                    disabled={registerSuccess}
                    value={registerUsername}
                    onChange={(e) => setRegisterUsername(e.target.value)}
                    placeholder="Kullanıcı adınızı belirleyin"
                    className="w-full bg-[#0a0a0a]/90 border border-blue-500/15 text-white placeholder-gray-600 pl-11 pr-4 py-3 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all duration-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.6)] focus:shadow-[0_0_15px_rgba(37,99,235,0.15)] text-sm font-sans"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-[10px] font-bold tracking-widest text-gray-400 uppercase" htmlFor="reg-email">
                  E-Posta Adresi
                </label>
                <div className="relative group/input">
                  <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 text-[18px] group-focus-within/input:text-blue-400 transition-colors duration-200">
                    alternate_email
                  </span>
                  <input
                    id="reg-email"
                    type="email"
                    required
                    disabled={registerSuccess}
                    value={registerEmail}
                    onChange={(e) => setRegisterEmail(e.target.value)}
                    placeholder="e-posta@haberbot.com"
                    className="w-full bg-[#0a0a0a]/90 border border-blue-500/15 text-white placeholder-gray-600 pl-11 pr-4 py-3 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all duration-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.6)] focus:shadow-[0_0_15px_rgba(37,99,235,0.15)] text-sm font-sans"
                  />
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-[10px] font-bold tracking-widest text-gray-400 uppercase" htmlFor="reg-password">
                  Şifre
                </label>
                <div className="relative group/input">
                  <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 text-[18px] group-focus-within/input:text-blue-400 transition-colors duration-200">
                    lock
                  </span>
                  <input
                    id="reg-password"
                    type={showRegisterPassword ? "text" : "password"}
                    required
                    disabled={registerSuccess}
                    value={registerPassword}
                    onChange={(e) => setRegisterPassword(e.target.value)}
                    placeholder="En az 8 karakter"
                    className="w-full bg-[#0a0a0a]/90 border border-blue-500/15 text-white placeholder-gray-600 pl-11 pr-12 py-3 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all duration-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.6)] focus:shadow-[0_0_15px_rgba(37,99,235,0.15)] text-sm font-sans"
                  />
                  <button
                    type="button"
                    onClick={() => setShowRegisterPassword(!showRegisterPassword)}
                    className="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-gray-500 hover:text-white cursor-pointer select-none text-[18px] transition-colors duration-200"
                  >
                    {showRegisterPassword ? 'visibility_off' : 'visibility'}
                  </button>
                </div>

                {/* Password Strength Indicator Bar */}
                {registerPassword && (
                  <div className="space-y-1 pt-1.5 animate-fade-in font-sans">
                    <div className="flex justify-between items-center text-[10px]">
                      <span className="text-gray-500 font-bold uppercase tracking-wider">Şifre Gücü:</span>
                      <span className={`${strength.textClass} font-extrabold uppercase tracking-widest`}>{strength.label}</span>
                    </div>
                    <div className="h-1 w-full bg-white/5 rounded-full overflow-hidden">
                      <div className={`h-full ${strength.colorClass} ${strength.widthClass} transition-all duration-500 rounded-full`}></div>
                    </div>
                  </div>
                )}
              </div>

              <div className="space-y-2">
                <label className="text-[10px] font-bold tracking-widest text-gray-400 uppercase" htmlFor="reg-confirm">
                  Şifre Tekrarı
                </label>
                <div className="relative group/input">
                  <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-gray-500 text-[18px] group-focus-within/input:text-blue-400 transition-colors duration-200">
                    lock
                  </span>
                  <input
                    id="reg-confirm"
                    type={showConfirmPassword ? "text" : "password"}
                    required
                    disabled={registerSuccess}
                    value={registerConfirmPassword}
                    onChange={(e) => setRegisterConfirmPassword(e.target.value)}
                    placeholder="Şifrenizi tekrar girin"
                    className="w-full bg-[#0a0a0a]/90 border border-blue-500/15 text-white placeholder-gray-600 pl-11 pr-12 py-3 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all duration-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.6)] focus:shadow-[0_0_15px_rgba(37,99,235,0.15)] text-sm font-sans"
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="material-symbols-outlined absolute right-4 top-1/2 -translate-y-1/2 text-gray-500 hover:text-white cursor-pointer select-none text-[18px] transition-colors duration-200"
                  >
                    {showConfirmPassword ? 'visibility_off' : 'visibility'}
                  </button>
                </div>
              </div>

              <button
                type="submit"
                disabled={registerLoading || registerSuccess}
                className={`w-full bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-500 hover:to-blue-600 text-white font-sans py-3.5 rounded-xl font-bold uppercase tracking-wider text-xs hover:scale-[1.01] hover:shadow-[0_0_20px_rgba(37,99,235,0.35)] active:scale-[0.99] transition-all duration-300 flex items-center justify-center space-x-2 ${registerLoading || registerSuccess ? 'opacity-70 cursor-not-allowed' : ''}`}
              >
                {registerLoading ? (
                  <>
                    <svg className="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span>Kayıt Yapılıyor...</span>
                  </>
                ) : (
                  <span>Kayıt Ol</span>
                )}
              </button>

              <div className="text-center pt-2 border-t border-white/5 mt-4">
                <p className="text-xs text-gray-500 font-sans">
                  Zaten bir hesabınız var mı?{' '}
                  <button
                    type="button"
                    onClick={() => setActiveTab('login')}
                    className="text-blue-400 hover:text-blue-300 font-bold transition-colors duration-200"
                  >
                    Giriş Yapın
                  </button>
                </p>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}
