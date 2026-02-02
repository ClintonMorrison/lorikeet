import React from 'react';
import { Navigate } from 'react-router-dom';

export default function Logout({ services }) {
  services.sessionService.deleteSession();
  services.authService.logout();
  return <Navigate to="/login" replace />;
}
