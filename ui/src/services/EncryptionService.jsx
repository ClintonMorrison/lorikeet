import sha256 from 'crypto-js/sha256';
import AES from 'crypto-js/aes';
import UTF_8 from 'crypto-js/enc-utf8';

const PEPPER_1 = 'CC352C99A14616AD22678563ECDA5';
const PEPPER_2 = '7767B9225CF66B418DD2A39CBC4AA';


export default class EncryptionService {
  generateClientEncryptToken({ username, password, salt }) {
    return sha256(password + username + salt).toString();
  }

  // Either password or clientEncryptToken is required
  generateServerEncryptToken({ username, password, token }) {
    let t1 = token;
    if (!t1) {
      t1 = sha256(password + username + PEPPER_1).toString();
    }

    return sha256(t1 + username + PEPPER_2).toString();
  }

  encrypt({ text, secret }) {
    return AES.encrypt(text, secret).toString();
  }

  decrypt({ text, secret }) {
    return AES.decrypt(text, secret).toString(UTF_8);
  }
}
