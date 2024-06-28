/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type CreateUserWithOrganizationInput = {
  organizationName: Scalars['String'];
  userEmail: Scalars['String'];
  userPassword: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  authenticate: ViewerResult;
  createStorageDefinition?: Maybe<StorageDefinition>;
  createUserWithOrganization: ViewerResult;
  logout: ViewerResult;
  updateStorageDefinition?: Maybe<StorageDefinition>;
};


export type MutationAuthenticateArgs = {
  email: Scalars['String'];
  organizationId?: InputMaybe<Scalars['String']>;
  password: Scalars['String'];
};


export type MutationCreateStorageDefinitionArgs = {
  input: CreateStorageDefinitionInput;
};


export type MutationCreateUserWithOrganizationArgs = {
  input: CreateUserWithOrganizationInput;
};


export type MutationUpdateStorageDefinitionArgs = {
  input: UpdateStorageDefinitionInput;
};

export type Organization = {
  __typename?: 'Organization';
  id: Scalars['ID'];
  name: Scalars['String'];
  slug: Scalars['String'];
};

export type OrganizationUser = {
  __typename?: 'OrganizationUser';
  id: Scalars['ID'];
  organization: Organization;
  roles: Array<Role>;
  user: User;
};

export type OtherStorageConfig = {
  __typename?: 'OtherStorageConfig';
  _empty?: Maybe<Scalars['String']>;
};

export enum PermissionNameEnum {
  AdminAccess = 'AdminAccess',
  SiteOwnerAccess = 'SiteOwnerAccess'
}

export type Query = {
  __typename?: 'Query';
  viewer: Viewer;
};

export type Role = {
  __typename?: 'Role';
  key: Scalars['String'];
  name: Scalars['String'];
};

export type S3StorageConfig = {
  __typename?: 'S3StorageConfig';
  access: Scalars['String'];
  bucket: Scalars['String'];
  endpoint: Scalars['String'];
  secret: Scalars['String'];
};

export type StorageConfig = OtherStorageConfig | S3StorageConfig;

export type StorageDefinition = {
  __typename?: 'StorageDefinition';
  config: StorageConfig;
  id: Scalars['ID'];
  identifier: Scalars['String'];
  isEnabled: Scalars['Boolean'];
  priority: Scalars['Int'];
};

export enum StorageTypeEnum {
  Other = 'Other',
  S3 = 'S3'
}

export type User = {
  __typename?: 'User';
  email: Scalars['String'];
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type Viewer = {
  __typename?: 'Viewer';
  hasPermission: Scalars['Boolean'];
  id: Scalars['ID'];
  organizationUser?: Maybe<OrganizationUser>;
  storageDefinitions: Array<StorageDefinition>;
};


export type ViewerHasPermissionArgs = {
  permission: PermissionNameEnum;
};

export type ViewerResult = {
  __typename?: 'ViewerResult';
  viewer: Viewer;
};

export type CreateStorageDefinitionInput = {
  configJSON: Scalars['String'];
  identifier: Scalars['String'];
  isEnabled: Scalars['Boolean'];
  priority: Scalars['Int'];
  storageType: StorageTypeEnum;
};

export type UpdateStorageDefinitionInput = {
  configJSON?: InputMaybe<Scalars['String']>;
  identifier: Scalars['String'];
  isEnabled?: InputMaybe<Scalars['Boolean']>;
  priority?: InputMaybe<Scalars['Int']>;
};

export type CreateUserWithOrganizationMutationVariables = Exact<{
  input: CreateUserWithOrganizationInput;
}>;


export type CreateUserWithOrganizationMutation = { __typename?: 'Mutation', createUserWithOrganization: { __typename?: 'ViewerResult', viewer: { __typename?: 'Viewer', id: string, organizationUser?: { __typename?: 'OrganizationUser', id: string, user: { __typename?: 'User', id: string, email: string, name: string } } | null } } };

export type AuthenticateMutationVariables = Exact<{
  email: Scalars['String'];
  password: Scalars['String'];
}>;


export type AuthenticateMutation = { __typename?: 'Mutation', authenticate: { __typename?: 'ViewerResult', viewer: { __typename?: 'Viewer', id: string, organizationUser?: { __typename?: 'OrganizationUser', id: string, user: { __typename?: 'User', id: string, email: string, name: string } } | null } } };

export type AuthQueryVariables = Exact<{ [key: string]: never; }>;


export type AuthQuery = { __typename?: 'Query', viewer: { __typename?: 'Viewer', id: string, hasAdminAccess: boolean, hasSiteOwnerAccess: boolean, organizationUser?: { __typename?: 'OrganizationUser', id: string, user: { __typename?: 'User', id: string, email: string, name: string } } | null } };

export type LogoutMutationVariables = Exact<{ [key: string]: never; }>;


export type LogoutMutation = { __typename?: 'Mutation', logout: { __typename?: 'ViewerResult', viewer: { __typename?: 'Viewer', id: string, hasAdminAccess: boolean, hasSiteOwnerAccess: boolean, organizationUser?: { __typename?: 'OrganizationUser', id: string, user: { __typename?: 'User', id: string, email: string, name: string } } | null } } };

export type ListStorageDefQueryVariables = Exact<{ [key: string]: never; }>;


export type ListStorageDefQuery = { __typename?: 'Query', viewer: { __typename?: 'Viewer', storageDefinitions: Array<{ __typename: 'StorageDefinition', id: string, identifier: string, isEnabled: boolean, priority: number, config: { __typename?: 'OtherStorageConfig' } | { __typename?: 'S3StorageConfig', bucket: string, endpoint: string, access: string, secret: string } }> } };

export type StorageDefinitionFragmentFragment = { __typename: 'StorageDefinition', id: string, identifier: string, isEnabled: boolean, priority: number, config: { __typename?: 'OtherStorageConfig' } | { __typename?: 'S3StorageConfig', bucket: string, endpoint: string, access: string, secret: string } };

export const StorageDefinitionFragmentFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"StorageDefinitionFragment"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"StorageDefinition"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"identifier"}},{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"isEnabled"}},{"kind":"Field","name":{"kind":"Name","value":"priority"}},{"kind":"Field","name":{"kind":"Name","value":"config"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"S3StorageConfig"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"bucket"}},{"kind":"Field","name":{"kind":"Name","value":"endpoint"}},{"kind":"Field","name":{"kind":"Name","value":"access"}},{"kind":"Field","name":{"kind":"Name","value":"secret"}}]}}]}}]}}]} as unknown as DocumentNode<StorageDefinitionFragmentFragment, unknown>;
export const CreateUserWithOrganizationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"createUserWithOrganization"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateUserWithOrganizationInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createUserWithOrganization"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"organizationUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<CreateUserWithOrganizationMutation, CreateUserWithOrganizationMutationVariables>;
export const AuthenticateDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"authenticate"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"email"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"password"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"authenticate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"email"},"value":{"kind":"Variable","name":{"kind":"Name","value":"email"}}},{"kind":"Argument","name":{"kind":"Name","value":"password"},"value":{"kind":"Variable","name":{"kind":"Name","value":"password"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"organizationUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<AuthenticateMutation, AuthenticateMutationVariables>;
export const AuthDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Auth"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"organizationUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"Field","alias":{"kind":"Name","value":"hasAdminAccess"},"name":{"kind":"Name","value":"hasPermission"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"permission"},"value":{"kind":"EnumValue","value":"AdminAccess"}}]},{"kind":"Field","alias":{"kind":"Name","value":"hasSiteOwnerAccess"},"name":{"kind":"Name","value":"hasPermission"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"permission"},"value":{"kind":"EnumValue","value":"SiteOwnerAccess"}}]}]}}]}}]} as unknown as DocumentNode<AuthQuery, AuthQueryVariables>;
export const LogoutDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Logout"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"logout"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"organizationUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"Field","alias":{"kind":"Name","value":"hasAdminAccess"},"name":{"kind":"Name","value":"hasPermission"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"permission"},"value":{"kind":"EnumValue","value":"AdminAccess"}}]},{"kind":"Field","alias":{"kind":"Name","value":"hasSiteOwnerAccess"},"name":{"kind":"Name","value":"hasPermission"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"permission"},"value":{"kind":"EnumValue","value":"SiteOwnerAccess"}}]}]}}]}}]}}]} as unknown as DocumentNode<LogoutMutation, LogoutMutationVariables>;
export const ListStorageDefDocument = {"kind":"Document", "definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ListStorageDef"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"storageDefinitions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"StorageDefinitionFragment"}}]}}]}}]}},...StorageDefinitionFragmentFragmentDoc.definitions]} as unknown as DocumentNode<ListStorageDefQuery, ListStorageDefQueryVariables>;