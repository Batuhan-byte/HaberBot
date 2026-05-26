import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function LoginPage() {
  const [apiKey, setApiKey] = useState('');
  const navigate = useNavigate();

  const handleLogin = (e) => {
    e.preventDefault();
    if (apiKey) {
      localStorage.setItem('admin_api_key', apiKey);
      navigate('/admin');
    }
  };

  return (
    <>
      <main className="w-full max-w-[400px] px-margin-mobile md:px-0 mx-auto mt-24">
        <div className="mb-12 text-center">
          <h1 className="font-headline-lg text-headline-lg text-primary tracking-tighter">
            HaberBot Admin
          </h1>
        </div>

        <section className="space-y-6">
          <form className="space-y-4" onSubmit={handleLogin}>
            <div className="relative group">
              <label className="sr-only" htmlFor="password">API Key</label>
              <input 
                className="w-full bg-background border border-outline-variant text-on-surface px-4 py-3 rounded focus:outline-none focus:border-primary transition-colors duration-150 font-label-md text-label-md" 
                id="password" 
                name="password" 
                placeholder="Admin API Key" 
                required 
                type="password"
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
              />
            </div>
            <button className="w-full bg-primary text-on-primary font-body-md text-body-md py-3 rounded-lg font-bold hover:opacity-90 active:scale-[0.98] transition-all duration-150" type="submit">
              Giriş Yap
            </button>
          </form>
        </section>

        <footer className="mt-24 text-center">
          <p className="font-label-sm text-label-sm text-on-surface-variant opacity-40">
            © 2024 HaberBot AI. All rights reserved.
          </p>
        </footer>
      </main>

      <div className="fixed inset-0 pointer-events-none opacity-[0.03] bg-[url('https://www.transparenttextures.com/patterns/stardust.png')] z-50"></div>
    </>
  );
}
