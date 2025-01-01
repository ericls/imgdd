import { TOKEN_NAME } from "./sessionTokenInterceptor";

export function getSessionToken() {
  return localStorage.getItem(TOKEN_NAME);
}

export function removeSessionToken() {
  localStorage.removeItem(TOKEN_NAME);
}

export function addSessionHeaderToXMLHttpRequest(request: XMLHttpRequest) {
  if (!window.SESSION_HEADER_NAME) {
    return;
  }
  const token = getSessionToken();
  if (token) {
    request.setRequestHeader(window.SESSION_HEADER_NAME, token);
  }
}
