import { useMutation, useQuery } from "@apollo/client";
import { noop } from "lodash-es";
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
    hasAdminAccess: hasPermission(permission: AdminAccess)
    hasSiteOwnerAccess: hasPermission(permission: SiteOwnerAccess)
  }
}
`);

const LOGOUT_MUTATION = gql(`
mutation Logout {
  logout {
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
    hasAdminAccess: hasPermission(permission: AdminAccess)
    hasSiteOwnerAccess: hasPermission(permission: SiteOwnerAccess)
  }
  }
}
`);

const EMPTY_AUTH_QUERY_RESULT: AuthQuery = {
  viewer: {
    id: "viewer",
    organizationUser: null,
    hasAdminAccess: false,
    hasSiteOwnerAccess: false,
  },
};

const AuthContext = React.createContext<{
  data: AuthQuery | null | undefined;
  isLoading: boolean;
  logout: () => void;
}>({ data: EMPTY_AUTH_QUERY_RESULT, isLoading: true, logout: noop });

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { data, loading } = useQuery(AUTH_QUERY);
  const [logout] = useMutation(LOGOUT_MUTATION);
  const value = { data, isLoading: loading, logout };
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return React.useContext(AuthContext);
}
