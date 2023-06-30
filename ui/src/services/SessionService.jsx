export default class AuthService {
  constructor({ apiService, authService }) {
    this.apiService = apiService;
    this.authService = authService;
  }

  createSession({ recaptchaResult }) {
    const username = this.authService.getUsername();
    const decryptToken = this.authService.getServerToken();

    return this.apiService.post("session", {
      username,
      decryptToken,
      recaptchaResult
    }, this.authService.getAuthedHeaders());
  }

  deleteSession() {
    return this.apiService.del("session", this.authService.getAuthedHeaders());
  }
}