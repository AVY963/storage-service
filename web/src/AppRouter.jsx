import React from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "./pages/Login";
import FilesPage from "./pages/Files";
import { useAuth } from "./AuthContext";

function PrivateRoute({ children }) {
  const { user, loading } = useAuth();
  if (loading) return <div className="text-center mt-10">Загрузка...</div>;
  return user ? children : <Navigate to="/login" replace />;
}

export default function AppRouter() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/files" element={<PrivateRoute><FilesPage /></PrivateRoute>} />
        <Route path="*" element={<Navigate to="/files" replace />} />
      </Routes>
    </Router>
  );
}