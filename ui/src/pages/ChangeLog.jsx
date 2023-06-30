import React from 'react';
import { Helmet } from "react-helmet-async";

export default class ChangeLog extends React.Component {
  render() {
    return (
      <div className="cp-change-log">
        <Helmet>
          <title>Change Log - Lorikeet</title>
        </Helmet>

        <h1>Change Log</h1>

        <h2>Saturday, July 12, 2023</h2>
        <ul className="browser-default">
          <li>security improvements</li>
          <li>bug fixes</li>
          <li>updated assets</li>
          <li>posted notice about account migration</li>
        </ul>


        <h2>Saturday, Dec 3, 2022</h2>
        <ul className="browser-default">
          <li>add Recaptcha</li>
          <li>security improvements</li>
          <li>bug fixes</li>
        </ul>


        <h2>Friday, May 18, 2022</h2>
        <ul className="browser-default">
          <li>fixed bug with default sort order not being applied</li>
          <li>improved how server handles sessions</li>
        </ul>



        <h2>Monday, November 9, 2020</h2>
        <ul className="browser-default">
          <li>fixed bug with text area contrast on dark mode</li>
        </ul>


        <h2>Monday, October 5, 2020</h2>
        <ul className="browser-default">
          <li>added changelog</li>
          <li>added dark mode</li>
          <li>fixed some server-side bugs related to concurrency</li>
          <li>added more options for generating passwords</li>
        </ul>
      </div>
    );

  }
}