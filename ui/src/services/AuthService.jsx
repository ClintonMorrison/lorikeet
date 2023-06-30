import sha256 from 'crypto-js/sha256';
import AES from 'crypto-js/aes';
import UTF_8 from 'crypto-js/enc-utf8';
import _ from 'lodash';

const PEPPER_1 = 'CC352C99A14616AD22678563ECDA5';
const PEPPER_2 = '7767B9225CF66B418DD2A39CBC4AA';

export default class AuthService {
  constructor() {
  }

  firstHash(password) {
    const username = this.getUsername();
    return sha256(password + username + PEPPER_1).toString();
  }

  secondHash(token) {
    const username = this.getUsername();
    return sha256(token + username + PEPPER_2).toString();
  }

  doubleHash(password) {
    return this.secondHash(this.firstHash(password));
  }

  passwordMatchesSession(password) {
    return password && this.firstHash(password) === this.getToken();
  }

  setCredentials({ username, password }) {
    sessionStorage.setItem('username', _.trim(username));
    this.setPassword(password);
  }

  setPassword(password) {
    sessionStorage.setItem('token', this.firstHash(password));
  }

  sessionExists() {
    return !!(this.getUsername() && this.getToken());
  }

  getUsername() {
    return sessionStorage.getItem('username');
  }

  getToken() {
    return sessionStorage.getItem('token');
  }

  logout() {
    sessionStorage.clear();
  }

  getHashedToken() {
    const token = this.getToken();
    if (!token) {
      return null;
    }

    return this.secondHash(token);
  }

  encryptWithToken(text) {
    const token = this.getToken();
    return AES.encrypt(text, token).toString();
  }

  encryptWithUpdatedPassword(text, updatedPassword) {
    const token = this.firstHash(updatedPassword);
    return AES.encrypt(text, token).toString();
  }

  decryptWithToken(text) {
    const token = this.getToken();
    return AES.decrypt(text, token).toString(UTF_8);
  }

  getAuthedHeaders() {
    const username = this.getUsername();
    const encoded = btoa(`${username}:`);
    return { 'Authorization': `Basic ${encoded}` };
  }

  getRegisterHeaders() {
    const username = this.getUsername();
    const decryptToken = this.getHashedToken();
    const encoded = btoa(`${username}:${decryptToken}`);
    return { 'Authorization': `Basic ${encoded}` };
  }
}
