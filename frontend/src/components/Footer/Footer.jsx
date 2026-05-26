import React from 'react';
import './Footer.css';

const Footer = () => {
  return (
    <footer className="site-footer">
      <div className="footer-container">
        <p className="footer-copy">
          &copy; {new Date().getFullYear()} HaberBot. Tüm hakları saklıdır.
        </p>
        <p className="footer-power">
          Google Gemini AI tarafından desteklenmektedir.
        </p>
      </div>
    </footer>
  );
};

export default Footer;
