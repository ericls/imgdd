import React from "react";
import { useForm } from "react-hook-form";

type S3StorageConfigData = {
  bucket: string;
  endpoint: string;
  access: string;
  secret: string;
};

type StorageProviderConfigData = S3StorageConfigData | { __other: string };

type StorageConfigData = {
  storageType: "S3" | "__other";
  providerConfig: StorageProviderConfigData;
};

type StorageConfigFormProps = {
  id?: string;
  initialValue?: StorageConfigData;
};

export function StorageConfigForm({
  initialValue,
  id,
}: StorageConfigFormProps) {
  const [storageTypeValue, setStorageTypeValue] = React.useState<
    "S3" | "__other"
  >(initialValue?.storageType || "S3");
  const commonFieldsForm = useForm<Omit<StorageConfigData, "providerConfig">>({
    defaultValues: { storageType: initialValue?.storageType || "S3" },
  });
  const providerConfigForm = useForm<StorageProviderConfigData>({
    defaultValues: initialValue?.providerConfig,
  });
  return <div>Form</div>;
}
