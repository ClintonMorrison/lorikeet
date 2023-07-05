import React from 'react';
import './MigrationNote.scss';

const MigrationNote = () => {
  return (
    <div class="cp-migration-note card-panel orange lighten-2">
      <h3 className="title">Important Notice</h3>
      <p>

        To improve the security of this service we are changing how accounts are stored, and will be <strong>removing inactive accounts</strong> as part of this process.
        Please sign in to your account by <strong>September 1, 2023</strong> to keep your account.
      </p>
    </div >
  )
}

export default MigrationNote
