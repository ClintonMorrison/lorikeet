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

  passwordMatchesSession(password) {
    const username = this.getUsername();
    const generatedClientTokenV1 = this.getClientToken({ username, password });

    return password && generatedClientTokenV1 === this.getClientToken({ version: 1 });
  }

  setCredentials({ username, password }) {
    if (username) {
      this.storageHelper.setUsername(_.trim(username));
    }

    if (password) {
      const tokenV1 = this.getClientToken({
        username: username || this.getUsername(),
        password,
        version: 1,
      });
      this.storageHelper.setClientTokenV1(tokenV1);
    }
  }

  setSalt({ password, salt }) {
    if (salt) {
      const tokenV2 = this.getClientToken({
        username: this.getUsername(),
        password,
        salt,
        version: 2,
      });
      this.storageHelper.setClientTokenV2(tokenV2);
    }
  }

  sessionExists() {
    return !!(this.getUsername() && this.getClientToken({ version: 1 }));
  }

  getUsername() {
    return this.storageHelper.getUsername();
  }

  logout() {
    sessionStorage.clear();
  }

  getClientToken({ username, password, version, salt }) {
    if (version === 2) {
      if (username && password && salt) {
        return this.encryptionService.generateClientEncryptTokenV2({ username, password, salt });
      }

      return this.storageHelper.getClientTokenV2();
    }

    if (username && password) {
      return this.encryptionService.generateClientEncryptTokenV1({ username, password })
    }

    return this.storageHelper.getClientTokenV1();
  }

  getServerToken({ password } = {}) {
    const username = this.getUsername();
    if (!username) {
      return null;
    }

    if (password) {
      return this.encryptionService.generateServerEncryptTokenV1({ username, password });
    }

    const token = this.getClientToken({ version: 1 });
    if (!token) {
      return null;
    }

    return this.encryptionService.generateServerEncryptTokenV1({ username, token });
  }

  encrypt({ text, password, salt, version }) {
    const username = this.getUsername();

    const secret = password ?
      this.getClientToken({ username, password, version, salt }) :
      this.getClientToken({ version });

    return this.encryptionService.encrypt({ text, secret });
  }

  decrypt({ text, version }) {
    const secret = this.getClientToken({ version });
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
