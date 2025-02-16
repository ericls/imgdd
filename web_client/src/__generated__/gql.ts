/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel-plugin for production.
 */
const documents = {
    "\nmutation createUserWithOrganization($input: CreateUserWithOrganizationInput!) {\n  createUserWithOrganization(\n    input: $input\n  ) {\n    viewer {\n      id\n      organizationUser {\n        id\n        user {\n          id\n          email\n          name\n        }\n      }\n    }\n  }\n}\n": types.CreateUserWithOrganizationDocument,
    "\nmutation authenticate($email: String!, $password: String!) {\n  authenticate(email: $email, password: $password) {\n    viewer {\n      id\n      organizationUser {\n        id\n        user {\n          id\n          email\n          name\n        }\n      }\n    }\n  }\n}\n": types.AuthenticateDocument,
    "\nmutation sendResetPasswordEmail($input: SendResetPasswordEmailInput!) {\n  sendResetPasswordEmail(input: $input) {\n    success\n  }\n}\n": types.SendResetPasswordEmailDocument,
    "\nmutation resetPassword($input: ResetPasswordInput!) {\n  resetPassword(input: $input) {\n    success\n  }\n}\n": types.ResetPasswordDocument,
    "\n  query ImagesQuery(\n    $orderBy: ImageOrderByInput\n    $filters: ImageFilterInput\n    $after: String\n    $before: String\n  ) {\n    viewer {\n      id\n      images(\n        orderBy: $orderBy\n        filters: $filters\n        after: $after\n        before: $before\n      ) {\n        pageInfo {\n          hasNextPage\n          hasPreviousPage\n          startCursor\n          endCursor\n          totalCount\n          currentCount\n        }\n        edges {\n          cursor\n          node {\n            id\n            url\n            name\n            nominalWidth\n            nominalHeight\n            nominalByteSize\n            createdAt\n            storedImages {\n              id\n            }\n          }\n        }\n      }\n    }\n  }\n": types.ImagesQueryDocument,
    "\n  mutation DeleteImage($input: DeleteImageInput!) {\n    deleteImage(input: $input) {\n      id\n    }\n  }\n": types.DeleteImageDocument,
    "\nquery Auth {\n  viewer {\n    id\n    organizationUser {\n      id\n      user {\n        id\n        email\n        name\n      }\n    }\n    hasAdminAccess: hasPermission(permission: AdminAccess)\n    hasSiteOwnerAccess: hasPermission(permission: SiteOwnerAccess)\n  }\n}\n": types.AuthDocument,
    "\nmutation Logout {\n  logout {\n    viewer {\n    id\n    organizationUser {\n      id\n      user {\n        id\n        email\n        name\n      }\n    }\n    hasAdminAccess: hasPermission(permission: AdminAccess)\n    hasSiteOwnerAccess: hasPermission(permission: SiteOwnerAccess)\n  }\n  }\n}\n": types.LogoutDocument,
    "\nmutation StorageDefTableConnectivityCellMutation(\n    $input: checkStorageDefinitionConnectivityInput!\n  ) {\n    checkStorageDefinitionConnectivity(input: $input) {\n      ok\n      error\n    }\n  }\n": types.StorageDefTableConnectivityCellMutationDocument,
    "\n  query ListStorageDef {\n    viewer {\n      id\n      storageDefinitions {\n        ...StorageDefinitionFragment\n      }\n    }\n  }\n": types.ListStorageDefDocument,
    "\n  query GetStorageDef($id: ID!) {\n    viewer {\n      id\n      getStorageDefinition(id: $id) {\n        ...StorageDefinitionFragment\n      }\n    }\n  }\n": types.GetStorageDefDocument,
    "\n  fragment StorageDefinitionFragment on StorageDefinition {\n    id\n    identifier\n    __typename\n    isEnabled\n    priority\n    config {\n      ... on S3StorageConfig {\n        bucket\n        endpoint\n        access\n        secret\n      }\n      ... on FSStorageConfig {\n        mediaRoot\n      }\n      ... on WebDAVStorageConfig {\n        url\n        username\n        password\n        pathPrefix\n      }\n    }\n  }\n": types.StorageDefinitionFragmentFragmentDoc,
    "\n  mutation CreateStorageDef($input: createStorageDefinitionInput!) {\n    createStorageDefinition(input: $input) {\n      ...StorageDefinitionFragment\n    }\n  }\n": types.CreateStorageDefDocument,
    "\n  mutation UpdateStorageDef($input: updateStorageDefinitionInput!) {\n    updateStorageDefinition(input: $input) {\n      ...StorageDefinitionFragment\n    }\n  }\n": types.UpdateStorageDefDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation createUserWithOrganization($input: CreateUserWithOrganizationInput!) {\n  createUserWithOrganization(\n    input: $input\n  ) {\n    viewer {\n      id\n      organizationUser {\n        id\n        user {\n          id\n          email\n          name\n        }\n      }\n    }\n  }\n}\n"): (typeof documents)["\nmutation createUserWithOrganization($input: CreateUserWithOrganizationInput!) {\n  createUserWithOrganization(\n    input: $input\n  ) {\n    viewer {\n      id\n      organizationUser {\n        id\n        user {\n          id\n          email\n          name\n        }\n      }\n    }\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation authenticate($email: String!, $password: String!) {\n  authenticate(email: $email, password: $password) {\n    viewer {\n      id\n      organizationUser {\n        id\n        user {\n          id\n          email\n          name\n        }\n      }\n    }\n  }\n}\n"): (typeof documents)["\nmutation authenticate($email: String!, $password: String!) {\n  authenticate(email: $email, password: $password) {\n    viewer {\n      id\n      organizationUser {\n        id\n        user {\n          id\n          email\n          name\n        }\n      }\n    }\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation sendResetPasswordEmail($input: SendResetPasswordEmailInput!) {\n  sendResetPasswordEmail(input: $input) {\n    success\n  }\n}\n"): (typeof documents)["\nmutation sendResetPasswordEmail($input: SendResetPasswordEmailInput!) {\n  sendResetPasswordEmail(input: $input) {\n    success\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation resetPassword($input: ResetPasswordInput!) {\n  resetPassword(input: $input) {\n    success\n  }\n}\n"): (typeof documents)["\nmutation resetPassword($input: ResetPasswordInput!) {\n  resetPassword(input: $input) {\n    success\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query ImagesQuery(\n    $orderBy: ImageOrderByInput\n    $filters: ImageFilterInput\n    $after: String\n    $before: String\n  ) {\n    viewer {\n      id\n      images(\n        orderBy: $orderBy\n        filters: $filters\n        after: $after\n        before: $before\n      ) {\n        pageInfo {\n          hasNextPage\n          hasPreviousPage\n          startCursor\n          endCursor\n          totalCount\n          currentCount\n        }\n        edges {\n          cursor\n          node {\n            id\n            url\n            name\n            nominalWidth\n            nominalHeight\n            nominalByteSize\n            createdAt\n            storedImages {\n              id\n            }\n          }\n        }\n      }\n    }\n  }\n"): (typeof documents)["\n  query ImagesQuery(\n    $orderBy: ImageOrderByInput\n    $filters: ImageFilterInput\n    $after: String\n    $before: String\n  ) {\n    viewer {\n      id\n      images(\n        orderBy: $orderBy\n        filters: $filters\n        after: $after\n        before: $before\n      ) {\n        pageInfo {\n          hasNextPage\n          hasPreviousPage\n          startCursor\n          endCursor\n          totalCount\n          currentCount\n        }\n        edges {\n          cursor\n          node {\n            id\n            url\n            name\n            nominalWidth\n            nominalHeight\n            nominalByteSize\n            createdAt\n            storedImages {\n              id\n            }\n          }\n        }\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation DeleteImage($input: DeleteImageInput!) {\n    deleteImage(input: $input) {\n      id\n    }\n  }\n"): (typeof documents)["\n  mutation DeleteImage($input: DeleteImageInput!) {\n    deleteImage(input: $input) {\n      id\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nquery Auth {\n  viewer {\n    id\n    organizationUser {\n      id\n      user {\n        id\n        email\n        name\n      }\n    }\n    hasAdminAccess: hasPermission(permission: AdminAccess)\n    hasSiteOwnerAccess: hasPermission(permission: SiteOwnerAccess)\n  }\n}\n"): (typeof documents)["\nquery Auth {\n  viewer {\n    id\n    organizationUser {\n      id\n      user {\n        id\n        email\n        name\n      }\n    }\n    hasAdminAccess: hasPermission(permission: AdminAccess)\n    hasSiteOwnerAccess: hasPermission(permission: SiteOwnerAccess)\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation Logout {\n  logout {\n    viewer {\n    id\n    organizationUser {\n      id\n      user {\n        id\n        email\n        name\n      }\n    }\n    hasAdminAccess: hasPermission(permission: AdminAccess)\n    hasSiteOwnerAccess: hasPermission(permission: SiteOwnerAccess)\n  }\n  }\n}\n"): (typeof documents)["\nmutation Logout {\n  logout {\n    viewer {\n    id\n    organizationUser {\n      id\n      user {\n        id\n        email\n        name\n      }\n    }\n    hasAdminAccess: hasPermission(permission: AdminAccess)\n    hasSiteOwnerAccess: hasPermission(permission: SiteOwnerAccess)\n  }\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation StorageDefTableConnectivityCellMutation(\n    $input: checkStorageDefinitionConnectivityInput!\n  ) {\n    checkStorageDefinitionConnectivity(input: $input) {\n      ok\n      error\n    }\n  }\n"): (typeof documents)["\nmutation StorageDefTableConnectivityCellMutation(\n    $input: checkStorageDefinitionConnectivityInput!\n  ) {\n    checkStorageDefinitionConnectivity(input: $input) {\n      ok\n      error\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query ListStorageDef {\n    viewer {\n      id\n      storageDefinitions {\n        ...StorageDefinitionFragment\n      }\n    }\n  }\n"): (typeof documents)["\n  query ListStorageDef {\n    viewer {\n      id\n      storageDefinitions {\n        ...StorageDefinitionFragment\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  query GetStorageDef($id: ID!) {\n    viewer {\n      id\n      getStorageDefinition(id: $id) {\n        ...StorageDefinitionFragment\n      }\n    }\n  }\n"): (typeof documents)["\n  query GetStorageDef($id: ID!) {\n    viewer {\n      id\n      getStorageDefinition(id: $id) {\n        ...StorageDefinitionFragment\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  fragment StorageDefinitionFragment on StorageDefinition {\n    id\n    identifier\n    __typename\n    isEnabled\n    priority\n    config {\n      ... on S3StorageConfig {\n        bucket\n        endpoint\n        access\n        secret\n      }\n      ... on FSStorageConfig {\n        mediaRoot\n      }\n      ... on WebDAVStorageConfig {\n        url\n        username\n        password\n        pathPrefix\n      }\n    }\n  }\n"): (typeof documents)["\n  fragment StorageDefinitionFragment on StorageDefinition {\n    id\n    identifier\n    __typename\n    isEnabled\n    priority\n    config {\n      ... on S3StorageConfig {\n        bucket\n        endpoint\n        access\n        secret\n      }\n      ... on FSStorageConfig {\n        mediaRoot\n      }\n      ... on WebDAVStorageConfig {\n        url\n        username\n        password\n        pathPrefix\n      }\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation CreateStorageDef($input: createStorageDefinitionInput!) {\n    createStorageDefinition(input: $input) {\n      ...StorageDefinitionFragment\n    }\n  }\n"): (typeof documents)["\n  mutation CreateStorageDef($input: createStorageDefinitionInput!) {\n    createStorageDefinition(input: $input) {\n      ...StorageDefinitionFragment\n    }\n  }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n  mutation UpdateStorageDef($input: updateStorageDefinitionInput!) {\n    updateStorageDefinition(input: $input) {\n      ...StorageDefinitionFragment\n    }\n  }\n"): (typeof documents)["\n  mutation UpdateStorageDef($input: updateStorageDefinitionInput!) {\n    updateStorageDefinition(input: $input) {\n      ...StorageDefinitionFragment\n    }\n  }\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;