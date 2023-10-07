import _ from 'lodash';

class StorageHelper {
  getUsername() {
    return sessionStorage.getItem('username');
  }
  setUsername(value) {
    sessionStorage.setItem('username', value);
  }

  getClientTokenV2() {
    return sessionStorage.getItem('tokenV2');
  }
  setClientTokenV2(value) {
    return sessionStorage.setItem('tokenV2', value);
  }
}

export default class AuthService {
  constructor({ encryptionService }) {
    this.encryptionService = encryptionService;
    this.storageHelper = new StorageHelper();
  }

  passwordMatchesSession({ password, salt }) {
    const username = this.getUsername();
    const generatedClientTokenV2 = this.getClientToken({ username, password, salt });

    return password && generatedClientTokenV2 === this.getClientToken();
  }

  setUsername({ username }) {
    if (username) {
      this.storageHelper.setUsername(_.trim(username));
    }
  }

  setSalt({ password, salt }) {
    if (salt) {
      const tokenV2 = this.getClientToken({
        username: this.getUsername(),
        password,
        salt,
      });
      this.storageHelper.setClientTokenV2(tokenV2);
    }
  }

  sessionExists() {
    return !!(this.getUsername() && this.getClientToken());
  }

  getUsername() {
    return this.storageHelper.getUsername();
  }

  logout() {
    sessionStorage.clear();
  }

  getClientToken({ username, password, salt } = {}) {
    if (username && password && salt) {
      return this.encryptionService.generateClientEncryptTokenV2({ username, password, salt });
    }

    return this.storageHelper.getClientTokenV2();
  }

  getServerToken({ password } = {}) {
    const username = this.getUsername();
    if (!username || !password) {
      return null;
    }

    return this.encryptionService.generateServerEncryptTokenV1({ username, password });
  }

  encrypt({ text, password, salt }) {
    const username = this.getUsername();

    const secret = password ?
      this.getClientToken({ username, password, salt }) :
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

  getRegisterHeaders({ password }) {
    const username = this.getUsername();
    const decryptToken = this.getServerToken({ password });
    const encoded = btoa(`${username}:${decryptToken}`);
    return { 'Authorization': `Basic ${encoded}` };
  }
}
