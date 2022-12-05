import React from 'react';
import _ from 'lodash';
import { Helmet } from "react-helmet-async";

import TextField from '../components/forms/TextField';

import './Register.scss';
import { validatePassword } from "../utils/validation";
import PasswordRequirements from "../components/PasswordRequirements";
import ReCaptcha from '../components/ReCaptcha';

export default class Register extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: "",
      password: "",
      usernameError: "",
      passwordError: "",
      recaptchaResult: "",
      resetRecaptcha: false,
      passwordValidation: {
        valid: true,
        sufficientLength: true,
        containsLower: true,
        containsUpper: true,
        containsDigit: true,
        containsSpecial: true
      },
      darkMode: props.services.preferencesService.isDarkModeEnabled(),
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

    const validation = this.state.passwordValidation;
    if (!validation.valid) {
      this.setState({
        passwordError: 'Password does not meet requirements',
      });
      isOkay = false;
    }

    if (!this.state.recaptchaResult) {
      isOkay = false;
    }

    if (!isOkay) {
      return;
    }

    const { username, password, recaptchaResult } = this.state;
    this.props.services.documentService.createDocument({ username, password, recaptchaResult })
      .then((resp) => {
        this.props.history.push('/passwords')
      })
      .catch(err => {
        console.log({ ...err });
        this.resetRecaptcha();
        const errorMessage = _.get(err, 'response.data.error', 'An error occurred.');
        if (errorMessage) {
          this.setState({ usernameError: errorMessage });
        }
        this.props.services.authService.logout();
      });
  }

  clearErrors() {
    this.setState({
      usernameError: "",
      passwordError: "",
    });
  }

  updateUsername(username) {
    this.clearErrors();
    this.setState({ username });
  }

  updatePassword(password) {
    this.clearErrors();
    const passwordValidation = validatePassword(password);
    this.setState({ password, passwordValidation });
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
      <div className="cp-register">
        <Helmet>
          <title>Sign Up - Lorikeet</title>
        </Helmet>

        <h1>Sign Up</h1>

        <div className="row">
          <form className="col s12">
            <div className="row">
              <div className="col s12">
                <p>
                  By signing up you agree to our <a href="/terms">Terms of Service</a>.
                </p>

                <p>
                  Enter a username and a strong password for your new account.
                  Please write down your account information and keep it safe.

                  Because of how your data will be encrypted, it will not be possible to regain
                  control of your account if you forget.
                </p>
                <PasswordRequirements result={this.state.passwordValidation} />
              </div>
            </div>

            <TextField
              label="Username"
              id="username"
              value={this.state.username}
              error={this.state.usernameError}
              onChange={val => this.updateUsername(val)} />

            <TextField
              label="Password"
              id="password"
              type="password"
              value={this.state.password}
              error={this.state.passwordError}
              onChange={val => this.updatePassword(val)} />

            <ReCaptcha
              onChange={result => this.updateRecaptchaResult(result)}
              reset={this.state.resetRecaptcha}
              darkMode={this.state.darkMode} />

            <div className="row">
              <div className="input-field col s12">
                <button
                  className="btn waves-effect waves-light"
                  type="submit"
                  name="action"
                  disabled={!this.state.recaptchaResult}
                  onClick={(e) => this.submit(e)}>
                  Register
                </button>
              </div>
            </div>

          </form>
        </div>
      </div>
    );

  }
}