import React from 'react';
import { Link } from 'react-router-dom';
import { Helmet } from "react-helmet-async";

import './Home.scss';
import MainLogo from "../components/icons/MainLogo";
import Heart from "../components/icons/Heart";
import Bloom from "../components/icons/Bloom";
import PadLock from "../components/icons/PadLock";
import MigrationNote from '../components/MigrationNote';

export default function Home({ services }) {

  const cta = services.authService.sessionExists() ?
    <Link to="/passwords" className="sign-up-link waves-effect waves-light btn-large btn">View Your Passwords</Link> :
    <Link to="/register" className="sign-up-link waves-effect waves-light btn-large btn">Sign Up Now</Link>;


  return (
    <div className="cp-home">
      <Helmet>
        <title>Lorikeet</title>
      </Helmet>

      <div className="bird-banner">
        <MainLogo />
      </div>

      <div className="heading">
        <h1>Lorikeet</h1>
        <p className="subtitle">A secure online password manager.</p>
      </div>

      <div className="row">
        <div className="col s12 m4">
          <div className="center promo promo-example">
            <Bloom />
            <h5 className="promo-caption">Easy</h5>
            <p className="light center">
              You can stop keeping track of your passwords. It's easy to manage your passwords with Lorikeet.
            </p>
          </div>
        </div>

        <div className="col s12 m4">
          <div className="center promo promo-example">
            <PadLock />
            <h5 className="promo-caption">Secure</h5>
            <p className="light center">
              With strong AES encryption on the client-side and server-side, you don't need to worry about your passwords.
            </p>
          </div>
        </div>

        <div className="col s12 m4">
          <div className="center promo promo-example">
            <Heart />
            <h5 className="promo-caption">Free</h5>
            <p className="light center">
              Lorikeet is free to use, and <a href="https://github.com/ClintonMorrison/lorikeet">open source</a>.
              If you like it you can <a href="https://ko-fi.com/T6T0VOWY">support us on Ko-fi</a>.
            </p>
          </div>
        </div>
      </div>
      {cta}
      <div className="row">
        <MigrationNote />
      </div>
    </div>
  );
}