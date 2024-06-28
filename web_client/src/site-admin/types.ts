import { gql } from "~src/__generated__";

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

type storageType = "S3" | "__";
