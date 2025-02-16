import {
  StorageTypeEnum,
  type GetStorageDefQuery,
} from "~src/__generated__/graphql";
import { type StorageType } from "~src/site-admin/types";

type ProviderConfig = NonNullable<
  GetStorageDefQuery["viewer"]["getStorageDefinition"]
>["config"];

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

type WebDAVStorageConfigData = {
  url: string;
  username: string;
  password: string;
};
const EmptyWebDAVConfig: WebDAVStorageConfigData = {
  url: "",
  username: "",
  password: "",
};

export type StorageProviderConfigData =
  | S3StorageConfigData
  | FSSotrageConfigData
  | WebDAVStorageConfigData
  | { __other: string };

const EmptyConfigs = {
  S3: EmptyS3Config,
  FS: EmptyFSConfig,
  WebDAV: EmptyWebDAVConfig,
  __other: { __other: "" },
} as const;
export const StorageProviders: {
  [key in StorageType]: {
    emptyConfig: StorageProviderConfigData;
    enum: StorageTypeEnum;
    mask: (
      data: StorageProviderConfigData | { __typename?: string },
    ) => StorageProviderConfigData;
  };
} = {
  S3: {
    emptyConfig: EmptyConfigs.S3,
    enum: StorageTypeEnum.S3,
    mask: maskStorageProviderConfigData("S3"),
  },
  FS: {
    emptyConfig: EmptyConfigs.FS,
    enum: StorageTypeEnum.Fs,
    mask: maskStorageProviderConfigData("FS"),
  },
  WebDAV: {
    emptyConfig: EmptyConfigs.WebDAV,
    enum: StorageTypeEnum.WebDav,
    mask: maskStorageProviderConfigData("WebDAV"),
  },
  __other: {
    emptyConfig: { __other: "" },
    enum: StorageTypeEnum.Other,
    mask: maskStorageProviderConfigData("__other"),
  },
} as const;

function maskStorageProviderConfigData(
  storageType: StorageType,
): (
  data: StorageProviderConfigData | { __typename?: string },
) => StorageProviderConfigData {
  const mask = EmptyConfigs[storageType];
  return (data) => {
    if (Object.keys(data).length === 0) {
      return mask;
    }
    return Object.keys(mask).reduce((acc, key) => {
      const k = key as keyof typeof mask;
      return { ...acc, [k]: data[k] };
    }, {} as StorageProviderConfigData);
  };
}

export function getStorageTypeFromConfig(config: ProviderConfig): StorageType {
  if (config.__typename === "S3StorageConfig") {
    return "S3";
  }
  if (config.__typename === "FSStorageConfig") {
    return "FS";
  }
  if (config.__typename === "WebDAVStorageConfig") {
    return "WebDAV";
  }
  return "__other";
}
