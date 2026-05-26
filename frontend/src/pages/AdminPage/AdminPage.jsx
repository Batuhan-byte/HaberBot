import React, { useState, useEffect } from 'react';
import Sidebar from '../../components/Admin/Sidebar';
import TopBar from '../../components/Admin/TopBar';

export default function AdminPage() {
  const [progress, setProgress] = useState(68);

  useEffect(() => {
    const interval = setInterval(() => {
      setProgress(prev => {
        let next = prev + Math.random() * 0.5;
        if (next > 100) next = 0;
        return next;
      });
    }, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="font-body-md text-body-md antialiased min-h-screen">
      <Sidebar />
      <TopBar />
      
      {/* Main Canvas */}
      <main className="ml-64 pt-24 px-container-padding pb-container-padding">
        
        {/* Header & Action Section */}
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-8">
          <div>
            <h1 className="font-display-lg text-display-lg text-on-surface mb-2 tracking-tight">Topics Management</h1>
            <p className="text-body-lg text-on-surface-variant max-w-2xl">Configure and monitor the core interest areas of your AI news aggregation engine.</p>
          </div>
          <button className="pill-shaped px-8 py-3 font-label-md text-label-md font-bold text-white gradient-button rounded-full flex items-center gap-2">
            <span className="material-symbols-outlined text-[18px]">add</span>
            Add New Topic
          </button>
        </div>

        {/* System Controls Section */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-section-gap">
          <div className="md:col-span-2 glass-card p-glass-padding">
            <div className="flex items-center justify-between mb-6">
              <h2 className="font-headline-md text-headline-md text-on-surface">System Controls</h2>
              <div className="flex items-center gap-2 px-3 py-1 bg-surface-container-highest/50 rounded-full">
                <span className="status-pulse"></span>
                <span className="text-label-sm text-tertiary">System Live</span>
              </div>
            </div>
            
            <div className="flex flex-wrap gap-4">
              <button className="gradient-button px-10 py-4 rounded-full font-label-md text-label-md font-bold text-white flex items-center gap-3 shadow-lg">
                <span className="material-symbols-outlined">rss_feed</span>
                Fetch News
              </button>
              <button className="gradient-button px-10 py-4 rounded-full font-label-md text-label-md font-bold text-white flex items-center gap-3 shadow-lg">
                <span className="material-symbols-outlined">psychology</span>
                Process AI
              </button>
            </div>
            
            <div className="mt-8">
              <div className="flex justify-between items-end mb-2">
                <span className="text-label-sm text-on-surface-variant">AI Processing Progress</span>
                <span className="text-label-sm text-primary">{Math.floor(progress)}%</span>
              </div>
              <div className="w-full h-2 bg-surface-container-highest rounded-full overflow-hidden relative">
                <div className="h-full bg-primary rounded-full relative" style={{ width: `${progress}%` }}>
                  <div className="absolute inset-0 progress-shimmer"></div>
                </div>
              </div>
              <div className="flex justify-between mt-3 text-[12px] text-on-surface-variant/60">
                <span>Analysis: Large Language Models Trends</span>
                <span>Last run: 5 minutes ago</span>
              </div>
            </div>
          </div>

          {/* Quick Stats / Status Card */}
          <div className="glass-card p-glass-padding flex flex-col justify-between">
            <div>
              <h3 className="text-label-md font-bold text-on-surface-variant mb-4 uppercase tracking-widest">Health Monitor</h3>
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-body-md text-on-surface-variant">Active Topics</span>
                  <span className="font-headline-md text-headline-md text-on-surface">12</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-body-md text-on-surface-variant">Daily Articles</span>
                  <span className="font-headline-md text-headline-md text-on-surface">1.2k</span>
                </div>
              </div>
            </div>
            <div className="mt-6 pt-6 border-t border-outline-variant/10">
              <div className="flex items-center gap-2 text-on-surface-variant">
                <span className="material-symbols-outlined text-[18px]">history</span>
                <span className="text-label-sm">Next scheduled fetch in 12m</span>
              </div>
            </div>
          </div>
        </div>

        {/* Topics Table Section */}
        <div className="glass-card overflow-hidden">
          <div className="px-glass-padding py-6 border-b border-outline-variant/10 flex justify-between items-center">
            <h2 className="font-headline-md text-headline-md text-on-surface">Active Topics</h2>
            <div className="flex gap-2">
              <button className="p-2 hover:bg-white/5 rounded-lg border border-outline-variant/10 text-on-surface-variant">
                <span className="material-symbols-outlined">filter_list</span>
              </button>
              <button className="p-2 hover:bg-white/5 rounded-lg border border-outline-variant/10 text-on-surface-variant">
                <span className="material-symbols-outlined">more_vert</span>
              </button>
            </div>
          </div>
          
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="text-label-sm text-on-surface-variant/70 uppercase tracking-widest border-b border-outline-variant/10 bg-white/2">
                  <th className="px-glass-padding py-4 font-bold">Title</th>
                  <th className="px-glass-padding py-4 font-bold">Slug</th>
                  <th className="px-glass-padding py-4 font-bold">Source</th>
                  <th className="px-glass-padding py-4 font-bold">Status</th>
                  <th className="px-glass-padding py-4 font-bold text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/5">
                
                {/* Row 1: AI */}
                <tr className="hover:bg-white/2 transition-colors group">
                  <td className="px-glass-padding py-5">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-xl bg-primary-container/20 flex items-center justify-center text-primary group-hover:scale-110 transition-transform">
                        <span className="material-symbols-outlined" style={{ fontVariationSettings: "'FILL' 1" }}>psychology</span>
                      </div>
                      <span className="text-body-md font-bold text-on-surface">Artificial Intelligence</span>
                    </div>
                  </td>
                  <td className="px-glass-padding py-5">
                    <code className="text-label-sm bg-surface-container-highest px-2 py-1 rounded text-tertiary">ai-news</code>
                  </td>
                  <td className="px-glass-padding py-5 text-on-surface-variant">Multi-Source API</td>
                  <td className="px-glass-padding py-5">
                    <span className="flex items-center gap-2 text-label-sm text-tertiary">
                      <span className="w-1.5 h-1.5 rounded-full bg-tertiary"></span>
                      Active
                    </span>
                  </td>
                  <td className="px-glass-padding py-5 text-right">
                    <button className="p-2 hover:text-primary transition-colors">
                      <span className="material-symbols-outlined">edit</span>
                    </button>
                    <button className="p-2 hover:text-error transition-colors">
                      <span className="material-symbols-outlined">delete</span>
                    </button>
                  </td>
                </tr>

                {/* Row 2: Technology */}
                <tr className="hover:bg-white/2 transition-colors group">
                  <td className="px-glass-padding py-5">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-xl bg-secondary-container/20 flex items-center justify-center text-secondary group-hover:scale-110 transition-transform">
                        <span className="material-symbols-outlined" style={{ fontVariationSettings: "'FILL' 1" }}>devices</span>
                      </div>
                      <span className="text-body-md font-bold text-on-surface">Technology</span>
                    </div>
                  </td>
                  <td className="px-glass-padding py-5">
                    <code className="text-label-sm bg-surface-container-highest px-2 py-1 rounded text-tertiary">tech-trends</code>
                  </td>
                  <td className="px-glass-padding py-5 text-on-surface-variant">Global RSS Hub</td>
                  <td className="px-glass-padding py-5">
                    <span className="flex items-center gap-2 text-label-sm text-tertiary">
                      <span className="w-1.5 h-1.5 rounded-full bg-tertiary"></span>
                      Active
                    </span>
                  </td>
                  <td className="px-glass-padding py-5 text-right">
                    <button className="p-2 hover:text-primary transition-colors">
                      <span className="material-symbols-outlined">edit</span>
                    </button>
                    <button className="p-2 hover:text-error transition-colors">
                      <span className="material-symbols-outlined">delete</span>
                    </button>
                  </td>
                </tr>

                {/* Row 3: Cybersecurity */}
                <tr className="hover:bg-white/2 transition-colors group">
                  <td className="px-glass-padding py-5">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-xl bg-tertiary-container/10 flex items-center justify-center text-tertiary group-hover:scale-110 transition-transform">
                        <span className="material-symbols-outlined" style={{ fontVariationSettings: "'FILL' 1" }}>security</span>
                      </div>
                      <span className="text-body-md font-bold text-on-surface">Cybersecurity</span>
                    </div>
                  </td>
                  <td className="px-glass-padding py-5">
                    <code className="text-label-sm bg-surface-container-highest px-2 py-1 rounded text-tertiary">cyber-sec</code>
                  </td>
                  <td className="px-glass-padding py-5 text-on-surface-variant">Threat Intelligence Feed</td>
                  <td className="px-glass-padding py-5">
                    <span className="flex items-center gap-2 text-label-sm text-on-surface-variant/40">
                      <span className="w-1.5 h-1.5 rounded-full bg-on-surface-variant/40"></span>
                      Paused
                    </span>
                  </td>
                  <td className="px-glass-padding py-5 text-right">
                    <button className="p-2 hover:text-primary transition-colors">
                      <span className="material-symbols-outlined">edit</span>
                    </button>
                    <button className="p-2 hover:text-error transition-colors">
                      <span className="material-symbols-outlined">delete</span>
                    </button>
                  </td>
                </tr>
                
              </tbody>
            </table>
          </div>
          
          <div className="px-glass-padding py-4 border-t border-outline-variant/10 flex items-center justify-between">
            <span className="text-label-sm text-on-surface-variant">Showing 3 of 12 topics</span>
            <div className="flex gap-2">
              <button className="px-4 py-2 bg-surface-container-highest/30 rounded-lg text-label-sm text-on-surface-variant hover:bg-white/10 transition-colors">Previous</button>
              <button className="px-4 py-2 bg-surface-container-highest/30 rounded-lg text-label-sm text-on-surface-variant hover:bg-white/10 transition-colors">Next</button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
