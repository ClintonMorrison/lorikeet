import AES from 'crypto-js/aes';
import UTF_8 from 'crypto-js/enc-utf8';
import _ from 'lodash';

export default class AuthService {
  constructor({ encryptionService }) {
    this.encryptionService = encryptionService;
  }

  firstHash(password) {
    const username = this.getUsername();
    return this.encryptionService.generateClientEncryptTokenV1({ username, password });
  }

  passwordMatchesSession(password) {
    return password && this.firstHash(password) === this.getToken();
  }

  setCredentials({ username, password }) {
    if (username) {
      sessionStorage.setItem('username', _.trim(username));
    }

    if (password) {
      sessionStorage.setItem('token', this.firstHash(password));
    }
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

  getServerToken({ password } = {}) {
    const username = this.getUsername();
    if (!username) {
      return null;
    }

    if (password) {
      return this.encryptionService.generateServerEncryptTokenV1({ username, password });
    }

    const token = this.getToken();
    if (!token) {
      return null;
    }

    return this.encryptionService.generateServerEncryptTokenV1({ username, token });
  }

  encrypt({ text, password }) {
    const token = password ?
      this.firstHash(password) : // TODO
      this.getToken();

    return AES.encrypt(text, token).toString();
  }

  decrypt({ text }) {
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
    const decryptToken = this.getServerToken();
    const encoded = btoa(`${username}:${decryptToken}`);
    return { 'Authorization': `Basic ${encoded}` };
  }
}
