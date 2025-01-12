export type RenderingImageItem = {
  id: string;
  url: string;
  name: string;
  nominalWidth: number;
  nominalHeight: number;
  nominalByteSize: number;
  createdAt: string;
};

export type ImageItemRenderer = (image: RenderingImageItem) => React.ReactNode;
