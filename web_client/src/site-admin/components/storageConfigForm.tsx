import { useMutation } from "@apollo/client";
import React from "react";
import { useForm } from "react-hook-form";
import { InputWithLabel } from "~src/ui/input";
import { SelectWithLabel } from "~src/ui/select";
import {
  createStorageDefMutation,
  StorageType,
  updateStorageDefMutation,
} from "../types";
import { Button } from "~src/ui/button";
import { useTranslation } from "react-i18next";
import {
  StorageProviderConfigData,
  StorageProviders,
} from "~src/site-admin/storageProviderDefs";

type StorageConfigData = {
  storageType: StorageType;
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

function WebDavProviderConfigForm({
  form,
}: {
  form: ReturnType<typeof useForm<StorageProviderConfigData>>;
}) {
  const { t } = useTranslation();
  return (
    <>
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label={t("storageConfigForm.url")}
        {...form.register("url", { required: true })}
      />
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label={t("storageConfigForm.username")}
        {...form.register("username", { required: false })}
      />
      <InputWithLabel
        containerClassName="flex flex-col gap-1 max-w-full"
        label={t("storageConfigForm.password")}
        type="password"
        {...form.register("password", { required: false })}
      />
    </>
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
      priority: initialValue?.priority ?? 1,
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
    const configJSON = JSON.stringify(
      StorageProviders[storageTypeValue].mask(providerConfigData),
    );
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
            storageType: StorageProviders[storageTypeValue].enum,
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
    if (storageTypeValue === "WebDAV") {
      return <WebDavProviderConfigForm form={providerConfigForm} />;
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
            <option value="WebDAV">{t("storageTypeNameTitle.WebDAV")}</option>
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
