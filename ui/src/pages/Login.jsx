import React from 'react';
import _ from 'lodash';
import { Helmet } from "react-helmet-async";

import TextField from '../components/forms/TextField';
import ReCaptcha from '../components/ReCaptcha';
import { isLocalDev } from '../utils/validation';
import MigrationNote from '../components/MigrationNote';

export default class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: "",
      password: "",
      usernameError: "",
      passwordError: "",
      recaptchaResult: "",
      resetRecaptcha: false,
      darkMode: props.services.preferencesService.isDarkModeEnabled(),
      isLocalDev: isLocalDev(),
    };
  }

  submit(e) {
    e.preventDefault();

    let isOkay = true;

    if (!this.state.username) {
      this.setState({ usernameError: "Username cannot be empty" });
      isOkay = false;
    }

    if (!this.state.password) {
      this.setState({ passwordError: "Password cannot be empty" });
      isOkay = false;
    }

    if (!this.state.recaptchaResult && !this.state.isLocalDev) {
      isOkay = false;
    }

    if (!isOkay) {
      return;
    }

    const { username, password, recaptchaResult } = this.state;
    this.props.services.authService.setCredentials({ username, password });
    this.props.services.sessionService.createSession({ recaptchaResult })
      .then((resp) => {
        this.resetRecaptcha();
        this.props.services.authService.setCredentials({ password, salt: resp.data.salt });
        return this.props.services.documentService.loadDocument({ password });
      })
      .then(({ document, needsMigration }) => {
        if (needsMigration) {
          return this.props.services.documentService.migrateDocument({ document, password }).then(() => {
            this.props.history.push('/passwords');
          })
        } else {
          this.props.history.push('/passwords');
        }
      })
      .catch(err => {
        this.resetRecaptcha();
        console.log(err);
        const errorMessage = _.get(err, 'response.data.error', 'An error occurred.');
        if (errorMessage) {
          this.setState({ usernameError: errorMessage, passwordError: ' ' });
        }
        this.props.services.authService.logout();
      });
  }

  clearErrors() {
    this.setState({
      usernameError: "",
      passwordError: ""
    });
  }

  updateUsername(username) {
    this.clearErrors();
    this.setState({ username });
  }

  updatePassword(password) {
    this.clearErrors();
    this.setState({ password });
  }

  updateRecaptchaResult(recaptchaResult) {
    this.clearErrors();
    this.setState({ recaptchaResult, resetRecaptcha: false });
  }

  resetRecaptcha() {
    this.setState({ recaptchaResult: '', resetRecaptcha: true });
  }

  render() {
    return (
      <div className="cp-login">
        <Helmet>
          <title>Login - Lorikeet</title>
        </Helmet>

        <h1>Login</h1>

        <div className="row">
          <form className="col s12">
            <div className="row">
              <div className="col s12">
                Enter your username and password to login. <a href="/register">Don't have an account yet?</a>
              </div>
            </div>

            <TextField
              label="Username"
              id="username"
              autoComplete="username"
              value={this.state.username}
              error={this.state.usernameError}
              onChange={e => this.updateUsername(e)} />

            <TextField
              label="Password"
              id="password"
              autoComplete="password"
              type="password"
              value={this.state.password}
              error={this.state.passwordError}
              onChange={e => this.updatePassword(e)} />

            <ReCaptcha
              onChange={result => this.updateRecaptchaResult(result)}
              reset={this.state.resetRecaptcha}
              darkMode={this.state.darkMode} />

            <div className="row">
              <div className="input-field col s12">
                <button
                  className="btn waves-effect waves-light"
                  disabled={!this.state.recaptchaResult && !this.state.isLocalDev}
                  type="submit"
                  name="action"
                  onClick={(e) => this.submit(e)}>
                  Login
                </button>
              </div>
            </div>

          </form>
        </div>

        <div className="row">
          <MigrationNote />
        </div>
      </div>
    );

  }
}