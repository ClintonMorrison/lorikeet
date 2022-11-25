import React from 'react';
import Footer from './components/Footer';
import Routes from './Routes';

// Services
import AuthService from "./services/AuthService";
import APIService from "./services/APIService";
import DocumentService from "./services/DocumentService";
import SessionService from "./services/SessionService";

import './App.scss';
import PreferencesService from './services/PreferencesService';

// Instantiate services
const authService = new AuthService();

const apiService = new APIService({
  authService,
  baseURL: `${window.location.origin}/api/`
});


const documentService = new DocumentService({
  authService,
  apiService
});

const sessionService = new SessionService({
  apiService,
  authService
})

const preferencesService = new PreferencesService({
  onDarkModeChanged: () => { }
});

const services = {
  apiService,
  authService,
  documentService,
  sessionService,
  preferencesService
};

window.services = services;

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      darkMode: preferencesService.isDarkModeEnabled()
    }

    preferencesService.onDarkModeChanged = (darkMode) => {
      this.setState({ darkMode: darkMode });
    }
  }
  render() {
    return (
      <div className={`cp-app ${this.state.darkMode ? 'dark-mode' : ''}`}>
        <Routes services={services} />
        <Footer />
      </div>
    );
  }
}

