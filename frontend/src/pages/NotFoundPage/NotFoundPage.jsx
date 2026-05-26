import React from 'react';
import { Link } from 'react-router-dom';
import './NotFoundPage.css';

const NotFoundPage = () => {
  return (
    <div className="not-found-page">
      <div className="not-found-container">
        <span className="not-found-emoji">🔍</span>
        <h1 className="not-found-code">404</h1>
        <h2 className="not-found-title">Sayfa Bulunamadı</h2>
        <p className="not-found-desc">
          Aradığınız sayfa silinmiş, taşınmış veya hiç var olmamış olabilir.
        </p>
        <Link to="/" className="home-link-btn">
          Ana Sayfaya Dön
        </Link>
      </div>
    </div>
  );
};

export default NotFoundPage;
