import React from 'react';
import _ from 'lodash';
import { Helmet } from "react-helmet";

import TextField from '../components/forms/TextField';

export default class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      username: "",
      password: "",
      usernameError: "",
      passwordError: ""
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

    if (!isOkay) {
      return;
    }

    const { username, password } = this.state;
    this.props.services.authService.setCredentials({ username, password });
    this.props.services.documentService.loadDocument()
      .then(() => {
        this.props.history.push('/passwords');
      })
      .catch(err => {
        console.log(err);
        const errorMessage = _.get(err, 'response.data.error', 'An error occurred.');
        if (errorMessage) {
          this.setState({ usernameError: errorMessage, passwordError: ' ' });
        }
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

            <div className="row">
              <div className="input-field col s12">
                <button
                  className="btn waves-effect waves-light"
                  type="submit"
                  name="action"
                  onClick={(e) => this.submit(e)}>
                  Login
                </button>
              </div>
            </div>

          </form>
        </div>
      </div>
    );

  }
}