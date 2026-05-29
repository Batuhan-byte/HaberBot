import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { api } from '../../services/api';

export default function RegisterPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleRegister = async (e) => {
    e.preventDefault();
    setError('');
    
    if (password.length < 8) {
      setError('Şifre en az 8 karakter uzunluğunda olmalıdır.');
      return;
    }
    if (password !== confirmPassword) {
      setError('Şifreler uyuşmuyor.');
      return;
    }

    setLoading(true);
    try {
      await api.register(username, password);
      setSuccess(true);
      setTimeout(() => {
        navigate('/login');
      }, 2000);
    } catch (err) {
      setError(err.message || 'Kayıt başarısız. Lütfen başka bir kullanıcı adı deneyin.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative min-h-screen flex items-center justify-center bg-black overflow-hidden py-12 px-4 sm:px-6 lg:px-8">
      {/* Premium Ambient Background Glowing Blobs */}
      <div className="absolute top-1/4 left-1/4 -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-blue-600/10 rounded-full blur-[120px] pointer-events-none z-0"></div>
      <div className="absolute bottom-1/4 right-1/4 translate-x-1/2 translate-y-1/2 w-96 h-96 bg-indigo-600/8 rounded-full blur-[140px] pointer-events-none z-0"></div>

      <main className="w-full max-w-[440px] mx-auto z-10 relative">
        {/* Glassmorphic Premium Card Container */}
        <div className="bg-[#050505]/75 backdrop-blur-xl border border-blue-500/15 p-8 sm:p-10 rounded-2xl shadow-2xl relative overflow-hidden transition-all duration-300 hover:border-blue-500/25 hover:shadow-[0_20px_50px_rgba(0,0,0,0.8),_0_0_30px_rgba(37,99,235,0.03)]">
          {/* Subtle top border illumination */}
          <div className="absolute top-0 inset-x-0 h-[1px] bg-gradient-to-r from-transparent via-blue-500/30 to-transparent"></div>
          
          <div className="mb-8 text-center">
            {/* Syne/Outfit inspired premium brand title */}
            <h1 className="text-3xl bg-gradient-to-r from-blue-500 via-[#3b82f6] to-cyan-400 bg-clip-text text-transparent tracking-tight font-black font-sans" style={{ fontFamily: "'Outfit', sans-serif" }}>
              Kayıt Ol
            </h1>
            <p className="text-gray-400 text-sm mt-3 font-sans">
              HaberBot topluluğuna katılarak yorum yapmaya başlayın.
            </p>
          </div>

          {error && (
            <div className="bg-red-500/10 border border-red-500/20 text-red-400 px-4 py-3 rounded-xl text-sm mb-6 flex items-center space-x-2 animate-shake">
              <span className="font-bold text-red-500">⚠️</span>
              <span>{error}</span>
            </div>
          )}

          {success && (
            <div className="bg-green-500/10 border border-green-500/20 text-green-400 px-4 py-3 rounded-xl text-sm mb-6 flex items-center space-x-2 animate-pulse">
              <span className="font-bold text-green-500">✓</span>
              <span>Kayıt başarılı! Giriş sayfasına yönlendiriliyorsunuz...</span>
            </div>
          )}

          <form className="space-y-5" onSubmit={handleRegister}>
            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-gray-400 uppercase" htmlFor="username">
                Kullanıcı Adı
              </label>
              <input 
                className="w-full bg-[#0a0a0a]/90 border border-blue-500/15 text-white placeholder-gray-600 px-4 py-3 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all duration-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.6)] focus:shadow-[0_0_15px_rgba(37,99,235,0.15)] text-sm font-sans" 
                id="username" 
                name="username" 
                placeholder="Kullanıcı adınızı seçin" 
                required 
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                disabled={success}
              />
            </div>

            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-gray-400 uppercase" htmlFor="password">
                Şifre
              </label>
              <input 
                className="w-full bg-[#0a0a0a]/90 border border-blue-500/15 text-white placeholder-gray-600 px-4 py-3 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all duration-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.6)] focus:shadow-[0_0_15px_rgba(37,99,235,0.15)] text-sm font-sans" 
                id="password" 
                name="password" 
                placeholder="En az 8 karakter" 
                required 
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={success}
              />
            </div>

            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-gray-400 uppercase" htmlFor="confirmPassword">
                Şifre Tekrarı
              </label>
              <input 
                className="w-full bg-[#0a0a0a]/90 border border-blue-500/15 text-white placeholder-gray-600 px-4 py-3 rounded-xl focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500/30 transition-all duration-300 shadow-[inset_0_2px_4px_rgba(0,0,0,0.6)] focus:shadow-[0_0_15px_rgba(37,99,235,0.15)] text-sm font-sans" 
                id="confirmPassword" 
                name="confirmPassword" 
                placeholder="Şifrenizi tekrar girin" 
                required 
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                disabled={success}
              />
            </div>

            <button 
              className={`w-full bg-gradient-to-r from-blue-600 to-blue-500 hover:from-blue-500 hover:to-blue-600 text-white font-sans py-3.5 rounded-xl font-bold uppercase tracking-wider text-xs hover:scale-[1.01] hover:shadow-[0_0_20px_rgba(37,99,235,0.35)] active:scale-[0.99] transition-all duration-300 flex items-center justify-center space-x-2 ${loading || success ? 'opacity-70 cursor-not-allowed' : ''}`}
              type="submit"
              disabled={loading || success}
            >
              {loading ? (
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
          </form>

          <div className="mt-8 text-center border-t border-white/5 pt-6">
            <p className="text-xs text-gray-400 font-sans">
              Zaten hesabınız var mı?{' '}
              <Link to="/login" className="text-blue-400 hover:text-blue-300 font-bold transition-colors duration-200">
                Giriş Yapın
              </Link>
            </p>
          </div>
        </div>

        <footer className="mt-8 text-center">
          <p className="text-[10px] uppercase tracking-widest text-gray-600">
            © 2026 HaberBot AI. Tüm hakları saklıdır.
          </p>
        </footer>
      </main>

      <div className="fixed inset-0 pointer-events-none opacity-[0.02] bg-[url('https://www.transparenttextures.com/patterns/stardust.png')] z-50"></div>
    </div>
  );
}
