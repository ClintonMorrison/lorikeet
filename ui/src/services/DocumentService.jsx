import _ from 'lodash';
import { downloadAsJSON, downloadAsCSV, downloadAsText } from "../utils/download";
import moment from "moment/moment";

export default class AuthService {
  constructor({ apiService, authService }) {
    this.apiService = apiService;
    this.authService = authService;
  }

  createDocument({ username, password }) {
    this.authService.setCredentials({ username, password });

    const initialDocument = JSON.stringify({ passwords: [] });
    const encryptedDocument = this.authService.encryptWithToken(initialDocument);

    return this.apiService.post("document", {
      document: encryptedDocument,
      password: this.authService.doubleHash(password)
    }, this.authService.getHeaders());
  }

  loadDocument() {
    return this.apiService.get("document", {}, this.authService.getHeaders()).then(resp => {
      const encryptedDocument = _.get(resp, "data.document") || '{}';
      const decryptedDocument = this.authService.decryptWithToken(encryptedDocument);
      this.document = JSON.parse(decryptedDocument);
      return this.document;
    });
  }

  updateDocument({ document, password }) {
    const unencryptedDocument = JSON.stringify(document);

    const encryptedDocument = password ?
      this.authService.encryptWithToken(unencryptedDocument, this.authService.firstHash(password)) :
      this.authService.encryptWithToken(unencryptedDocument);

    return this.apiService.put("document", {
      document: encryptedDocument,
      password: password ? this.authService.doubleHash(password) : undefined
    }, this.authService.getHeaders());
  }

  deleteDocument() {
    return this.apiService.del("document", this.authService.getHeaders());
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
    }).then(() => this.authService.setPassword(password));
  }

  downloadDocument(type) {
    let handler = downloadAsText;
    if (type === 'json') handler = downloadAsJSON;
    if (type === 'csv') handler = downloadAsCSV;

    const extension = type === 'text' ? 'txt' : type;

    this.loadDocument().then(document => handler(document, `passwords.${extension}`));
  }
}