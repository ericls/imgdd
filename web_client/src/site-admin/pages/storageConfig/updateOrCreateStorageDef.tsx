import React from "react";
import { useParams } from "react-router-dom";
import { StorageConfigForm } from "~src/site-admin/components/storageConfigForm";

export function UpdateOrCreateStorageDef() {
  const { id: maybeId } = useParams<{ id: string }>();
  const id = maybeId === "new" ? undefined : maybeId;
  return (
    <div>
      <h1>UpdateOrCreateStorageDef</h1>
      <StorageConfigForm id={id} />
    </div>
  );
}
