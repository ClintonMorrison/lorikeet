import axios from 'axios';

export default class APIService {
  constructor({ baseURL, authService }) {
    this.authService = authService;

    axios.defaults.baseURL = baseURL;
    // axios.defaults.headers.common['Authorization'] = this.getAuthHeader();
    axios.defaults.headers.common['Accept'] = 'application/json';
  }

  handleAuthError(error) {
    console.error(error);
    if (error?.response?.status === 401) {
      this.authService.logout();
      setTimeout(() => window.location.assign('/login'), 500);
    }
    return error;
  }

  get(path, params, headers) {
    return axios({
      method: 'get',
      url: `/${path}`,
      params,
      headers
    });
  }

  post(path, data, headers) {
    return axios({
      method: 'post',
      url: `/${path}`,
      data,
      headers
    });
  }

  put(path, data, headers) {
    return axios({
      method: 'put',
      url: `/${path}`,
      data,
      headers
    });
  }

  del(path, headers) {
    return axios({
      method: 'delete',
      url: `/${path}`,
      headers
    });
  }
}
