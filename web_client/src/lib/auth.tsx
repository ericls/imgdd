import { useQuery } from "@apollo/client";
import React from "react";
import { gql } from "~src/__generated__/gql";
import { AuthQuery } from "~src/__generated__/graphql";

const AUTH_QUERY = gql(`
query Auth {
  viewer {
    id
    organizationUser {
      id
      user {
        id
        email
        name
      }
    }
  }
}
`);

const EMPTY_AUTH_QUERY_RESULT: AuthQuery = {
  viewer: {
    id: "viewer",
    organizationUser: null,
  },
};

const AuthContext = React.createContext<{
  data: AuthQuery | null | undefined;
  isLoading: boolean;
}>({ data: EMPTY_AUTH_QUERY_RESULT, isLoading: true });

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { data, loading } = useQuery(AUTH_QUERY);
  const value = { data, isLoading: loading };
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return React.useContext(AuthContext);
}
