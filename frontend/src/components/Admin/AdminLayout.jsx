import { NavLink } from 'react-router-dom';
import '../../pages/AdminPage/AdminContent.css';

const baseLinkClass = 'flex items-center gap-3 py-2.5 px-3 rounded transition-all duration-150 cursor-pointer active:opacity-80 font-label-md text-label-md';
const inactiveLinkClass = 'text-zinc-400 hover:bg-zinc-900 hover:text-white';
const activeLinkClass = 'bg-blue-950/40 text-blue-400 border-r-2 border-blue-500 font-bold';

const getNavLinkClass = ({ isActive }) => `${baseLinkClass} ${isActive ? activeLinkClass : inactiveLinkClass}`;

export default function AdminLayout({
  title,
  subtitle,
  actions,
  children,
}) {
  return (
    <div className="flex min-h-screen bg-black text-white">
      <aside className="h-screen w-64 fixed left-0 top-0 bg-zinc-950 border-r border-zinc-800 flex flex-col py-6 px-4 gap-4 z-50">
        <div className="flex flex-col gap-1 mb-6 px-2">
          <h1 className="font-headline-md text-headline-md text-primary tracking-tight font-extrabold">
            Haber<span className="text-blue-500">Bot</span> Admin
          </h1>
          <p className="font-label-sm text-label-sm text-zinc-500">v1.1.0 — Admin Panel</p>
        </div>
        <nav className="flex flex-col gap-1 flex-grow">
          <NavLink to="/admin" end className={getNavLinkClass}>
            <span className="material-symbols-outlined text-[20px]" data-icon="dashboard">dashboard</span>
            Dashboard
          </NavLink>
          <NavLink to="/admin/sources" className={getNavLinkClass}>
            <span className="material-symbols-outlined text-[20px]" data-icon="rss_feed">rss_feed</span>
            RSS Sources
          </NavLink>
          <NavLink to="/admin/content" className={getNavLinkClass}>
            <span className="material-symbols-outlined text-[20px]" data-icon="article">article</span>
            Content Management
          </NavLink>
          <div className="flex items-center gap-3 py-2.5 px-3 rounded text-zinc-600 cursor-not-allowed opacity-50 font-label-md text-label-md select-none">
            <span className="material-symbols-outlined text-[20px]">group</span>
            Users
            <span className="text-[8px] bg-zinc-900 text-zinc-500 px-1.5 py-0.5 rounded font-mono uppercase ml-auto tracking-wider font-extrabold">Kilitli</span>
          </div>
          <div className="flex items-center gap-3 py-2.5 px-3 rounded text-zinc-600 cursor-not-allowed opacity-50 font-label-md text-label-md select-none">
            <span className="material-symbols-outlined text-[20px]">terminal</span>
            Logs
            <span className="text-[8px] bg-zinc-900 text-zinc-500 px-1.5 py-0.5 rounded font-mono uppercase ml-auto tracking-wider font-extrabold">Kilitli</span>
          </div>
          <NavLink to="/" className={getNavLinkClass}>
            <span className="material-symbols-outlined text-[20px]" data-icon="home">home</span>
            Siteye Dön
          </NavLink>
        </nav>
      </aside>

      <main className="flex-grow ml-64 p-8 min-h-screen flex flex-col bg-[#050505]">
        <header className="flex justify-between items-end mb-8">
          <div>
            <h2 className="text-3xl font-extrabold text-white tracking-tight">{title}</h2>
            {subtitle ? (
              <p className="text-zinc-400 mt-1 text-sm font-medium">{subtitle}</p>
            ) : null}
          </div>
          {actions ? (
            <div className="flex items-center gap-4">
              {actions}
            </div>
          ) : null}
        </header>

        {children}
      </main>
    </div>
  );
}
