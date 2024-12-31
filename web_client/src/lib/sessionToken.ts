import { TOKEN_NAME } from "./sessionTokenInterceptor";

export function getSessionToken() {
  return localStorage.getItem(TOKEN_NAME);
}

export function removeSessionToken() {
  localStorage.removeItem(TOKEN_NAME);
}
