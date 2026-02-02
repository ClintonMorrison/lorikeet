import React from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";

import Navigation from './components/Navigation';

import Home from './pages/Home';
import About from './pages/About';
import Register from './pages/Register';
import Login from './pages/Login';
import Logout from './pages/Logout';

import Account from "./pages/Account";
import Passwords from './pages/Passwords';
import View from "./pages/View";
import ChangeLog from "./pages/ChangeLog";
import Terms from "./pages/Terms";
import Privacy from "./pages/Privacy";

export default class AppRouter extends React.Component {
  render() {
    const { services } = this.props;
    return (
      <Router>
        <header>
          <Navigation services={services} />
        </header>

        <main className="container">
          <Routes>
            <Route path="/" element={<Home services={services} />} />
            <Route path="/change-log" element={<ChangeLog services={services} />} />
            <Route path="/terms" element={<Terms services={services} />} />
            <Route path="/privacy" element={<Privacy services={services} />} />
            <Route path="/about" element={<About services={services} />} />
            <Route path="/login" element={<Login services={services} />} />
            <Route path="/logout" element={<Logout services={services} />} />
            <Route path="/register" element={<Register services={services} />} />
            <Route path="/account" element={<Account services={services} />} />
            <Route path="/passwords" element={<Passwords services={services} />} />
            <Route path="/passwords/:id" element={<View services={services} />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </main>

      </Router>
    );
  }
}
