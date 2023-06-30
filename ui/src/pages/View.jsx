import React from 'react';
import { Redirect } from 'react-router-dom';
import { Helmet } from "react-helmet-async";
import _ from 'lodash';

import Details from "../components/passwords/Details";
import Loader from "../components/Loader";
import moment from "moment/moment";

export default class View extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      document: null
    };
  }

  updatePassword(item) {
    const { id } = item;
    item.updated = moment().toISOString();
    item.lastUsed = moment().toISOString();

    this.props.services.documentService.loadDocument().then(({ document }) => {
      const indexToUpdate = _.findIndex(document.passwords, { id });
      document.passwords[indexToUpdate] = item;
      return this.props.services.documentService.updateDocument({ document });
    }).then(() => {
      this.props.history.push("/passwords");
    }).catch((err) => {
      console.log('error', err);
      this.props.history.push("/logout");
    });
  }

  deletePassword(password) {
    const { id } = password;
    this.props.services.documentService.loadDocument().then(({ document }) => {
      const indexToDelete = _.findIndex(document.passwords, { id });
      document.passwords.splice(indexToDelete, 1);
      this.props.services.documentService.updateDocument({ document });
    }).then(() => {
      this.props.history.push("/passwords");
    }).catch((err) => {
      console.log('error', err);
      this.props.history.push("/logout");
    });
  }

  getPasswords() {
    return _.get(this.state, 'document.passwords', []);
  }

  componentDidMount() {
    return this.props.services.documentService.loadDocument().then(({ document }) => {
      this.setState({ document });
    }).catch(() => {
      this.props.history.push("/logout");
    });
  }

  render() {
    if (!this.state.document) {
      return (
        <Loader />
      );
    }

    const { id } = this.props.match.params;
    const passwords = this.getPasswords();
    const item = _.find(passwords, { id });

    if (!item) {
      return <Redirect to="/passwords" />;
    }

    return (
      <div className="cp-view">
        <Helmet>
          <title>Passwords - Lorikeet</title>
        </Helmet>

        <Details
          item={item}
          createPassword={() => this.createPassword()}
          updatePassword={p => this.updatePassword(p)}
          deletePassword={p => this.deletePassword(p)} />
      </div>
    );
  }
}