import React from 'react';
import { Link } from 'react-router-dom';

export default function Sidebar() {
  return (
    <aside className="h-screen w-64 fixed left-0 top-0 bg-surface-container-low/30 backdrop-blur-xl border-r border-outline-variant/10 z-50 flex flex-col p-glass-padding transition-all duration-300 ease-in-out">
      <div className="mb-10 flex flex-col">
        <span className="font-headline-md text-headline-md font-bold text-primary">HaberBot Admin</span>
        <span className="text-label-sm text-on-surface-variant/60 uppercase tracking-widest mt-1">AI Command Center</span>
      </div>
      
      <nav className="flex flex-col gap-2 flex-grow">
        {/* Active: Topics */}
        <Link to="/admin" className="flex items-center gap-3 px-4 py-3 bg-primary-container/20 text-primary rounded-xl transition-all">
          <span className="material-symbols-outlined" style={{ fontVariationSettings: "'FILL' 1" }}>topic</span>
          <span className="font-label-md text-label-md">Topics</span>
        </Link>
        
        {/* Inactive: Settings */}
        <Link to="#" className="flex items-center gap-3 px-4 py-3 text-on-surface-variant hover:text-on-surface hover:bg-white/5 rounded-xl transition-colors">
          <span className="material-symbols-outlined">settings</span>
          <span className="font-label-md text-label-md">Settings</span>
        </Link>
      </nav>
      
      <div className="mt-auto pt-6 border-t border-outline-variant/10 flex items-center gap-3">
        <div className="w-10 h-10 rounded-full bg-surface-container-highest flex items-center justify-center overflow-hidden">
          <img alt="HaberBot Logo" className="w-full h-full object-cover" src="https://lh3.googleusercontent.com/aida-public/AB6AXuC-xGWMQJCu4M6owGofjp3Jj_FAQzt1aDLzY_culGHdszaRMf7HCBnRjUABf0jcCJVm7vJ_horVkMbKYMiUswKlyiyGgL5jnQFFAkMZ85OhoMSgNBTi_IP2yXxlr8coNgO0kiz-lhHsQwg28c8EfDuz94GKGQKL6xWl2S20DhEkNWjQtNE5Mh-XrcddLv0aru0eu0zHj_7SLhUwcgRj0ClZopN6_zumbqE1Q7ftAEy53xvONG_t6sDvlPEeIpbAmhOh23iEi2RsmfAd"/>
        </div>
        <div className="flex flex-col">
          <span className="text-label-md font-bold text-on-surface">Admin User</span>
          <span className="text-[10px] text-on-surface-variant uppercase">Super Admin</span>
        </div>
      </div>
    </aside>
  );
}
