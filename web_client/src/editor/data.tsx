import { useMutation } from "@apollo/client/react";
import { gql } from "~src/__generated__";

const ApplyWatermarkDoc = gql(`
  mutation ApplyWatermark($input: ApplyWatermarkInput!) {
    applyWatermark(input: $input) {
      image {
        id
        url
        name
        identifier
        nominalWidth
        nominalHeight
        nominalByteSize
        MIMEType
        parent {
          id
          name
        }
        changes
      }
    }
  }
`);

export function useApplyWatermark() {
  const [execute, { loading, error, data }] = useMutation(ApplyWatermarkDoc);
  return { execute, loading, error, data };
}
