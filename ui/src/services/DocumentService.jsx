import _ from 'lodash';
import { downloadAsJSON, downloadAsCSV, downloadAsText } from "../utils/download";
import moment from "moment/moment";

const defaultEmptyDocument = { passwords: [] };

export default class AuthService {
  constructor({ apiService, authService }) {
    this.apiService = apiService;
    this.authService = authService;
  }

  createDocument({ username, password, recaptchaResult }) {
    this.authService.setCredentials({ username, password });

    return this.apiService.post("document", {
      document: '',
      password: this.authService.getServerToken({ password }),
      recaptchaResult,
    }, this.authService.getRegisterHeaders());
  }

  loadDocument() {
    return this.apiService.get("document", {}, this.authService.getAuthedHeaders()).then(resp => {
      const encryptedDocument = _.get(resp, "data.document") || '';

      const decryptedDocument = encryptedDocument ?
        this.authService.decrypt({ text: encryptedDocument }) :
        JSON.stringify(defaultEmptyDocument);

      this.document = JSON.parse(decryptedDocument);

      return this.document;
    }).catch(e => {
      this.apiService.handleAuthError(e);
    });
  }

  updateDocument({ document, password }) {
    const unencryptedDocument = JSON.stringify(document);

    const encryptedDocument = password ?
      this.authService.encrypt({ text: unencryptedDocument, password }) :
      this.authService.encrypt({ text: unencryptedDocument });

    return this.apiService.put("document", {
      document: encryptedDocument,
      password: password ? this.authService.getServerToken({ password }) : undefined
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
      created: moment().toISOString(),
      updated: moment().toISOString(),
      lastUsed: moment().toISOString()
    };
  }

  updatePassword(password) {
    return this.loadDocument().then(document => {
      return this.updateDocument({ document, password });
    }).then((resp) => {
      this.authService.setCredentials({ password })
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

    this.loadDocument().then(document => handler(document, `passwords.${extension}`));
  }
}