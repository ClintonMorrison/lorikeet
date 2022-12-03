import React from 'react';

import './PasswordRequirements.scss';

export default function PasswordRequirements({ result }) {
  return (
    <p className="cp-password-requirements">
      Your password must:
      <ul className="browser-default password-requirements">
        <li className={result.sufficientLength ? '' : 'invalid'}>be at least 8 characters long</li>
        <li className={result.containsLower ? '' : 'invalid'}>contains at least 1 lowercase letter</li>
        <li className={result.containsUpper ? '' : 'invalid'}>contains at least 1 capital letter</li>
        <li className={result.containsDigit ? '' : 'invalid'}>contains at least 1 number</li>
        <li className={result.containsSpecial ? '' : 'invalid'}>contains at least 1 special character (! @ # ? etc.)</li>
      </ul>
    </p>
  );
}
