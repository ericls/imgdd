import { ImageItemMenuConfig } from "./menu";

export type RenderingImageItem = {
  id: string;
  url: string;
  name: string;
  nominalWidth: number;
  nominalHeight: number;
  nominalByteSize: number;
  createdAt: string;
  parent?: {
    id: string;
    name: string;
  } | null;
  createdBy?: {
    id: string;
    user: {
      id: string;
      avatarUrl: string;
    };
  } | null;
} & {
  menuConfig?: ImageItemMenuConfig;
};

export type ImageItemRenderer = (image: RenderingImageItem) => React.ReactNode;
