import React from 'react';
import { Helmet } from "react-helmet-async";
import MigrationNote from '../components/MigrationNote';

export default class About extends React.Component {
  render() {
    return (
      <div className="cp-about">
        <Helmet>
          <title>About - Lorikeet</title>
        </Helmet>

        <h1>About</h1>

        <div className="row">
          <MigrationNote />
        </div>

        <p>
          Lorikeet is a personal project created by <a href="https://clintonmorrison.com">Clinton Morrison</a>{' '}
          and Emma MacDonald.

          It is designed to be a simple and secure online password manager.
          We weren't happy with many password managers out there, so with our love of making cool new things together, we made Lorikeet.
        </p>

        <p>
          All the code for Lorikeet is on <a href="https://github.com/ClintonMorrison/lorikeet">GitHub</a>.
          Lorikeet is built with <a href="https://golang.org/">Golang</a>, <a href="https://reactjs.org/">React</a>, and <a href="https://materializecss.com">Materialize CSS</a>.
        </p>

        <h2>Security</h2>

        <p>
          Security was our main focus in Lorikeet, which led to some unconventional decisions.
          Unlike other online services, you cannot reset your password if you forget it.
          This is because of how your passwords are encrypted; we cannot reset your password
          because we do not store your password. We don't have the ability to decrypt your passwords.
        </p>


        <p>
          If you enter incorrect credentials several times in a row, your account will be locked for a few hours.
          This helps keep your account safe by making it difficult for others to guess your password.
        </p>

        <p>
          Neither your login information nor your stored passwords are ever sent to the server.

          When you login, an obfuscated copy of your Lorikeet username and password are stored temporarily in your browser session storage.
          A special token derived from your credentials is sent to the server when you login.
        </p>

        <p>
          When you add passwords, they are encrypted in the browser using your Lorikeet credentials.
          Our servers only ever receive this encrypted password data, and never your real passwords.
          On the server it is encrypted a second time with your token, your username, and a salt.

          When we retrieve your passwords, they are partially decrypted on the server with your token.

          In the browser, the original passwords are recovered using your credentials.
        </p>

        <p>
          We use strong AES encryption, and SHA256 hashing.
          All requests to the server are also sent over HTTPS.
          You can find more information about how Lorikeet works on{' '}
          <a href="https://ClintonMorrison.com/blog/creating_password_manager">Clinton's blog</a> and on the{' '}
          <a href="https://github.com/ClintonMorrison/lorikeet">GitHub</a> page.
        </p>

        <h2>Contact</h2>
        <p>
          If you have feedback, suggestions, or comments, we would love to hear from you.
          You can send us an email at{' '}
          <a href="mailto:contact@lorikeet.ca">contact@lorikeet.ca</a>.
        </p>

        <h2>Support Us</h2>

        <p>
          We created Lorikeet for fun; not for profit. It is currently not monetized in any way.
          If you like Lorikeet, you can <a href="https://ko-fi.com/T6T0VOWY">support us on Ko-fi</a>.
        </p>
      </div>
    );

  }
}