import React from 'react';
import { Outlet } from 'react-router-dom';
import Header from '../Header/Header';
import Footer from '../Footer/Footer';
import ConsoleWidget from '../ConsoleWidget/ConsoleWidget';
import './Layout.css';

const Layout = () => {
  return (
    <div className="app-layout">
      <Header />
      <main className="main-content">
        <Outlet />
      </main>
      <ConsoleWidget />
      <Footer />
    </div>
  );
};

export default Layout;
