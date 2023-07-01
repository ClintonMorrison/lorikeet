import _ from 'lodash';
import { downloadAsJSON, downloadAsCSV, downloadAsText } from "../utils/download";
import moment from "moment/moment";

const defaultEmptyDocument = { passwords: [] };

export default class AuthService {
  constructor({ apiService, authService }) {
    this.apiService = apiService;
    this.authService = authService;
    this.storageVersion = 0;
    this.clientEncryptVersion = 0;
    this.salt = null;
    this.document = null;
  }

  createDocument({ username, password, recaptchaResult }) {
    this.storageVersion = 2; // use new version for all new users
    this.clientEncryptVersion = 2;
    this.authService.setCredentials({ username, password });

    return this.apiService.post("document", {
      document: '',
      password: this.authService.getServerToken({ password }),
      recaptchaResult,
    }, this.authService.getRegisterHeaders()).then(resp => {
      this.salt = _.get(resp, 'data.salt') || '';
      this.authService.setCredentials({ password, salt: this.salt });
    });
  }

  loadDocument() {
    return this.apiService.get("document", {}, this.authService.getAuthedHeaders()).then(resp => {
      const encryptedDocument = _.get(resp, "data.document") || '';
      this.storageVersion = _.get(resp, 'data.storageVersion') || 1;
      this.clientEncryptVersion = _.get(resp, 'data.clientEncryptVersion') || 1;
      this.salt = _.get(resp, 'data.salt') || '';

      const decryptedDocument = encryptedDocument ?
        this.authService.decrypt({ text: encryptedDocument, version: this.clientEncryptVersion }) :
        JSON.stringify(defaultEmptyDocument);

      this.document = JSON.parse(decryptedDocument);
      return {
        document: this.document,
        salt: this.salt,
        version: this.clientEncryptVersion,
        needsMigration: this.storageVersion !== 2 || this.clientEncryptVersion !== 2
      };
    }).catch(e => {
      this.apiService.handleAuthError(e);
    });
  }

  updateDocument({ document, password, migrate }) {
    const unencryptedDocument = JSON.stringify(document);
    const version = this.clientEncryptVersion;
    const salt = this.salt;

    const encryptedDocument = password ?
      this.authService.encrypt({ text: unencryptedDocument, salt, password, version }) :
      this.authService.encrypt({ text: unencryptedDocument, version });

    return this.apiService.put("document", {
      document: encryptedDocument,
      password: password ? this.authService.getServerToken({ password }) : undefined,
      migrate,
    }, this.authService.getAuthedHeaders());
  }

  // Updates document from v1 to v2
  async migrateDocument({ document, password }) {
    const response = await this.updateDocument({ document, migrate: true });
    this.storageVersion = response.data.storageVersion;
    this.clientEncryptVersion = response.data.clientEncryptVersion;
    this.salt = response.data.salt;
    this.authService.setCredentials({ password, salt: this.salt, });

    if (this.storageVersion !== 2 || this.clientEncryptVersion !== 2) {
      throw new Error('Version not updated properly, aborting migration')
    }

    return await this.updateDocument({ document, password });
  }

  deleteDocument() {
    return this.apiService.del("document", this.authService.getAuthedHeaders())
      .catch(e => {
        this.apiService.handleAuthError(e);
      });
  }

  createPassword(id) {
    return {
      id,
      title: '',
      username: '',
      password: '',
      email: '',
      website: '',
      notes: '',
      created: moment().toISOString(),
      updated: moment().toISOString(),
      lastUsed: moment().toISOString()
    };
  }

  updatePassword(password) {
    return this.loadDocument().then(({ document }) => {
      return this.updateDocument({ document, password });
    }).then((resp) => {
      this.authService.setCredentials({ password, salt: resp.data.salt })
    }).catch(e => {
      this.apiService.handleAuthError(e);
    });
    ;
  }

  downloadDocument(type) {
    let handler = downloadAsText;
    if (type === 'json') handler = downloadAsJSON;
    if (type === 'csv') handler = downloadAsCSV;

    const extension = type === 'text' ? 'txt' : type;

    this.loadDocument().then(({ document }) => handler(document, `passwords.${extension}`));
  }
}