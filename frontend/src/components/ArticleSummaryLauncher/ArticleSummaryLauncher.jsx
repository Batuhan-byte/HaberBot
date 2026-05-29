import React from 'react';
import './ArticleSummaryLauncher.css';

/**
 * Sticky/Floating Premium CTA Button to launch the AI Summary modal.
 * Guides non-authenticated users to login seamlessly by integrating
 * with the site's existing authentication event listener.
 */
export default function ArticleSummaryLauncher({ onClick }) {
  return (
    <div className="summary-launcher-container">
      <button 
        className="summary-cta-button" 
        onClick={onClick}
        aria-label="Yapay Zeka Özeti Oku"
        id="summary-launcher-button"
      >
        <span className="material-symbols-outlined summary-cta-sparkle" aria-hidden="true">auto_awesome</span>
        <span>Yapay Zeka Özeti</span>
      </button>
    </div>
  );
}
