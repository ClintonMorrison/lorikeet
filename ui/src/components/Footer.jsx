import React from 'react';

import './Footer.scss';

export default function Footer() {
  return (
    <footer className="cp-footer page-footer">
      <div className="container">
        <div className="row">
          <div className="col s8 links">
            <span className="item grey-text text-lighten-3">© 2023</span>
            <a className="item grey-text text-lighten-3" href="/terms">Terms of Service</a>
            <a className="item grey-text text-lighten-3" href="/privacy">Privacy Statement</a>
            <a className="item grey-text text-lighten-3" href="/change-log">Change Log</a>
          </div>
          <div className="col s4">
            <a className="grey-text text-lighten-4 right" href="https://ko-fi.com/T6T0VOWY">Support Us</a>
          </div>
        </div>
      </div>
    </footer>
  );
}