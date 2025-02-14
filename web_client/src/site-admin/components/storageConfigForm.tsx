import { useMutation } from "@apollo/client";
import React from "react";
import { useForm } from "react-hook-form";
import { InputWithLabel } from "~src/ui/input";
import { SelectWithLabel } from "~src/ui/select";
import { createStorageDefMutation, updateStorageDefMutation } from "../types";
import { StorageTypeEnum } from "~src/__generated__/graphql";
import { Button } from "~src/ui/button";
import { useTranslation } from "react-i18next";

type S3StorageConfigData = {
  bucket: string;
  endpoint: string;
  access: string;
  secret: string;
};

const EmptyS3Config: S3StorageConfigData = {
  bucket: "",
  endpoint: "",
  access: "",
  secret: "",
};

type FSSotrageConfigData = {
  mediaRoot: string;
};

const EmptyFSConfig: FSSotrageConfigData = {
  mediaRoot: "",
};

type StorageProviderConfigData =
  | S3StorageConfigData
  | FSSotrageConfigData
  | { __other: string };

type StorageConfigData = {
  storageType: "S3" | "FS" | "__other";
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
  const { t } = useTranslation();
  return (
    <>
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label={t("storageConfigForm.bucket")}
        {...form.register("bucket", { required: true })}
      />
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label={t("storageConfigForm.endpoint")}
        {...form.register("endpoint", { required: true })}
      />
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label={t("storageConfigForm.access")}
        {...form.register("access", { required: true })}
      />
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label={t("storageConfigForm.secret")}
        type="password"
        {...form.register("secret", {
          required: true,
        })}
      />
    </>
  );
}

function FSSotrageConfigForm({
  form,
}: {
  form: ReturnType<typeof useForm<StorageProviderConfigData>>;
}) {
  const { t } = useTranslation();
  return (
    <InputWithLabel
      containerClassName="flex flex-col gap-1 max-w-full"
      label={t("storageConfigForm.mediaRootPath")}
      {...form.register("mediaRoot", { required: true })}
    />
  );
}

export function StorageConfigForm({
  initialValue,
  id,
  afterSave,
}: StorageConfigFormProps) {
  const { t } = useTranslation();
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
    createStorageDefMutation,
  );
  const [updateStorageDef, { loading: updating }] = useMutation(
    updateStorageDefMutation,
  );
  const isSubmitting = creating || updating;
  const onSubmit = React.useCallback(() => {
    if (isSubmitting) {
      return;
    }
    const commonData = commonFieldsForm.getValues();
    const providerConfigData = providerConfigForm.getValues();
    let task: Promise<{ id?: string }>;
    const providerConfigDataMask: StorageProviderConfigData = (() => {
      if (storageTypeValue === "S3") {
        return EmptyS3Config;
      } else if (storageTypeValue === "FS") {
        return EmptyFSConfig;
      } else {
        return { __other: "" };
      }
    })();
    const maskedProviderConfigData = Object.keys(providerConfigDataMask).reduce(
      (acc, key) => {
        const k = key as keyof StorageProviderConfigData;
        acc[k] = providerConfigData[k];
        return acc;
      },
      {} as StorageProviderConfigData,
    );
    const configJSON = JSON.stringify(maskedProviderConfigData);
    if (id) {
      task = updateStorageDef({
        variables: {
          input: {
            identifier: commonData.identifier,
            priority: commonData.priority,
            isEnabled: commonData.isEnabled,
            configJSON,
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
                : commonData.storageType == "FS"
                  ? StorageTypeEnum.Fs
                  : StorageTypeEnum.Other,
            configJSON,
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
    storageTypeValue,
    updateStorageDef,
    createStorageDef,
    afterSave,
  ]);

  const providerConfigFields = React.useMemo(() => {
    if (storageTypeValue === "S3") {
      return <S3ProviderConfigForm form={providerConfigForm} />;
    }
    if (storageTypeValue === "FS") {
      return <FSSotrageConfigForm form={providerConfigForm} />;
    }
    return null;
  }, [storageTypeValue, providerConfigForm]);
  return (
    <div>
      <form>
        <div className="flex flex-col gap-4 mt-6">
          <SelectWithLabel
            containerClassName="flex flex-col gap-1 max-w-full"
            label={t("storageConfigForm.storageType")}
            value={storageTypeValue}
            {...commonFieldsForm.register("storageType", {
              value: storageTypeValue,
            })}
            disabled={!!id}
          >
            <option value="S3">S3</option>
            <option value="FS">{t("storageTypeNameTitle.fs")}</option>
          </SelectWithLabel>
          <InputWithLabel
            containerClassName="flex flex-col gap-1 max-w-full"
            label={t("storageConfigForm.identifier")}
            {...commonFieldsForm.register("identifier", { required: true })}
            disabled={!!id}
          />
          <InputWithLabel
            containerClassName="flex flex-col gap-1 max-w-full"
            label={t("storageConfigForm.enabled")}
            type="checkbox"
            {...commonFieldsForm.register("isEnabled")}
          />
          <InputWithLabel
            containerClassName="flex flex-col gap-1 max-w-full"
            label={t("storageConfigForm.priority")}
            {...commonFieldsForm.register("priority", { valueAsNumber: true })}
            type="number"
          />
          {providerConfigFields}
        </div>
      </form>
      <Button
        onClick={onSubmit}
        className="mt-6 w-full"
        disabled={isSubmitting}
      >
        {t("common.buttonLabel.save")}
      </Button>
    </div>
  );
}
