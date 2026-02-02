import _ from 'lodash';
import { downloadAsJSON, downloadAsCSV, downloadAsText } from "../utils/download";
import dayjs from 'dayjs';

const defaultEmptyDocument = { passwords: [] };

export default class AuthService {
  constructor({ apiService, authService }) {
    this.apiService = apiService;
    this.authService = authService;
    this.salt = null;
    this.document = null;
  }

  createDocument({ username, password, recaptchaResult }) {
    this.authService.setUsername({ username });

    const headers = this.authService.getRegisterHeaders({ password })

    return this.apiService.post("document", {
      document: '',
      password: this.authService.getServerToken({ password }),
      recaptchaResult,
    }, headers).then(resp => {
      this.salt = _.get(resp, 'data.salt') || '';
      this.authService.setSalt({ password, salt: this.salt });
    });
  }

  loadDocument() {
    return this.apiService.get("document", {}, this.authService.getAuthedHeaders()).then(resp => {
      const encryptedDocument = _.get(resp, "data.document") || '';
      this.salt = _.get(resp, 'data.salt') || '';

      const decryptedDocument = encryptedDocument ?
        this.authService.decrypt({ text: encryptedDocument }) :
        JSON.stringify(defaultEmptyDocument);

      this.document = JSON.parse(decryptedDocument);
      return {
        document: this.document,
        salt: this.salt,
      };
    }).catch(e => {
      this.apiService.handleAuthError(e);
    });
  }

  updateDocument({ document, password }) {
    const unencryptedDocument = JSON.stringify(document);
    const salt = this.salt;

    const encryptedDocument = password ?
      this.authService.encrypt({ text: unencryptedDocument, salt, password }) :
      this.authService.encrypt({ text: unencryptedDocument });

    return this.apiService.put("document", {
      document: encryptedDocument,
      password: password ? this.authService.getServerToken({ password }) : undefined,
    }, this.authService.getAuthedHeaders());
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
      created: dayjs().toISOString(),
      updated: dayjs().toISOString(),
      lastUsed: dayjs().toISOString()
    };
  }

  updatePassword(password) {
    return this.loadDocument().then(({ document }) => {
      return this.updateDocument({ document, password });
    }).then((resp) => {
      this.authService.setSalt({ password, salt: resp.data.salt })
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