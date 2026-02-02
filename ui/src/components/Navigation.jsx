import React from 'react';
import { Link } from 'react-router-dom';

import './Navigation.scss';

const loggedOutItems = [
  <li key="about"><Link to="/about">About</Link></li>,
  <li key="register"><Link to="/register">Sign Up</Link></li>, // className="btn waves-effect waves-light"
  <li key="login"><Link to="/login">Login</Link></li>,
];

const loggedInItems = [
  <li key="account"><Link to="/account">My Account</Link></li>,
  <li key="passwords"><Link to="/passwords">My Passwords</Link></li>,
  <li key="logout"><Link to="/logout">Logout</Link></li>
];


export class Navigation extends React.Component {
  constructor(props) {
    super(props);
    this.sidebarRef = React.createRef();
  }

  componentDidMount() {
    setTimeout(() => {
      window.M.Sidenav.init(this.sidebarRef.current);
    }, 0);
  }

  componentDidUpdate() {
    const instance = window.M.Sidenav.getInstance(this.sidebarRef.current);
    if (instance) {
      instance.close();
    }
  }

  render() {
    const { services } = this.props;
    const loggedIn = services.authService.sessionExists();

    const items = loggedIn ? loggedInItems : loggedOutItems;

    return (
      <div className="cp-navigation">
        <nav>
          <div className="nav-wrapper">
            <Link to="/" className="brand-logo">
              Lorikeet
            </Link>
            <button type="button" data-target="mobile-demo" className="sidenav-trigger right hide-on-med-and-up">
              <i className="material-icons">menu</i>
            </button>
            <ul className="right hide-on-small-and-down">
              {items}
            </ul>
          </div>
        </nav>

        <ul className="sidenav" id="mobile-demo" ref={this.sidebarRef}>
          {items}
        </ul>
      </div>
    );
  }
}

export default Navigation;