import React from 'react';
import './MigrationNote.scss';

const MigrationNote = () => {
  return (
    <div class="cp-migration-note card-panel orange lighten-2">
      <h3 className="title">Important Notice</h3>
      <p>

        To improve security for this service, we are changing how our users' accounts are stored -- including the <strong>removal of inactive accounts</strong>.
        Be sure to login any time before the deadline of <strong>September 1, 2023</strong>, so that we can continue to keep your passwords secure.
        Thank you for continuing to use Lorikeet!
      </p>
      <p>
        -- Clinton & Emma
      </p>
    </div >
  )
}

export default MigrationNote
