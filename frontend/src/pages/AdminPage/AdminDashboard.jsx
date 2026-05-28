import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AdminLayout from '../../components/Admin/AdminLayout';

export default function AdminDashboard() {
  const navigate = useNavigate();

  useEffect(() => {
    const key = localStorage.getItem('admin_api_key');
    if (!key) {
      navigate('/login');
    }
  }, [navigate]);

  const headerActions = (
    <>
      <button className="p-2 border border-zinc-800 hover:border-blue-500 transition-colors rounded">
        <span className="material-symbols-outlined text-zinc-400" data-icon="notifications">notifications</span>
      </button>
      <div className="h-10 w-10 bg-zinc-900 rounded-full overflow-hidden border border-zinc-800">
        <img alt="User Profile" className="w-full h-full object-cover" data-alt="A clean, professional headshot of a technical administrator for a high-performance software company. The subject is a person with a focused and intelligent expression. The lighting is sophisticated and directional, creating a modern professional aesthetic. The background is a minimalist, dark tech-office environment with subtle blue and grey highlights. The overall mood is one of reliability and expertise in a high-tech monochrome environment." src="https://lh3.googleusercontent.com/aida-public/AB6AXuDgNLlMcaykhEC2Ra9c252-Y5XfxrfGQXElo2-S0NCmUxdLCxuN2Fw_fxnpyLPcDvhRHAn6AB3_aslMpHyNMm5FL-iPqWnnMm59sJISEQgY7WFNfBNnaSPBoX9Dtn2vzPELzHjFlXNE1G4dYoqq8INFGUFvKEeF9xBz_4AAjYeZ-bIw0is1xBVr4GAJVbb8W-HvW_ZIn2cHqSUxAYZB9nuR_R7W59e9xmCwD0FbHJ-RJUE6auVqsdUWfgExQWuypXy1XYtzds1GitY"/>
      </div>
    </>
  );

  return (
    <AdminLayout
      title="HaberBot Overview"
      subtitle="Real-time content aggregation and processing engine metrics."
      actions={headerActions}
    >

      <section className="grid grid-cols-1 md:grid-cols-3 gap-gutter mb-8">
        <div className="bg-zinc-950 border border-zinc-800 p-6 rounded-xl transition-all hover:border-zinc-700">
          <div className="flex justify-between items-start mb-4">
            <p className="text-xs font-bold uppercase tracking-wider text-zinc-500">BUGÜN ÇEKİLEN HABERLER</p>
            <span className="material-symbols-outlined text-blue-400 text-[20px]" data-icon="auto_awesome">auto_awesome</span>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-extrabold text-white">42</span>
            <span className="text-xs font-bold text-emerald-400">+12%</span>
          </div>
          <div className="mt-4 h-[2px] w-full bg-zinc-900">
            <div className="h-full bg-blue-500 w-[65%]"></div>
          </div>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 p-6 rounded-xl transition-all hover:border-zinc-700">
          <div className="flex justify-between items-start mb-4">
            <p className="text-xs font-bold uppercase tracking-wider text-zinc-500">AKTİF KAYNAKLAR</p>
            <span className="material-symbols-outlined text-blue-400 text-[20px]" data-icon="hub">hub</span>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-extrabold text-white">12</span>
            <span className="text-xs font-bold text-zinc-400">RSS/API</span>
          </div>
          <div className="mt-4 flex gap-1">
            <div className="h-1 w-4 bg-blue-500"></div>
            <div className="h-1 w-4 bg-blue-500"></div>
            <div className="h-1 w-4 bg-blue-500"></div>
            <div className="h-1 w-4 bg-zinc-800"></div>
            <div className="h-1 w-4 bg-zinc-800"></div>
          </div>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 p-6 rounded-xl transition-all hover:border-zinc-700">
          <div className="flex justify-between items-start mb-4">
            <p className="text-xs font-bold uppercase tracking-wider text-zinc-500">GEMINI API KOTASI</p>
            <span className="material-symbols-outlined text-blue-400 text-[20px]" data-icon="memory">memory</span>
          </div>
          <div className="flex items-baseline gap-2">
            <span className="text-3xl font-extrabold text-white">%24</span>
            <span className="text-xs font-bold text-zinc-400">KULLANILDI</span>
          </div>
          <p className="mt-4 text-xs font-medium text-zinc-500">Reset in 4 days</p>
        </div>
      </section>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-gutter">
        <div className="lg:col-span-2 bg-zinc-950 border border-zinc-800 rounded-xl p-6">
          <div className="flex justify-between items-center mb-8">
            <h3 className="text-lg font-bold text-white">Aktivite Grafiği</h3>
            <div className="flex gap-2">
              <button className="px-3 py-1 bg-zinc-900 text-blue-400 border border-zinc-800 text-xs font-bold">24H</button>
              <button className="px-3 py-1 text-zinc-400 border border-transparent hover:border-zinc-800 text-xs font-bold transition-all">7D</button>
            </div>
          </div>
          <div className="relative h-[300px] w-full flex items-end">
            <div className="absolute inset-0 grid grid-rows-4 w-full">
              <div className="border-b border-zinc-900 w-full"></div>
              <div className="border-b border-zinc-900 w-full"></div>
              <div className="border-b border-zinc-900 w-full"></div>
              <div className="border-b border-zinc-900 w-full"></div>
            </div>

            <svg className="w-full h-full relative z-10 overflow-visible" viewbox="0 0 800 300">
              <path className="chart-path" d="M0 250 Q 100 200 200 220 T 400 120 T 600 180 T 800 50" fill="none" stroke="white" strokeWidth="2"></path>
              <path d="M0 250 Q 100 200 200 220 T 400 120 T 600 180 T 800 50 L 800 300 L 0 300 Z" fill="url(#gradient)" opacity="0.1"></path>
              <defs>
                <lineargradient id="gradient" x1="0%" x2="0%" y1="0%" y2="100%">
                  <stop offset="0%" style={{}}></stop>
                  <stop offset="100%" style={{}}></stop>
                </lineargradient>
              </defs>

              <circle cx="200" cy="220" fill="white" r="4"></circle>
              <circle cx="400" cy="120" fill="white" r="4"></circle>
              <circle cx="600" cy="180" fill="white" r="4"></circle>
              <circle cx="800" cy="50" fill="white" r="4"></circle>
            </svg>
            <div className="absolute bottom-[-24px] left-0 right-0 flex justify-between text-xs font-bold text-zinc-500">
              <span>00:00</span>
              <span>06:00</span>
              <span>12:00</span>
              <span>18:00</span>
              <span>23:59</span>
            </div>
          </div>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 rounded-xl p-6">
          <h3 className="text-lg font-bold text-white mb-6">Recent Events</h3>
          <div className="space-y-6">
            <div className="flex gap-4">
              <div className="h-2 w-2 rounded-full bg-green-500 mt-2 shrink-0"></div>
              <div>
                <p className="text-sm font-semibold text-white">New York Times Sync</p>
                <p className="text-xs text-zinc-500">Completed • 2m ago</p>
              </div>
            </div>
            <div className="flex gap-4">
              <div className="h-2 w-2 rounded-full bg-blue-500 mt-2 shrink-0"></div>
              <div>
                <p className="text-sm font-semibold text-white">Gemini Analysis</p>
                <p className="text-xs text-zinc-500">Running analysis on 12 articles</p>
              </div>
            </div>
            <div className="flex gap-4">
              <div className="h-2 w-2 rounded-full bg-red-500 mt-2 shrink-0"></div>
              <div>
                <p className="text-sm font-semibold text-white">API Connection Error</p>
                <p className="text-xs text-zinc-500">Source: TechCrunch • 15m ago</p>
              </div>
            </div>
            <div className="flex gap-4">
              <div className="h-2 w-2 rounded-full bg-zinc-700 mt-2 shrink-0"></div>
              <div>
                <p className="text-sm font-semibold text-white">Cron Job: DB Clean</p>
                <p className="text-xs text-zinc-500">Scheduled for 03:00 AM</p>
              </div>
            </div>
          </div>
          <button className="w-full mt-8 border border-zinc-800 py-2 text-xs font-bold text-zinc-400 hover:border-zinc-700 hover:text-white transition-colors">
            View Full System Log
          </button>
        </div>
      </div>

      <footer className="w-full py-12 mt-12 border-t border-zinc-800 flex justify-between items-center">
        <span className="text-xs font-bold text-zinc-500">© 2024 HaberBot AI. All rights reserved.</span>
        <div className="flex gap-6">
          <a className="text-xs font-bold text-zinc-500 hover:text-white underline transition-all" href="#">Terms</a>
          <a className="text-xs font-bold text-zinc-500 hover:text-white underline transition-all" href="#">Privacy</a>
          <a className="text-xs font-bold text-zinc-500 hover:text-white underline transition-all" href="#">API Docs</a>
          <a className="text-xs font-bold text-zinc-500 hover:text-white underline transition-all" href="#">Status</a>
        </div>
      </footer>
    </AdminLayout>
  );
}
