import axios from 'axios';
export const CancelToken = axios.CancelToken;

export const auth = axios.create({
  baseURL: `/api/auth/`,
  withCredentials: true,
  validateStatus: (status) => status >= 200 && status < 500,
});

export const profile = axios.create({
  baseURL: `/api/profile/`,
  withCredentials: true,
  validateStatus: (status) => status >= 200 && status < 500,
});

export const passwd = axios.create({
  baseURL: `/api/passwd/`,
  withCredentials: true,
  validateStatus: (status) => status >= 200 && status < 500,
});

export const email = axios.create({
  baseURL: `/api/email/`,
  withCredentials: true,
  validateStatus: (status) => status >= 200 && status < 500,
});

export const search = axios.create({
  baseURL: `/api/search/`,
  withCredentials: true,
  validateStatus: (status) => status >= 200 && status < 500,
});
export const movies = axios.create({
  baseURL: `/api/movies/`,
  withCredentials: true,
  validateStatus: (status) => status >= 200 && status < 500,
});
export const test = axios.create({
  baseURL: `/api/test/`,
  withCredentials: true,
  validateStatus: (status) => status >= 200 && status < 500,
});

export const main = axios.create({
  baseURL: `/api/`,
  withCredentials: true,
  validateStatus: (status) => status >= 200 && status < 500,
});

// GET /api/info

// GET /api/auth/basic
// GET /api/auth/oauth42

// GET /api/profile/get (в зависимости от того есть ли параметр user_id  в урле или нет возвращает твоего или чужого юзера)
// PUT /api/profile/create
// PATCH /api/profile/patch
// DELETE /api/profile/delete

// PATCH /api/email/patch
// POST /api/email/confirm
// POST /api/email/resend

// PATCH /api/passwd/patch
// POST /api/passwd/repair
