import { ApolloClient, createHttpLink, InMemoryCache } from "@apollo/client";
import { getSessionToken } from "./lib/sessionToken";

const apolloLinkWithToken = createHttpLink({
  uri: (operation) => {
    return `/query?operation=${operation.operationName}`;
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
