import sha256 from 'crypto-js/sha256';
import AES from 'crypto-js/aes';
import UTF_8 from 'crypto-js/enc-utf8';
import _ from 'lodash';

const PEPPER_1 = 'CC352C99A14616AD22678563ECDA5';
const PEPPER_2 = '7767B9225CF66B418DD2A39CBC4AA';


export default class EncryptionService {
  // V1 encryption
  generateClientEncryptTokenV1({ username, password }) {
    return sha256(password + username + PEPPER_1).toString();
  }

  // Either password or clientEncryptToken is required
  generateServerEncryptTokenV1({ username, password, clientEncryptToken }) {
    let t1 = clientEncryptToken;
    if (!t1) {
      t1 = this.generateClientEncryptTokenV1({ username, password });
    }

    return sha256(t1 + username + PEPPER_2).toString();
  }

  // V2 encryption
}
