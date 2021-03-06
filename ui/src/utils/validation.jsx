import _ from 'lodash';

const digitsChars = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9'];

const specialChars = [
  '`', '~', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-',
  '_', '+', '=', '[', ']', '|', '\\', ';', ':', '\'', '"',
  '<', ',', '.', '>', '/', '?'
];


export function validatePassword(password) {
  const containsLower = _.some(password, c => c.toLowerCase() === c);
  const containsUpper = _.some(password, c => c.toUpperCase() === c);
  const containsDigit = _.some(password, c => digitsChars.includes(c));
  const containsSpecial = _.some(password, c => specialChars.includes(c));
  const valid = containsLower && containsUpper && containsDigit && containsSpecial;

  return {
    valid,
    containsLower,
    containsUpper,
    containsDigit,
    containsSpecial
  };
}