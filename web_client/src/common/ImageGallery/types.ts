import { ImageItemMenuConfig } from "./menu";

export type RenderingImageItem = {
  id: string;
  url: string;
  name: string;
  nominalWidth: number;
  nominalHeight: number;
  nominalByteSize: number;
  createdAt: string;
} & {
  menuConfig?: ImageItemMenuConfig;
};

export type ImageItemRenderer = (image: RenderingImageItem) => React.ReactNode;
