import React from 'react';
import _ from 'lodash';
import dayjs from 'dayjs';
import { v4 as uuid } from "uuid";
import { Helmet } from "react-helmet-async";

import Collection from "../components/passwords/Collection";
import Loader from "../components/Loader";
import TextField from "../components/forms/TextField";
import Sort from "../components/passwords/Sort";
import EmptyState from '../components/EmptyState';
import { withRouter } from '../utils/withRouter';

import './Passwords.scss';

class Passwords extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      query: '',
      sort: localStorage.getItem('sort') || 'last_used',
      document: null
    };
  }

  setSort(sort) {
    this.setState({ sort });
    localStorage.setItem('sort', sort);
  }

  applyFilter(password) {
    const query = this.state.query.toLowerCase().trim();

    if (!query) {
      return true;
    }

    const fields = ['title', 'notes', 'website'];

    return fields.some(field => password[field].toLowerCase().trim().indexOf(query) !== -1);
  }

  updateLastUsedDate(id) {
    const { documentService } = this.props.services;
    documentService.loadDocument().then(({ document }) => {
      const indexToUpdate = _.findIndex(document.passwords, { id });
      document.passwords[indexToUpdate].lastUsed = dayjs().toISOString();
      documentService.updateDocument({ document });
    });
  }

  getPasswords() {
    const { sort } = this.state;

    let passwords = _.chain(this.state)
      .get('document.passwords', [])
      .filter(p => this.applyFilter(p))
      .sortBy(this.getSortBys());

    if (sort !== 'title') {
      passwords = passwords.reverse();
    }

    return passwords.value();
  }

  getSortBys() {
    const { sort } = this.state;

    if (sort === 'title') {
      return [
        p => p.title.trim().toLowerCase(),
        p => dayjs(p.created).toISOString()
      ];
    } else if (sort === 'last_used') {
      return [
        p => dayjs(p.lastUsed).toISOString(),
        p => dayjs(p.created).toISOString()
      ];
    } else if (sort === 'updated') {
      return [
        p => dayjs(p.updated).toISOString(),
        p => dayjs(p.created).toISOString()
      ];
    } else if (sort === 'created') {
      return [
        p => dayjs(p.created).toISOString(),
        p => dayjs(p.updated).toISOString()
      ];
    }

    return null;
  }

  createPassword(e) {
    e.preventDefault();
    const { documentService } = this.props.services;

    const id = uuid();
    const password = documentService.createPassword(id);

    this.props.services.documentService.loadDocument().then(({ document }) => {
      document.passwords = [...document.passwords, password];
      return documentService.updateDocument({ document });
    }).then(() => {
      this.props.navigate(`/passwords/${id}`);
    }).catch((err) => {
      console.log('error', err);
      this.props.navigate("/logout");
    });
  }

  handleClear(e) {
    e.preventDefault();
    this.setState({ query: '' });
  }

  componentDidMount() {
    return this.props.services.documentService.loadDocument().then(({ document }) => {
      this.setState({ document });
    }).catch(() => {
      this.props.navigate("/logout");
    });
  }

  renderActions() {
    return (
      <div className="row">
        <div className="col s12">
          <div className="top-actions">
            <button
              onClick={e => this.createPassword(e)}
              className="waves-effect waves-light btn">
              <i className="material-icons left">add</i>
              Add
            </button>
          </div>

        </div>
      </div>
    );
  }

  renderFilters() {
    const passwords = this.getPasswords();

    if (!this.state.query && passwords.length === 0) {
      return null;
    }

    return (
      <form className="row" onSubmit={(e) => e.preventDefault()}>
        <TextField
          className="col s12 m6"
          label="Search"
          id="search"
          value={this.state.query}
          onChange={query => this.setState({ query })} />
        <Sort className="col s12 m6" onChange={sort => this.setSort(sort)} value={this.state.sort} />
      </form>
    );
  }

  renderContent() {
    const passwords = this.getPasswords();

    if (passwords.length === 0) {
      if (this.state.query) {
        const subtitle = (
          <p>
            We couldn't find any passwords matching your search.{' '}
            <a href="#clear" onClick={e => this.handleClear(e)}>Clear search</a>.
          </p>
        );
        return (
          <EmptyState
            img="/noSearchResults.png"
            title="We're coming up empty"
            subtitle={subtitle} />
        );
      }

      return (
        <EmptyState
          img="/emptyState.png"
          title="Nothing here but us birds!"
          subtitle="Add your first password and it will show up here." />
      );
    }

    return (
      <Collection
        passwords={passwords}
        updateLastUsedDate={id => this.updateLastUsedDate(id)} />
    );
  }


  render() {
    if (!this.state.document) {
      return (
        <div className="cp-passwords">
          <Loader />
        </div>
      );
    }

    return (
      <div className="cp-passwords">
        <Helmet>
          <title>Passwords - Lorikeet</title>
        </Helmet>

        <h1>Passwords</h1>
        {this.renderFilters()}
        {this.renderActions()}
        {this.renderContent()}
      </div>
    )
  }
}

export default withRouter(Passwords);