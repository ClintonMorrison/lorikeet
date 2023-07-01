import React from 'react';
import './MigrationNote.scss';

const MigrationNote = () => {
  return (
    <div class="cp-migration-note card-panel orange lighten-2">
      <h3 className="title">An Important Notice</h3>
      <p>
        To improve the security of this service we changing how accounts are stored, and will be removing inactive accounts as part of this process.
        {' '}
        Please sign in to your account by <strong>September 1, 2023</strong> to keep your account.
      </p>
    </div >
  )
}

export default MigrationNote
