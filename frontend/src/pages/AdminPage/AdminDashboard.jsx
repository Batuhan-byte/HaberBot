import React, { useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';

export default function AdminDashboard() {
  const navigate = useNavigate();

  useEffect(() => {
    const key = localStorage.getItem('admin_api_key');
    if (!key) {
      navigate('/login');
    }
  }, [navigate]);

  return (
    <>
      <aside className="h-screen w-64 fixed left-0 top-0 bg-surface-container-lowest border-r border-outline-variant flex flex-col py-6 px-4 gap-4 z-50">
<div className="mb-8 px-2">
<h1 className="font-headline-md text-headline-md text-primary">Admin Console</h1>
<p className="font-label-sm text-label-sm text-on-surface-variant opacity-60">v1.0.4</p>
</div>
<nav className="flex-1 space-y-1">

<Link className="flex items-center gap-3 px-3 py-2 bg-surface-container-high text-primary border-r-2 border-primary cursor-pointer transition-all duration-150 font-label-md text-label-md" to="/admin">
<span className="material-symbols-outlined" data-icon="dashboard">dashboard</span>
<span>Dashboard</span>
</Link>
<Link className="flex items-center gap-3 px-3 py-2 text-on-surface-variant hover:bg-surface-container hover:text-primary cursor-pointer transition-all duration-150 font-label-md text-label-md" to="/admin/sources">
<span className="material-symbols-outlined" data-icon="rss_feed">rss_feed</span>
<span>RSS Sources</span>
</Link>
<Link className="flex items-center gap-3 px-3 py-2 text-on-surface-variant hover:bg-surface-container hover:text-primary cursor-pointer transition-all duration-150 font-label-md text-label-md" to="/admin/content">
<span className="material-symbols-outlined" data-icon="article">article</span>
<span>Content</span>
</Link>
<Link className="flex items-center gap-3 px-3 py-2 text-on-surface-variant hover:bg-surface-container hover:text-primary cursor-pointer transition-all duration-150 font-label-md text-label-md" to="/admin">
<span className="material-symbols-outlined" data-icon="group">group</span>
<span>Users</span>
</Link>
<Link className="flex items-center gap-3 px-3 py-2 text-on-surface-variant hover:bg-surface-container hover:text-primary cursor-pointer transition-all duration-150 font-label-md text-label-md" to="/admin">
<span className="material-symbols-outlined" data-icon="terminal">terminal</span>
<span>Logs</span>
</Link>
</nav>
<div className="mt-auto space-y-1">
<Link to="/admin/sources" className="w-full mb-4 bg-primary text-background font-label-md text-label-md py-3 rounded hover:opacity-90 active:scale-[0.98] transition-transform flex items-center justify-center">
                Add New Source
            </Link>
<Link className="flex items-center gap-3 px-3 py-2 text-on-surface-variant hover:bg-surface-container hover:text-primary transition-all font-label-md text-label-md" to="/admin">
<span className="material-symbols-outlined" data-icon="help">help</span>
<span>Support</span>
</Link>
<Link className="flex items-center gap-3 px-3 py-2 text-on-surface-variant hover:bg-surface-container hover:text-primary transition-all font-label-md text-label-md" to="/admin">
<span className="material-symbols-outlined" data-icon="settings">settings</span>
<span>Settings</span>
</Link>
</div>
</aside>

<main className="ml-64 min-h-screen p-margin-desktop bg-background overflow-x-hidden">

<header className="flex justify-between items-center mb-8">
<div>
<h2 className="font-headline-lg text-headline-lg text-primary tracking-tight">HaberBot Overview</h2>
<p className="font-body-md text-body-md text-on-surface-variant">Real-time content aggregation and processing engine metrics.</p>
</div>
<div className="flex items-center gap-4">
<button className="p-2 border border-outline-variant hover:border-primary transition-colors">
<span className="material-symbols-outlined text-on-surface-variant" data-icon="notifications">notifications</span>
</button>
<div className="h-10 w-10 bg-surface-container-high rounded-full overflow-hidden border border-outline-variant">
<img alt="User Profile" className="w-full h-full object-cover" data-alt="A clean, professional headshot of a technical administrator for a high-performance software company. The subject is a person with a focused and intelligent expression. The lighting is sophisticated and directional, creating a modern professional aesthetic. The background is a minimalist, dark tech-office environment with subtle blue and grey highlights. The overall mood is one of reliability and expertise in a high-tech monochrome environment." src="https://lh3.googleusercontent.com/aida-public/AB6AXuDgNLlMcaykhEC2Ra9c252-Y5XfxrfGQXElo2-S0NCmUxdLCxuN2Fw_fxnpyLPcDvhRHAn6AB3_aslMpHyNMm5FL-iPqWnnMm59sJISEQgY7WFNfBNnaSPBoX9Dtn2vzPELzHjFlXNE1G4dYoqq8INFGUFvKEeF9xBz_4AAjYeZ-bIw0is1xBVr4GAJVbb8W-HvW_ZIn2cHqSUxAYZB9nuR_R7W59e9xmCwD0FbHJ-RJUE6auVqsdUWfgExQWuypXy1XYtzds1GitY"/>
</div>
</div>
</header>

<section className="grid grid-cols-1 md:grid-cols-3 gap-gutter mb-8">

<div className="bg-surface-container-lowest border border-outline-variant p-6 rounded transition-all hover:border-outline">
<div className="flex justify-between items-start mb-4">
<p className="font-label-sm text-label-sm text-on-surface-variant">BUGÜN ÇEKİLEN HABERLER</p>
<span className="material-symbols-outlined text-primary text-[20px]" data-icon="auto_awesome">auto_awesome</span>
</div>
<div className="flex items-baseline gap-2">
<span className="font-display-lg text-display-lg text-primary">42</span>
<span className="font-label-sm text-label-sm text-green-500">+12%</span>
</div>
<div className="mt-4 h-[2px] w-full bg-surface-container-high">
<div className="h-full bg-primary w-[65%]"></div>
</div>
</div>

<div className="bg-surface-container-lowest border border-outline-variant p-6 rounded transition-all hover:border-outline">
<div className="flex justify-between items-start mb-4">
<p className="font-label-sm text-label-sm text-on-surface-variant">AKTİF KAYNAKLAR</p>
<span className="material-symbols-outlined text-primary text-[20px]" data-icon="hub">hub</span>
</div>
<div className="flex items-baseline gap-2">
<span className="font-display-lg text-display-lg text-primary">12</span>
<span className="font-label-sm text-label-sm text-on-surface-variant">RSS/API</span>
</div>
<div className="mt-4 flex gap-1">
<div className="h-1 w-4 bg-primary"></div>
<div className="h-1 w-4 bg-primary"></div>
<div className="h-1 w-4 bg-primary"></div>
<div className="h-1 w-4 bg-outline-variant"></div>
<div className="h-1 w-4 bg-outline-variant"></div>
</div>
</div>

<div className="bg-surface-container-lowest border border-outline-variant p-6 rounded transition-all hover:border-outline">
<div className="flex justify-between items-start mb-4">
<p className="font-label-sm text-label-sm text-on-surface-variant">GEMINI API KOTASI</p>
<span className="material-symbols-outlined text-primary text-[20px]" data-icon="memory">memory</span>
</div>
<div className="flex items-baseline gap-2">
<span className="font-display-lg text-display-lg text-primary">%24</span>
<span className="font-label-sm text-label-sm text-on-surface-variant">KULLANILDI</span>
</div>
<p className="mt-4 font-label-sm text-label-sm text-on-surface-variant opacity-60">Reset in 4 days</p>
</div>
</section>

<div className="grid grid-cols-1 lg:grid-cols-3 gap-gutter">

<div className="lg:col-span-2 bg-surface-container-lowest border border-outline-variant rounded p-6">
<div className="flex justify-between items-center mb-8">
<h3 className="font-headline-md text-headline-md text-primary">Aktivite Grafiği</h3>
<div className="flex gap-2">
<button className="px-3 py-1 bg-surface-container-high text-primary border border-outline-variant font-label-sm text-label-sm">24H</button>
<button className="px-3 py-1 text-on-surface-variant border border-transparent hover:border-outline-variant font-label-sm text-label-sm transition-all">7D</button>
</div>
</div>
<div className="relative h-[300px] w-full flex items-end">

<div className="absolute inset-0 grid grid-rows-4 w-full">
<div className="border-b border-surface-container-high w-full"></div>
<div className="border-b border-surface-container-high w-full"></div>
<div className="border-b border-surface-container-high w-full"></div>
<div className="border-b border-surface-container-high w-full"></div>
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
<div className="absolute bottom-[-24px] left-0 right-0 flex justify-between font-label-sm text-label-sm text-on-surface-variant">
<span>00:00</span>
<span>06:00</span>
<span>12:00</span>
<span>18:00</span>
<span>23:59</span>
</div>
</div>
</div>

<div className="bg-surface-container-lowest border border-outline-variant rounded p-6">
<h3 className="font-headline-md text-headline-md text-primary mb-6">Recent Events</h3>
<div className="space-y-6">
<div className="flex gap-4">
<div className="h-2 w-2 rounded-full bg-green-500 mt-2 shrink-0"></div>
<div>
<p className="font-body-md text-body-md text-primary">New York Times Sync</p>
<p className="font-label-sm text-label-sm text-on-surface-variant opacity-60">Completed • 2m ago</p>
</div>
</div>
<div className="flex gap-4">
<div className="h-2 w-2 rounded-full bg-primary mt-2 shrink-0"></div>
<div>
<p className="font-body-md text-body-md text-primary">Gemini Analysis</p>
<p className="font-label-sm text-label-sm text-on-surface-variant opacity-60">Running analysis on 12 articles</p>
</div>
</div>
<div className="flex gap-4">
<div className="h-2 w-2 rounded-full bg-error mt-2 shrink-0"></div>
<div>
<p className="font-body-md text-body-md text-primary">API Connection Error</p>
<p className="font-label-sm text-label-sm text-on-surface-variant opacity-60">Source: TechCrunch • 15m ago</p>
</div>
</div>
<div className="flex gap-4">
<div className="h-2 w-2 rounded-full bg-surface-container-high mt-2 shrink-0"></div>
<div>
<p className="font-body-md text-body-md text-primary">Cron Job: DB Clean</p>
<p className="font-label-sm text-label-sm text-on-surface-variant opacity-60">Scheduled for 03:00 AM</p>
</div>
</div>
</div>
<button className="w-full mt-8 border border-outline-variant py-2 font-label-sm text-label-sm hover:border-primary transition-colors">
                    View Full System Log
                </button>
</div>
</div>

<footer className="w-full py-12 mt-12 border-t border-outline-variant flex justify-between items-center">
<span className="font-label-sm text-label-sm text-on-surface-variant">© 2024 HaberBot AI. All rights reserved.</span>
<div className="flex gap-6">
<a className="font-label-sm text-label-sm text-on-surface-variant hover:text-primary underline transition-all" href="#">Terms</a>
<a className="font-label-sm text-label-sm text-on-surface-variant hover:text-primary underline transition-all" href="#">Privacy</a>
<a className="font-label-sm text-label-sm text-on-surface-variant hover:text-primary underline transition-all" href="#">API Docs</a>
<a className="font-label-sm text-label-sm text-on-surface-variant hover:text-primary underline transition-all" href="#">Status</a>
</div>
</footer>
</main>
    </>
  );
}
