import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import HomePage from './pages/HomePage/HomePage';
import ArticlePage from './pages/ArticlePage/ArticlePage';
import TopicPage from './pages/TopicPage/TopicPage';
import SearchPage from './pages/SearchPage/SearchPage';
import LoginPage from './pages/LoginPage/LoginPage';
import RegisterPage from './pages/RegisterPage/RegisterPage';
import ProtectedRoute from './components/ProtectedRoute';
import AdminDashboard from './pages/AdminPage/AdminDashboard';
import AdminSources from './pages/AdminPage/AdminSources';
import AdminContent from './pages/AdminPage/AdminContent';
import NotFoundPage from './pages/NotFoundPage/NotFoundPage';
import AuthModal from './components/AuthModal/AuthModal';
import './App.css';

// Initialize TanStack Query client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes cache
      refetchOnWindowFocus: false,
      retry: 2,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          {/* Public Routes */}
          <Route path="/" element={<HomePage />} />
          <Route path="/kategori/:slug" element={<TopicPage />} />
          <Route path="/haber/:id" element={<ArticlePage />} />
          <Route path="/arama" element={<SearchPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          
          {/* Admin Routes */}
          <Route path="/admin" element={
            <ProtectedRoute allowedRoles={['Admin']}>
              <AdminDashboard />
            </ProtectedRoute>
          } />
          <Route path="/admin/sources" element={
            <ProtectedRoute allowedRoles={['Admin']}>
              <AdminSources />
            </ProtectedRoute>
          } />
          <Route path="/admin/content" element={
            <ProtectedRoute allowedRoles={['Admin']}>
              <AdminContent />
            </ProtectedRoute>
          } />
          
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
        
        {/* Global Auth Modal for tabbed Login/Register Popup */}
        <AuthModal />
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;

