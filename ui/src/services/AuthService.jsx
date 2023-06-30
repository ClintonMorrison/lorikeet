import _ from 'lodash';

class StorageHelper {
  getUsername() {
    return sessionStorage.getItem('username');
  }
  setUsername(value) {
    sessionStorage.setItem('username', value);
  }

  getClientTokenV1() {
    return sessionStorage.getItem('token');
  }
  setClientTokenV1(value) {
    return sessionStorage.setItem('token', value);
  }
}

export default class AuthService {
  constructor({ encryptionService }) {
    this.encryptionService = encryptionService;
    this.storageHelper = new StorageHelper();
  }

  passwordMatchesSession(password) {
    const username = this.getUsername();
    const generatedClientTokenV1 = this.encryptionService.generateClientEncryptTokenV1({ username, password });

    return password && generatedClientTokenV1 === this.getClientToken();
  }

  setCredentials({ username, password }) {
    if (username) {
      this.storageHelper.setUsername(_.trim(username));
    }

    if (password) {
      const tokenV1 = this.encryptionService.generateClientEncryptTokenV1({
        username: username || this.getUsername(),
        password
      });
      this.storageHelper.setClientTokenV1(tokenV1);
    }
  }

  sessionExists() {
    return !!(this.getUsername() && this.getClientToken());
  }

  getUsername() {
    return this.storageHelper.getUsername();
  }

  getClientToken() {
    return this.storageHelper.getClientTokenV1();
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

    const token = this.getClientToken();
    if (!token) {
      return null;
    }

    return this.encryptionService.generateServerEncryptTokenV1({ username, token });
  }

  encrypt({ text, password }) {
    const username = this.getUsername();

    const secret = password ?
      this.encryptionService.generateClientEncryptTokenV1({ username, password }) :
      this.getClientToken();

    return this.encryptionService.encrypt({ text, secret });
  }

  decrypt({ text }) {
    const secret = this.getClientToken();
    return this.encryptionService.decrypt({ text, secret });
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
