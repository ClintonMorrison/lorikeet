import React from 'react';
import { Helmet } from "react-helmet-async";

export default class Privacy extends React.Component {
  render() {
    return (
      <div className="cp-privacy">
        <Helmet>
          <title>Privacy - Lorikeet</title>
        </Helmet>


        <h1>Privacy Statement</h1>

        <p>
          Privacy is a priority of Lorikeet. We collect and store only information that is required to run this service.
        </p>

        <p>
          We store your password and account data securely on our servers.
          Your Lorikeet password, as well as all the credentials you store in Lorikeet,
          are only ever sent to our servers in an encrypted form. Your data can be unencrypted
          only with your Lorikeet password. It is not possible for us to decrypt or view your
          passwords on the server. Learn more on the <a href="/about">About page</a>.
        </p>

        <p>
          We may create copies or backups of your data to help restore your data
          in the event of unexpected software faults or server issues. After you delete
          your account we may still temporarily store a copy of your data in a backup.
          Backup archives are stored temporarily and deleted over time.
        </p>

        <p>
          There are only a few pieces of data we collect and store in an unencrypted form.
          We store your username as plain text. We also store a log of requests the server receives.
          These logs do not include data that would associate them with your account.
          They include only:

          <ul className="browser-default">
            <li>IP address</li>
            <li>Request metadata (HTTP request headers, method, and path)</li>
          </ul>

          These logs are stored only temporarily to help monitor the service and protect from attacks.
        </p>

        <p>
          Lorikeet may set cookies or values in web browser session storage when you log into the service.
        </p>

        <p>
          We do not share your data with any third parties.
        </p>
      </div>
    );

  }
}