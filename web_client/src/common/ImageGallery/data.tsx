import React from "react";
import { useLazyQuery, useMutation } from "@apollo/client";
import { gql } from "~src/__generated__";
import {
  ImageOrderByInput,
  ImagesQueryQuery,
  ImagesQueryQueryVariables,
  PaginationDirection,
} from "~src/__generated__/graphql";

const ImagesQueryDoc = gql(`
  query ImagesQuery(
    $orderBy: ImageOrderByInput
    $filters: ImageFilterInput
    $after: String
    $before: String
  ) {
    viewer {
      id
      images(
        orderBy: $orderBy
        filters: $filters
        after: $after
        before: $before
      ) {
        pageInfo {
          hasNextPage
          hasPreviousPage
          startCursor
          endCursor
          totalCount
          currentCount
        }
        edges {
          cursor
          node {
            id
            url
            name
            nominalWidth
            nominalHeight
            nominalByteSize
            createdAt
            storedImages {
              id
            }
          }
        }
      }
    }
  }
`);

export type UseImagesQueryOptions = {
  createdById?: string;
  nameContains?: string;
};

export function useImagesQuery({
  nameContains,
  createdById,
}: UseImagesQueryOptions) {
  const [orderBy, setOrderBy] = React.useState<ImageOrderByInput>({
    createdAt: PaginationDirection.Desc,
  });
  const [after, setAfter] = React.useState<string | null>(null);
  const [before, setBefore] = React.useState<string | null>(null);
  const variables: ImagesQueryQueryVariables = React.useMemo(() => {
    return {
      orderBy,
      filters: {
        nameContains,
        createdBy: createdById,
      },
      after,
      before,
    };
  }, [after, before, createdById, nameContains, orderBy]);
  const [
    execute,
    { data, loading, previousData, error, refetch, variables: dataVariables },
  ] = useLazyQuery(ImagesQueryDoc, {
    variables: variables,
    fetchPolicy: "network-only",
  });
  const currentPageInfo = React.useMemo<
    ImagesQueryQuery["viewer"]["images"]["pageInfo"]
  >(() => {
    return (
      data?.viewer.images.pageInfo ?? {
        hasNextPage: false,
        hasPreviousPage: false,
        startCursor: null,
        endCursor: null,
        totalCount: 0,
        currentCount: 0,
      }
    );
  }, [data]);
  const goNext = React.useCallback(() => {
    if (currentPageInfo.hasNextPage) {
      setBefore(null);
      setAfter(currentPageInfo.endCursor ?? null);
    }
  }, [currentPageInfo]);
  const goPrev = React.useCallback(() => {
    if (currentPageInfo.hasPreviousPage) {
      setAfter(null);
      setBefore(currentPageInfo.startCursor ?? null);
    }
  }, [currentPageInfo]);
  const hasPrev = currentPageInfo.hasPreviousPage;
  const hasNext = currentPageInfo.hasNextPage;
  if (
    !hasPrev &&
    !hasNext &&
    data?.viewer.images.edges.length === 0 &&
    !loading
  ) {
    if (variables.after || variables.before) {
      refetch({ ...variables, after: null, before: null });
    }
  }
  return {
    execute,
    data: data || previousData,
    loading,
    error,
    refetch,
    goNext,
    goPrev,
    hasPrev,
    hasNext,
    setOrderBy,
    dataVariables,
  };
}

const DeleteImageDoc = gql(`
  mutation DeleteImage($input: DeleteImageInput!) {
    deleteImage(input: $input) {
      id
    }
  }
`);

export function useDeleteImage(imageId: string) {
  const [execute, { loading, error, data }] = useMutation(DeleteImageDoc, {
    variables: {
      input: {
        id: imageId,
      },
    },
    refetchQueries: [ImagesQueryDoc],
    errorPolicy: "all",
  });
  return {
    execute,
    loading,
    error,
    data,
  };
}
