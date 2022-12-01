import React from 'react';
import _ from 'lodash';
import TextField from "../forms/TextField";

import './PasswordField.scss';
import GeneratorModal from './GeneratorModal';

export default class PasswordField extends React.Component {
  constructor(props) {
    super(props);
    this.state = { show: false };
    this.onPasswordGenerated = this.onPasswordGenerated.bind(this);
  }

  toggleShow(e) {
    e.preventDefault();
    this.setState({ show: !this.state.show });
  }

  onPasswordGenerated(password) {
    this.setState({ show: true });
    this.onChange(password);
    setTimeout(() => {
      window.M.updateTextFields();
    }, 0);
  }

  onChange(val) {
    this.props.updateItem('password', val);
  }

  getPassword() {
    return this.props.value;
  }

  renderGenerateButton() {
    return (
      <GeneratorModal onPasswordGenerated={this.onPasswordGenerated} />
    );
  }

  renderToggleButton() {
    return (
      <button
        className="waves-effect waves-light btn"
        onClick={e => this.toggleShow(e)}>
        <i className="material-icons left">{this.state.show ? 'visibility_off' : 'visibility'}</i>
        {this.state.show ? 'Hide' : 'Show'}
      </button>
    );
  }

  render() {
    const value = this.getPassword();
    return (
      <div className="cp-password-field">
        <TextField
          label="Password"
          id="password"
          type="text" // using "password" causes Chrome to prompt to save the password
          value={value}
          masked={!this.state.show}
          autoComplete="new-password"
          error={this.props.error}
          onChange={val => this.onChange(val)} />
        {this.props.value ? this.renderToggleButton() : this.renderGenerateButton()}
      </div>
    );
  }
}