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
  Time: any;
};

export type CreateUserWithOrganizationInput = {
  organizationName: Scalars['String'];
  userEmail: Scalars['String'];
  userPassword: Scalars['String'];
};

export type DeleteImageInput = {
  id: Scalars['ID'];
};

export type DeleteImageResult = {
  __typename?: 'DeleteImageResult';
  id?: Maybe<Scalars['ID']>;
};

export type FsStorageConfig = {
  __typename?: 'FSStorageConfig';
  mediaRoot: Scalars['String'];
};

export type Image = {
  __typename?: 'Image';
  MIMEType: Scalars['String'];
  createdAt: Scalars['Time'];
  id: Scalars['ID'];
  identifier: Scalars['String'];
  name: Scalars['String'];
  nominalByteSize: Scalars['Int'];
  nominalHeight: Scalars['Int'];
  nominalWidth: Scalars['Int'];
  revisions: Array<Image>;
  root?: Maybe<Image>;
  storedImages: Array<StoredImage>;
  url: Scalars['String'];
};

export type ImageEdge = {
  __typename?: 'ImageEdge';
  cursor: Scalars['String'];
  node: Image;
};

export type ImageFilterInput = {
  createdAtGt?: InputMaybe<Scalars['Time']>;
  createdAtLt?: InputMaybe<Scalars['Time']>;
  createdBy?: InputMaybe<Scalars['ID']>;
  nameContains?: InputMaybe<Scalars['String']>;
};

export type ImageOrderByInput = {
  createdAt?: InputMaybe<PaginationDirection>;
  id?: InputMaybe<PaginationDirection>;
  name?: InputMaybe<PaginationDirection>;
};

export type ImagePageInfo = {
  __typename?: 'ImagePageInfo';
  currentCount?: Maybe<Scalars['Int']>;
  endCursor?: Maybe<Scalars['String']>;
  hasNextPage: Scalars['Boolean'];
  hasPreviousPage: Scalars['Boolean'];
  startCursor?: Maybe<Scalars['String']>;
  totalCount?: Maybe<Scalars['Int']>;
};

export type ImagesResult = {
  __typename?: 'ImagesResult';
  edges: Array<ImageEdge>;
  pageInfo: ImagePageInfo;
};

export type Mutation = {
  __typename?: 'Mutation';
  authenticate: ViewerResult;
  checkStorageDefinitionConnectivity?: Maybe<StorageDefinitionConnectivityResult>;
  createStorageDefinition?: Maybe<StorageDefinition>;
  createUserWithOrganization: ViewerResult;
  deleteImage: DeleteImageResult;
  logout: ViewerResult;
  resetPassword: ResetPasswordResult;
  sendResetPasswordEmail: SendResetPasswordEmailResult;
  updateStorageDefinition?: Maybe<StorageDefinition>;
};


export type MutationAuthenticateArgs = {
  email: Scalars['String'];
  organizationId?: InputMaybe<Scalars['String']>;
  password: Scalars['String'];
};


export type MutationCheckStorageDefinitionConnectivityArgs = {
  input: CheckStorageDefinitionConnectivityInput;
};


export type MutationCreateStorageDefinitionArgs = {
  input: CreateStorageDefinitionInput;
};


export type MutationCreateUserWithOrganizationArgs = {
  input: CreateUserWithOrganizationInput;
};


export type MutationDeleteImageArgs = {
  input: DeleteImageInput;
};


export type MutationResetPasswordArgs = {
  input: ResetPasswordInput;
};


export type MutationSendResetPasswordEmailArgs = {
  input: SendResetPasswordEmailInput;
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

export type PageInfo = {
  __typename?: 'PageInfo';
  hasNextPage: Scalars['Boolean'];
  hasPreviousPage: Scalars['Boolean'];
  totalCount: Scalars['Int'];
};

export type PaginatedUsers = {
  __typename?: 'PaginatedUsers';
  nodes: Array<User>;
  pageInfo: PageInfo;
};

export enum PaginationDirection {
  Asc = 'asc',
  Desc = 'desc'
}

export enum PermissionNameEnum {
  AdminAccess = 'AdminAccess',
  SiteOwnerAccess = 'SiteOwnerAccess'
}

export type Query = {
  __typename?: 'Query';
  viewer: Viewer;
};

export type ResetPasswordInput = {
  message: Scalars['String'];
  password: Scalars['String'];
};

export type ResetPasswordResult = {
  __typename?: 'ResetPasswordResult';
  success: Scalars['Boolean'];
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

export type SendResetPasswordEmailInput = {
  email: Scalars['String'];
};

export type SendResetPasswordEmailResult = {
  __typename?: 'SendResetPasswordEmailResult';
  success: Scalars['Boolean'];
};

export type StorageConfig = FsStorageConfig | OtherStorageConfig | S3StorageConfig | WebDavStorageConfig;

export type StorageDefinition = {
  __typename?: 'StorageDefinition';
  config: StorageConfig;
  connectivity: Scalars['Boolean'];
  id: Scalars['ID'];
  identifier: Scalars['String'];
  isEnabled: Scalars['Boolean'];
  priority: Scalars['Int'];
};

export type StorageDefinitionConnectivityResult = {
  __typename?: 'StorageDefinitionConnectivityResult';
  error?: Maybe<Scalars['String']>;
  ok: Scalars['Boolean'];
};

export enum StorageTypeEnum {
  Fs = 'FS',
  Other = 'Other',
  S3 = 'S3',
  WebDav = 'WebDAV'
}

export type StoredImage = {
  __typename?: 'StoredImage';
  fileIdentifier: Scalars['String'];
  id: Scalars['ID'];
  storageDefinition: StorageDefinition;
};

export type User = {
  __typename?: 'User';
  email: Scalars['String'];
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type Viewer = {
  __typename?: 'Viewer';
  allUsers: Array<User>;
  getStorageDefinition?: Maybe<StorageDefinition>;
  hasPermission: Scalars['Boolean'];
  id: Scalars['ID'];
  images: ImagesResult;
  organizationUser?: Maybe<OrganizationUser>;
  paginatedAllUsers: PaginatedUsers;
  storageDefinitions: Array<StorageDefinition>;
};


export type ViewerAllUsersArgs = {
  limit?: InputMaybe<Scalars['Int']>;
  offset?: InputMaybe<Scalars['Int']>;
  search?: InputMaybe<Scalars['String']>;
};


export type ViewerGetStorageDefinitionArgs = {
  id: Scalars['ID'];
};


export type ViewerHasPermissionArgs = {
  permission: PermissionNameEnum;
};


export type ViewerImagesArgs = {
  after?: InputMaybe<Scalars['String']>;
  before?: InputMaybe<Scalars['String']>;
  filters?: InputMaybe<ImageFilterInput>;
  orderBy?: InputMaybe<ImageOrderByInput>;
};


export type ViewerPaginatedAllUsersArgs = {
  limit?: InputMaybe<Scalars['Int']>;
  offset?: InputMaybe<Scalars['Int']>;
  search?: InputMaybe<Scalars['String']>;
};

export type ViewerResult = {
  __typename?: 'ViewerResult';
  viewer: Viewer;
};

export type WebDavStorageConfig = {
  __typename?: 'WebDAVStorageConfig';
  password: Scalars['String'];
  pathPrefix: Scalars['String'];
  url: Scalars['String'];
  username: Scalars['String'];
};

export type CheckStorageDefinitionConnectivityInput = {
  id: Scalars['ID'];
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

export type SendResetPasswordEmailMutationVariables = Exact<{
  input: SendResetPasswordEmailInput;
}>;


export type SendResetPasswordEmailMutation = { __typename?: 'Mutation', sendResetPasswordEmail: { __typename?: 'SendResetPasswordEmailResult', success: boolean } };

export type ResetPasswordMutationVariables = Exact<{
  input: ResetPasswordInput;
}>;


export type ResetPasswordMutation = { __typename?: 'Mutation', resetPassword: { __typename?: 'ResetPasswordResult', success: boolean } };

export type ImagesQueryQueryVariables = Exact<{
  orderBy?: InputMaybe<ImageOrderByInput>;
  filters?: InputMaybe<ImageFilterInput>;
  after?: InputMaybe<Scalars['String']>;
  before?: InputMaybe<Scalars['String']>;
}>;


export type ImagesQueryQuery = { __typename?: 'Query', viewer: { __typename?: 'Viewer', id: string, images: { __typename?: 'ImagesResult', pageInfo: { __typename?: 'ImagePageInfo', hasNextPage: boolean, hasPreviousPage: boolean, startCursor?: string | null, endCursor?: string | null, totalCount?: number | null, currentCount?: number | null }, edges: Array<{ __typename?: 'ImageEdge', cursor: string, node: { __typename?: 'Image', id: string, url: string, name: string, nominalWidth: number, nominalHeight: number, nominalByteSize: number, createdAt: any, storedImages: Array<{ __typename?: 'StoredImage', id: string }> } }> } } };

export type DeleteImageMutationVariables = Exact<{
  input: DeleteImageInput;
}>;


export type DeleteImageMutation = { __typename?: 'Mutation', deleteImage: { __typename?: 'DeleteImageResult', id?: string | null } };

export type AuthQueryVariables = Exact<{ [key: string]: never; }>;


export type AuthQuery = { __typename?: 'Query', viewer: { __typename?: 'Viewer', id: string, hasAdminAccess: boolean, hasSiteOwnerAccess: boolean, organizationUser?: { __typename?: 'OrganizationUser', id: string, user: { __typename?: 'User', id: string, email: string, name: string } } | null } };

export type LogoutMutationVariables = Exact<{ [key: string]: never; }>;


export type LogoutMutation = { __typename?: 'Mutation', logout: { __typename?: 'ViewerResult', viewer: { __typename?: 'Viewer', id: string, hasAdminAccess: boolean, hasSiteOwnerAccess: boolean, organizationUser?: { __typename?: 'OrganizationUser', id: string, user: { __typename?: 'User', id: string, email: string, name: string } } | null } } };

export type StorageDefTableConnectivityCellMutationMutationVariables = Exact<{
  input: CheckStorageDefinitionConnectivityInput;
}>;


export type StorageDefTableConnectivityCellMutationMutation = { __typename?: 'Mutation', checkStorageDefinitionConnectivity?: { __typename?: 'StorageDefinitionConnectivityResult', ok: boolean, error?: string | null } | null };

export type ListStorageDefQueryVariables = Exact<{ [key: string]: never; }>;


export type ListStorageDefQuery = { __typename?: 'Query', viewer: { __typename?: 'Viewer', id: string, storageDefinitions: Array<{ __typename: 'StorageDefinition', id: string, identifier: string, isEnabled: boolean, priority: number, config: { __typename?: 'FSStorageConfig', mediaRoot: string } | { __typename?: 'OtherStorageConfig' } | { __typename?: 'S3StorageConfig', bucket: string, endpoint: string, access: string, secret: string } | { __typename?: 'WebDAVStorageConfig', url: string, username: string, password: string, pathPrefix: string } }> } };

export type GetStorageDefQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type GetStorageDefQuery = { __typename?: 'Query', viewer: { __typename?: 'Viewer', id: string, getStorageDefinition?: { __typename: 'StorageDefinition', id: string, identifier: string, isEnabled: boolean, priority: number, config: { __typename?: 'FSStorageConfig', mediaRoot: string } | { __typename?: 'OtherStorageConfig' } | { __typename?: 'S3StorageConfig', bucket: string, endpoint: string, access: string, secret: string } | { __typename?: 'WebDAVStorageConfig', url: string, username: string, password: string, pathPrefix: string } } | null } };

export type ListUsersQueryVariables = Exact<{
  limit?: InputMaybe<Scalars['Int']>;
  offset?: InputMaybe<Scalars['Int']>;
  search?: InputMaybe<Scalars['String']>;
}>;


export type ListUsersQuery = { __typename?: 'Query', viewer: { __typename?: 'Viewer', id: string, paginatedAllUsers: { __typename?: 'PaginatedUsers', nodes: Array<{ __typename?: 'User', id: string, email: string }>, pageInfo: { __typename?: 'PageInfo', totalCount: number, hasNextPage: boolean, hasPreviousPage: boolean } } } };

export type StorageDefinitionFragmentFragment = { __typename: 'StorageDefinition', id: string, identifier: string, isEnabled: boolean, priority: number, config: { __typename?: 'FSStorageConfig', mediaRoot: string } | { __typename?: 'OtherStorageConfig' } | { __typename?: 'S3StorageConfig', bucket: string, endpoint: string, access: string, secret: string } | { __typename?: 'WebDAVStorageConfig', url: string, username: string, password: string, pathPrefix: string } };

export type CreateStorageDefMutationVariables = Exact<{
  input: CreateStorageDefinitionInput;
}>;


export type CreateStorageDefMutation = { __typename?: 'Mutation', createStorageDefinition?: { __typename: 'StorageDefinition', id: string, identifier: string, isEnabled: boolean, priority: number, config: { __typename?: 'FSStorageConfig', mediaRoot: string } | { __typename?: 'OtherStorageConfig' } | { __typename?: 'S3StorageConfig', bucket: string, endpoint: string, access: string, secret: string } | { __typename?: 'WebDAVStorageConfig', url: string, username: string, password: string, pathPrefix: string } } | null };

export type UpdateStorageDefMutationVariables = Exact<{
  input: UpdateStorageDefinitionInput;
}>;


export type UpdateStorageDefMutation = { __typename?: 'Mutation', updateStorageDefinition?: { __typename: 'StorageDefinition', id: string, identifier: string, isEnabled: boolean, priority: number, config: { __typename?: 'FSStorageConfig', mediaRoot: string } | { __typename?: 'OtherStorageConfig' } | { __typename?: 'S3StorageConfig', bucket: string, endpoint: string, access: string, secret: string } | { __typename?: 'WebDAVStorageConfig', url: string, username: string, password: string, pathPrefix: string } } | null };

export const StorageDefinitionFragmentFragmentDoc = {"kind":"Document","definitions":[{"kind":"FragmentDefinition","name":{"kind":"Name","value":"StorageDefinitionFragment"},"typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"StorageDefinition"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"identifier"}},{"kind":"Field","name":{"kind":"Name","value":"__typename"}},{"kind":"Field","name":{"kind":"Name","value":"isEnabled"}},{"kind":"Field","name":{"kind":"Name","value":"priority"}},{"kind":"Field","name":{"kind":"Name","value":"config"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"S3StorageConfig"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"bucket"}},{"kind":"Field","name":{"kind":"Name","value":"endpoint"}},{"kind":"Field","name":{"kind":"Name","value":"access"}},{"kind":"Field","name":{"kind":"Name","value":"secret"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"FSStorageConfig"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"mediaRoot"}}]}},{"kind":"InlineFragment","typeCondition":{"kind":"NamedType","name":{"kind":"Name","value":"WebDAVStorageConfig"}},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"username"}},{"kind":"Field","name":{"kind":"Name","value":"password"}},{"kind":"Field","name":{"kind":"Name","value":"pathPrefix"}}]}}]}}]}}]} as unknown as DocumentNode<StorageDefinitionFragmentFragment, unknown>;
export const CreateUserWithOrganizationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"createUserWithOrganization"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CreateUserWithOrganizationInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createUserWithOrganization"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"organizationUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<CreateUserWithOrganizationMutation, CreateUserWithOrganizationMutationVariables>;
export const AuthenticateDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"authenticate"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"email"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"password"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"authenticate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"email"},"value":{"kind":"Variable","name":{"kind":"Name","value":"email"}}},{"kind":"Argument","name":{"kind":"Name","value":"password"},"value":{"kind":"Variable","name":{"kind":"Name","value":"password"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"organizationUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<AuthenticateMutation, AuthenticateMutationVariables>;
export const SendResetPasswordEmailDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"sendResetPasswordEmail"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"SendResetPasswordEmailInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"sendResetPasswordEmail"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}}]}}]}}]} as unknown as DocumentNode<SendResetPasswordEmailMutation, SendResetPasswordEmailMutationVariables>;
export const ResetPasswordDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"resetPassword"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ResetPasswordInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"resetPassword"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"success"}}]}}]}}]} as unknown as DocumentNode<ResetPasswordMutation, ResetPasswordMutationVariables>;
export const ImagesQueryDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ImagesQuery"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"orderBy"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"ImageOrderByInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"filters"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"ImageFilterInput"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"after"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"before"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"images"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"orderBy"},"value":{"kind":"Variable","name":{"kind":"Name","value":"orderBy"}}},{"kind":"Argument","name":{"kind":"Name","value":"filters"},"value":{"kind":"Variable","name":{"kind":"Name","value":"filters"}}},{"kind":"Argument","name":{"kind":"Name","value":"after"},"value":{"kind":"Variable","name":{"kind":"Name","value":"after"}}},{"kind":"Argument","name":{"kind":"Name","value":"before"},"value":{"kind":"Variable","name":{"kind":"Name","value":"before"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"pageInfo"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"hasNextPage"}},{"kind":"Field","name":{"kind":"Name","value":"hasPreviousPage"}},{"kind":"Field","name":{"kind":"Name","value":"startCursor"}},{"kind":"Field","name":{"kind":"Name","value":"endCursor"}},{"kind":"Field","name":{"kind":"Name","value":"totalCount"}},{"kind":"Field","name":{"kind":"Name","value":"currentCount"}}]}},{"kind":"Field","name":{"kind":"Name","value":"edges"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"cursor"}},{"kind":"Field","name":{"kind":"Name","value":"node"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"url"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"nominalWidth"}},{"kind":"Field","name":{"kind":"Name","value":"nominalHeight"}},{"kind":"Field","name":{"kind":"Name","value":"nominalByteSize"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"storedImages"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]}}]}}]}}]}}]} as unknown as DocumentNode<ImagesQueryQuery, ImagesQueryQueryVariables>;
export const DeleteImageDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteImage"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"DeleteImageInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"deleteImage"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<DeleteImageMutation, DeleteImageMutationVariables>;
export const AuthDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Auth"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"organizationUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"Field","alias":{"kind":"Name","value":"hasAdminAccess"},"name":{"kind":"Name","value":"hasPermission"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"permission"},"value":{"kind":"EnumValue","value":"AdminAccess"}}]},{"kind":"Field","alias":{"kind":"Name","value":"hasSiteOwnerAccess"},"name":{"kind":"Name","value":"hasPermission"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"permission"},"value":{"kind":"EnumValue","value":"SiteOwnerAccess"}}]}]}}]}}]} as unknown as DocumentNode<AuthQuery, AuthQueryVariables>;
export const LogoutDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"Logout"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"logout"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"organizationUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"user"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}},{"kind":"Field","alias":{"kind":"Name","value":"hasAdminAccess"},"name":{"kind":"Name","value":"hasPermission"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"permission"},"value":{"kind":"EnumValue","value":"AdminAccess"}}]},{"kind":"Field","alias":{"kind":"Name","value":"hasSiteOwnerAccess"},"name":{"kind":"Name","value":"hasPermission"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"permission"},"value":{"kind":"EnumValue","value":"SiteOwnerAccess"}}]}]}}]}}]}}]} as unknown as DocumentNode<LogoutMutation, LogoutMutationVariables>;
export const StorageDefTableConnectivityCellMutationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"StorageDefTableConnectivityCellMutation"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"checkStorageDefinitionConnectivityInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"checkStorageDefinitionConnectivity"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"ok"}},{"kind":"Field","name":{"kind":"Name","value":"error"}}]}}]}}]} as unknown as DocumentNode<StorageDefTableConnectivityCellMutationMutation, StorageDefTableConnectivityCellMutationMutationVariables>;
export const ListStorageDefDocument = {"kind":"Document", "definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ListStorageDef"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"storageDefinitions"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"StorageDefinitionFragment"}}]}}]}}]}},...StorageDefinitionFragmentFragmentDoc.definitions]} as unknown as DocumentNode<ListStorageDefQuery, ListStorageDefQueryVariables>;
export const GetStorageDefDocument = {"kind":"Document", "definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"GetStorageDef"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"getStorageDefinition"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"StorageDefinitionFragment"}}]}}]}}]}},...StorageDefinitionFragmentFragmentDoc.definitions]} as unknown as DocumentNode<GetStorageDefQuery, GetStorageDefQueryVariables>;
export const ListUsersDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"ListUsers"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"limit"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"offset"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"search"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"viewer"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"paginatedAllUsers"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"limit"},"value":{"kind":"Variable","name":{"kind":"Name","value":"limit"}}},{"kind":"Argument","name":{"kind":"Name","value":"offset"},"value":{"kind":"Variable","name":{"kind":"Name","value":"offset"}}},{"kind":"Argument","name":{"kind":"Name","value":"search"},"value":{"kind":"Variable","name":{"kind":"Name","value":"search"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"nodes"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"email"}}]}},{"kind":"Field","name":{"kind":"Name","value":"pageInfo"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"totalCount"}},{"kind":"Field","name":{"kind":"Name","value":"hasNextPage"}},{"kind":"Field","name":{"kind":"Name","value":"hasPreviousPage"}}]}}]}}]}}]}}]} as unknown as DocumentNode<ListUsersQuery, ListUsersQueryVariables>;
export const CreateStorageDefDocument = {"kind":"Document", "definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateStorageDef"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"createStorageDefinitionInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createStorageDefinition"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"StorageDefinitionFragment"}}]}}]}},...StorageDefinitionFragmentFragmentDoc.definitions]} as unknown as DocumentNode<CreateStorageDefMutation, CreateStorageDefMutationVariables>;
export const UpdateStorageDefDocument = {"kind":"Document", "definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateStorageDef"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"updateStorageDefinitionInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateStorageDefinition"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"FragmentSpread","name":{"kind":"Name","value":"StorageDefinitionFragment"}}]}}]}},...StorageDefinitionFragmentFragmentDoc.definitions]} as unknown as DocumentNode<UpdateStorageDefMutation, UpdateStorageDefMutationVariables>;