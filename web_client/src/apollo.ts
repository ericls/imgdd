import { ApolloClient, createHttpLink, InMemoryCache } from "@apollo/client";
import { getSessionToken } from "./lib/sessionToken";

const apolloLinkWithToken = createHttpLink({
  uri: (operation) => {
    const context = operation.getContext();
    const maybeCaptchaToken = context.captchaToken;
    let urlStr = `/query?operation=${operation.operationName}`;
    if (maybeCaptchaToken) {
      urlStr += `&captchaToken=${maybeCaptchaToken}`;
    }
    return urlStr;
  },
  fetch: async (uri, options) => {
    const token = getSessionToken();
    return fetch(uri, {
      ...options,
      headers: {
        ...(options || {}).headers,
        ...(token && window.SESSION_HEADER_NAME
          ? { [window.SESSION_HEADER_NAME]: token }
          : {}),
      },
    });
  },
});

export const apolloClient = new ApolloClient({
  cache: new InMemoryCache(),
  link: apolloLinkWithToken,
});
