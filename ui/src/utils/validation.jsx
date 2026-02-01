import _ from 'lodash';

const digitsChars = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'];

const specialChars = [
  '`', '~', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-',
  '_', '+', '=', '[', ']', '|', '\\', ';', ':', '\'', '"',
  '<', ',', '.', '>', '/', '?'
];

export function isLocalDev() {
  return window.location.host === 'localhost:3000'; // skip some validation for local dev
}


export function validatePassword(password) {
  const skipValidation = isLocalDev();
  const sufficientLength = password.length >= 8;
  const containsLower = _.some(password, c => c.toLowerCase() === c);
  const containsUpper = _.some(password, c => c.toUpperCase() === c);
  const containsDigit = _.some(password, c => digitsChars.includes(c));
  const containsSpecial = _.some(password, c => specialChars.includes(c));
  const valid = skipValidation || (containsLower && containsUpper && containsDigit && containsSpecial && sufficientLength);

  return {
    valid,
    sufficientLength,
    containsLower,
    containsUpper,
    containsDigit,
    containsSpecial
  };
}