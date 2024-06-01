import { ApolloClient, InMemoryCache } from '@apollo/client';


export const apolloClient = new ApolloClient({
  uri: (operation) => {
    return `/query?operation=${operation.operationName}`;
  },
  cache: new InMemoryCache(),
});