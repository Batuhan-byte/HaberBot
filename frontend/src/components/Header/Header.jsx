import React, { useState } from 'react';
import { Link, NavLink } from 'react-router-dom';
import { useTopics } from '../../hooks/useTopics';
import './Header.css';

const Header = () => {
  const { data, isLoading } = useTopics();
  const [isOpen, setIsOpen] = useState(false);

  const topics = data?.topics || [];

  return (
    <header className="site-header">
      <div className="header-container">
        <Link to="/" className="logo-link" onClick={() => setIsOpen(false)}>
          <span className="logo-emoji">🤖</span>
          <span className="logo-text">HaberBot</span>
        </Link>

        <button 
          className={`mobile-menu-toggle ${isOpen ? 'is-active' : ''}`}
          onClick={() => setIsOpen(!isOpen)}
          aria-label="Menüyü Aç/Kapat"
        >
          <span className="bar"></span>
          <span className="bar"></span>
          <span className="bar"></span>
        </button>

        <nav className={`header-nav ${isOpen ? 'is-open' : ''}`}>
          <NavLink 
            to="/" 
            className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
            onClick={() => setIsOpen(false)}
            end
          >
            Ana Sayfa
          </NavLink>
          
          {!isLoading && topics.map((topic) => (
            <NavLink 
              key={topic.ID}
              to={`/konu/${topic.Slug}`} 
              className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
              onClick={() => setIsOpen(false)}
            >
              {topic.Name}
            </NavLink>
          ))}
        </nav>
      </div>
    </header>
  );
};

export default Header;
