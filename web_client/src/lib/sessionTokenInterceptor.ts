export const TOKEN_NAME = "sessionToken";

(function () {
  if (!window.SESSION_HEADER_NAME) {
    return;
  }
  const originalFetch = window.fetch;

  window.fetch = async function (input, init) {
    const response = await originalFetch(input, init);
    const sessionToken = response.headers.get(window.SESSION_HEADER_NAME);
    if (sessionToken) {
      localStorage.setItem(TOKEN_NAME, sessionToken);
    }
    return response;
  };
})();
