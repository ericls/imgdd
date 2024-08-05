import { useMutation } from "@apollo/client";
import React from "react";
import { useForm } from "react-hook-form";
import { InputWithLabel } from "~src/ui/input";
import { SelectWithLabel } from "~src/ui/select";
import { createStorageDefMutation, updateStorageDefMutation } from "../types";
import { StorageTypeEnum } from "~src/__generated__/graphql";
import { Button } from "~src/ui/button";

type S3StorageConfigData = {
  bucket: string;
  endpoint: string;
  access: string;
  secret: string;
};

type StorageProviderConfigData = S3StorageConfigData | { __other: string };

type StorageConfigData = {
  storageType: "S3" | "__other";
  identifier: string;
  priority: number;
  isEnabled: boolean;
  providerConfig: StorageProviderConfigData;
};

type StorageConfigFormProps = {
  id?: string;
  initialValue?: StorageConfigData;
  afterSave?: (id?: string) => void;
};

function S3ProviderConfigForm({
  form,
}: {
  form: ReturnType<typeof useForm<StorageProviderConfigData>>;
}) {
  return (
    <>
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label="Bucket"
        {...form.register("bucket", { required: true })}
      />
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label="Endpoint"
        {...form.register("endpoint", { required: true })}
      />
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label="Access"
        {...form.register("access", { required: true })}
      />
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label="Secret"
        type="password"
        {...form.register("secret", {
          required: true,
        })}
      />
    </>
  );
}

export function StorageConfigForm({
  initialValue,
  id,
  afterSave,
}: StorageConfigFormProps) {
  const commonFieldsForm = useForm<Omit<StorageConfigData, "providerConfig">>({
    defaultValues: {
      storageType: initialValue?.storageType || "S3",
      identifier: initialValue?.identifier || "",
      priority: initialValue?.priority || 1,
      isEnabled: initialValue?.isEnabled ?? true,
    },
  });
  const storageTypeValue = commonFieldsForm.watch("storageType");
  const providerConfigForm = useForm<StorageProviderConfigData>({
    defaultValues: initialValue?.providerConfig,
  });
  const [createStorageDef, { loading: creating }] = useMutation(
    createStorageDefMutation
  );
  const [updateStorageDef, { loading: updating }] = useMutation(
    updateStorageDefMutation
  );
  const isSubmitting = creating || updating;
  const onSubmit = React.useCallback(() => {
    if (isSubmitting) {
      return;
    }
    const commonData = commonFieldsForm.getValues();
    const providerConfigData = providerConfigForm.getValues();
    let task: Promise<{ id?: string }>;
    if (id) {
      task = updateStorageDef({
        variables: {
          input: {
            identifier: commonData.identifier,
            priority: commonData.priority,
            isEnabled: commonData.isEnabled,
            configJSON: JSON.stringify(providerConfigData),
          },
        },
      }).then((res) => ({ id: res.data?.updateStorageDefinition?.id }));
    } else {
      task = createStorageDef({
        variables: {
          input: {
            identifier: commonData.identifier,
            storageType:
              commonData.storageType == "S3"
                ? StorageTypeEnum.S3
                : StorageTypeEnum.Other,
            configJSON: JSON.stringify(providerConfigData),
            isEnabled: true,
            priority: commonData.priority,
          },
        },
      }).then((res) => ({ id: res.data?.createStorageDefinition?.id }));
    }
    task
      .then(({ id }) => {
        afterSave?.(id);
      })
      .catch(() => {
        afterSave?.();
      });
  }, [
    isSubmitting,
    commonFieldsForm,
    providerConfigForm,
    id,
    updateStorageDef,
    createStorageDef,
    afterSave,
  ]);
  return (
    <div>
      <form>
        <div className="flex flex-col gap-4 mt-6">
          <SelectWithLabel
            containerClassName="flex flex-col gap-1 max-w-full"
            label="Storage Type"
            value={storageTypeValue}
            {...commonFieldsForm.register("storageType", {
              value: storageTypeValue,
            })}
            disabled={!!id}
          >
            <option value="S3">S3</option>
          </SelectWithLabel>
          <InputWithLabel
            containerClassName="flex flex-col gap-1 max-w-full"
            label="Identifier"
            {...commonFieldsForm.register("identifier", { required: true })}
            disabled={!!id}
          />
          <InputWithLabel
            containerClassName="flex flex-col gap-1 max-w-full"
            label="Enabled"
            type="checkbox"
            {...commonFieldsForm.register("isEnabled")}
          />
          <InputWithLabel
            containerClassName="flex flex-col gap-1 max-w-full"
            label="Priority"
            {...commonFieldsForm.register("priority", { valueAsNumber: true })}
            type="number"
          />
          {storageTypeValue === "S3" ? (
            <S3ProviderConfigForm form={providerConfigForm} />
          ) : null}
        </div>
      </form>
      <Button
        onClick={onSubmit}
        className="mt-6 w-full"
        disabled={isSubmitting}
      >
        Save
      </Button>
    </div>
  );
}
