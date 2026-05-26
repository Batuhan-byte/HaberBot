import React from 'react';

export default function TopBar() {
  return (
    <header className="fixed top-0 right-0 w-[calc(100%-16rem)] z-40 bg-surface-dim/20 backdrop-blur-md border-b border-outline-variant/10 flex justify-between items-center px-container-padding h-16">
      <div className="flex items-center gap-4 flex-grow max-w-xl">
        <div className="relative w-full">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[20px]">search</span>
          <input 
            className="w-full bg-surface-container-lowest/50 border border-outline-variant/30 rounded-full py-1.5 pl-10 pr-4 text-label-md focus:outline-none focus:border-primary/50 transition-colors text-white" 
            placeholder="Search topics, logs, or bots..." 
            type="text"
          />
        </div>
      </div>
      
      <div className="flex items-center gap-3">
        <button className="p-2 hover:bg-white/10 rounded-full transition-colors scale-95 active:scale-90 text-on-surface-variant">
          <span className="material-symbols-outlined">notifications</span>
        </button>
        <button className="p-2 hover:bg-white/10 rounded-full transition-colors scale-95 active:scale-90 text-on-surface-variant">
          <span className="material-symbols-outlined">account_circle</span>
        </button>
      </div>
    </header>
  );
}
