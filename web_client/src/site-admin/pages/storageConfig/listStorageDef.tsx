import { useQuery } from "@apollo/client";
import classNames from "classnames";
import React from "react";
import { gql } from "~src/__generated__";
import { HEADING_2 } from "~src/ui/classNames";
import { BlockLoader } from "~src/ui/loader";
import { DumbStorageDefTable } from "../../components/storageDefTable";
import { useNavigate } from "react-router-dom";

const listStorageDefQuery = gql(`
  query ListStorageDef {
    viewer {
      id
      storageDefinitions {
        ...StorageDefinitionFragment
      }
    }
  }
`);

export function ListStorageDef() {
  const navigate = useNavigate();
  const { data: storageDefs, loading } = useQuery(listStorageDefQuery, {
    fetchPolicy: "cache-and-network",
  });
  const onAddNew = React.useCallback(() => {
    navigate("/site-admin/storage/storage-def/new");
  }, [navigate]);
  const onEdit = React.useCallback(
    (id: string) => {
      navigate(`/site-admin/storage/storage-def/${id}`);
    },
    [navigate]
  );
  return (
    <div className="m-auto max-w-5xl">
      <h1 className={classNames(HEADING_2, "font-poppins")}>Storage Backend</h1>
      {!storageDefs && loading && <BlockLoader />}
      <div className="mt-4">
        <DumbStorageDefTable
          data={storageDefs?.viewer.storageDefinitions || []}
          onAddNew={onAddNew}
          onEdit={onEdit}
        />
      </div>
    </div>
  );
}
