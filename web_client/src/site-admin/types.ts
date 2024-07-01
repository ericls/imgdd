import { gql } from "~src/__generated__";

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const storageDefinitionFragment = gql(/* GraphQL */ `
  fragment StorageDefinitionFragment on StorageDefinition {
    id
    identifier
    __typename
    isEnabled
    priority
    config {
      ... on S3StorageConfig {
        bucket
        endpoint
        access
        secret
      }
    }
  }
`);

export const createStorageDefMutation = gql(/* GraphQL */ `
  mutation CreateStorageDef($input: createStorageDefinitionInput!) {
    createStorageDefinition(input: $input) {
      ...StorageDefinitionFragment
    }
  }
`);

export const updateStorageDefMutation = gql(/* GraphQL */ `
  mutation UpdateStorageDef($input: updateStorageDefinitionInput!) {
    updateStorageDefinition(input: $input) {
      ...StorageDefinitionFragment
    }
  }
`);

export type StorageType = "S3" | "__other";
