export type ImageItem = {
  id: string;
  url: string;
  name: string;
  nominalWidth: number;
  nominalHeight: number;
  nominalByteSize: number;
  createdAt: string;
};

export type ImageItemRenderer = (image: ImageItem) => React.ReactNode;
