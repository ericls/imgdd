import { useLazyQuery } from "@apollo/client";
import classNames from "classnames";
import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "react-toastify";
import { gql } from "~src/__generated__";
import { StorageConfigForm } from "~src/site-admin/components/storageConfigForm";
import { StorageType } from "~src/site-admin/types";
import { HEADING_1 } from "~src/ui/classNames";

const getStorageDefQuery = gql(/* GraphQL */ `
  query GetStorageDef($id: ID!) {
    viewer {
      id
      getStorageDefinition(id: $id) {
        ...StorageDefinitionFragment
      }
    }
  }
`);

export function UpdateOrCreateStorageDef() {
  const { id: maybeId } = useParams<{ id: string }>();
  const id = maybeId === "new" ? undefined : maybeId;
  const [execute, { data, loading, called }] = useLazyQuery(
    getStorageDefQuery,
    {
      fetchPolicy: "network-only",
    }
  );
  React.useEffect(() => {
    if (id) {
      execute({ variables: { id } });
    }
  }, [id, execute]);
  const formInitialValue = React.useMemo(() => {
    if (!id || !called || loading || !data?.viewer.getStorageDefinition) {
      return undefined;
    }
    const storageDef = data.viewer.getStorageDefinition;
    const providerConfig = storageDef.config;
    const storageType: StorageType =
      providerConfig.__typename === "S3StorageConfig" ? "S3" : "__other";
    return {
      storageType,
      identifier: storageDef.identifier,
      providerConfig:
        providerConfig.__typename === "S3StorageConfig"
          ? {
              bucket: providerConfig.bucket,
              endpoint: providerConfig.endpoint,
              access: providerConfig.access,
              secret: providerConfig.secret,
            }
          : { __other: "" },
    };
  }, [id, called, loading, data]);
  const navigate = useNavigate();
  const afterSave = React.useCallback(
    (id?: string) => {
      if (id) {
        toast.success("Storage definition saved");
        setTimeout(() => {
          navigate("/site-admin/storage/storage-def/list");
        }, 300);
      } else {
        toast.error("Failed to save storage definition");
      }
    },
    [navigate]
  );
  return (
    <div className="m-auto max-w-5xl">
      <h1 className={classNames(HEADING_1)}>
        {id ? "Update" : "Create"} Storage Definition
      </h1>
      <div>
        <StorageConfigForm
          id={id}
          initialValue={formInitialValue}
          key={`${id}-${called}-${loading}-${data?.viewer.getStorageDefinition}`}
          afterSave={afterSave}
        />
      </div>
    </div>
  );
}
