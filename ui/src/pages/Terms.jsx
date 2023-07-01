import React from 'react';
import { Helmet } from "react-helmet-async";

export default class Terms extends React.Component {
  render() {
    return (
      <div className="cp-terms">
        <Helmet>
          <title>Terms - Lorikeet</title>
        </Helmet>

        <h1>Terms of Service</h1>

        <p>
          Lorikeet is a small personal project. As such, these terms of service are somewhat informal.
          The gist of it is that we will do our best and try to be reasonable, but don't hold it against us if things go wrong!
        </p>

        <h2>Accounts & Data</h2>
        <p>
          To use this service you must create an account on the <a href="/register">registration page</a>.
          It is your responsibility to remember your username and password for this service.
          If you lose your credentials, we cannot help recover your account or data due to how your data is securely stored.
          You are also responsible for using a unique password and not sharing it with others.
          This is essential for keeping your data secure.
          We recommend you change your password for the service every few months.
        </p>

        <p>
          By using this service you share your password data with us.
          You agree to allow us to store, retrieve, backup, restore, and copy this data.
          You retain ownership and all rights to your data. We can delete it at your request (see the "Termination" section).

          All data associated with your account is stored in an encrypted form, except your username and server logs.
          See the <a href="/about">About page</a> and <a href="/privacy">Privacy Statement</a> for more information.
        </p>

        <h2>Disclaimer & Limitation of Liability</h2>

        <p>
          We will always make our best effort to keep
          the service available and working. However, as a free personal project,
          we cannot give a strong guarantee that the service will always be available, secure, and free of bugs.
          We backup your data regularly but may lose it as a result of server failure or software defects.
          We suggest you periodically export your data and store it outside of our service.
          We will do our best, but do not make any guarantees about the availability, functionality, or security of the service.
          Lorikeet is provided at on an "AS IS" and "AS AVAILABLE" basis. You use it at our own risk.
        </p>

        <p>
          By using Lorikeet you agree to not hold Lorikeet, its founders, employees, or any contributors to the project, liable
          for any losses or damages directly or indirectly related to using the service or site.
        </p>

        <h2>Termination</h2>
        <p>
          You may stop using the service at any time.
          You can request that we permanently delete your account and password data by logging
          in and clicking the "Delete All Data" button on the "My Account" page.
          Older copies of your data may temporarily be stored in backups.
          See our <a href="/privacy">Privacy Statement</a> for more information.
        </p>

        <h2>Changes</h2>
        <p>
          In the future we may make changes to this agreement, the privacy statement, or the service. We may also cease providing this service.
          For example, we may add paid features, or permanently shut down the service.
          For any such changes, we will provide 30 days notice on the website.
          By continuing to use the service you agree to these new terms.
        </p>

        <h2>Inactive Account Removal</h2>
        <p>
          On June 30, 2023 we posted a notice that inactive Lorikeet accounts will be removed. Users will have until September 1, 2023 to sign in to their account.
        </p>
        <p>
          Any accounts that were not accessed between June 30, 2023 and September 1, 2023 will be removed.
        </p>
      </div>
    );

  }
}