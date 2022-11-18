import React from 'react';
import { Link } from 'react-router-dom';
import TextField from "../forms/TextField";
import TextArea from "../forms/TextArea";

import './Details.scss';
import PasswordField from './PasswordField';

export default class Details extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      item: {
        title: '',
        username: '',
        password: '',
        email: '',
        website: '',
        notes: ''
      },
      errors: {}
    };

    this.updateItem = this.updateItem.bind(this);
  }

  updateItem(field, value) {
    this.setState({
      item: {
        ...this.state.item,
        [field]: value
      }
    });
  }

  componentDidMount() {
    this.setState({
      item: { ...this.props.item }
    });

  }

  handleSave(e) {
    e.preventDefault();
    this.props.updatePassword(this.state.item);
  }

  handleDelete(e) {
    e.preventDefault();

    if (window.confirm('Are you sure you want to delete this password? This cannot be undone.')) {
      return this.props.deletePassword(this.state.item);
    }
  }

  render() {
    const { item } = this.state;

    if (!item) {
      return null;
    }

    return (
      <div className="cp-details">
        <h1 className="title">{item.title || 'Untitled'}</h1>

        <form>
          <TextField
            autoFocus
            label="Title"
            id="title"
            value={item.title}
            error={this.state.errors.title}
            onChange={val => this.updateItem('title', val)} />

          <TextField
            label="Username"
            id="username"
            value={item.username}
            error={this.state.errors.username}
            autoComplete="off"
            onChange={val => this.updateItem('username', val)} />

          <TextField
            label="Email"
            id="email"
            type="text"
            value={item.email}
            error={this.state.errors.email}
            autoComplete="off"
            onChange={val => this.updateItem('email', val)} />

          <PasswordField
            value={item.password}
            updateItem={this.updateItem} />

          <TextField
            label="Website"
            id="website"
            type="text"
            value={item.website}
            error={this.state.errors.website}
            onChange={val => this.updateItem('website', val)} />

          <TextArea
            label="Notes"
            id="notes"
            type="text"
            value={item.notes}
            error={this.state.errors.notes}
            onChange={val => this.updateItem('notes', val)} />

          <div className="actions">
            <Link className="waves-effect waves-light btn-secondary btn" to="/passwords">Cancel</Link>
            <button
              className="waves-effect waves-light btn"
              onClick={(e) => this.handleSave(e)}>Save</button>
            <button
              className="waves-effect waves-light btn btn-negative delete-button"
              onClick={(e) => this.handleDelete(e)}>
              <i className="material-icons">delete</i>
            </button>
          </div>
        </form>
      </div>
    );
  }
}
