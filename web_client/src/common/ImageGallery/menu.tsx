import React from "react";
import { RenderingImageItem } from "./types";
import { MenuItem, MenuSections } from "~src/ui/menu";
import { useDeleteImage } from "./data";
import { copyText } from "~src/lib/copyText";
import { toast } from "react-toastify";
import { absoluteURL } from "~src/lib/url";

enum ImageMenuItemName {
  DOWNLOAD = "download",
  COPY_URL = "copy-url",
  DELETE = "delete",
}

export type ImageItemMenuConfig = {
  sections: { id: string; names: ImageMenuItemName[] }[];
};

export const DEFAULT_MENU_CONFIG: ImageItemMenuConfig = {
  sections: [
    {
      id: "actions",
      names: [ImageMenuItemName.DOWNLOAD, ImageMenuItemName.COPY_URL],
    },
    {
      id: "delete",
      names: [ImageMenuItemName.DELETE],
    },
  ],
};

type MenuItemGetterProps = {
  image: RenderingImageItem;
  onDelete?: () => PromiseLike<unknown>;
};

type MenuItemGetter = (props: MenuItemGetterProps) => MenuItem;

const MENU_ITEM_GETTERS: Record<ImageMenuItemName, MenuItemGetter> = {
  [ImageMenuItemName.DOWNLOAD]: getDownloadMenuItem,
  [ImageMenuItemName.COPY_URL]: getCopyURLMenuItem,
  [ImageMenuItemName.DELETE]: getDeleteMenuItem,
};

function getMenuItemByName(
  name: ImageMenuItemName,
  props: MenuItemGetterProps
): MenuItem {
  return MENU_ITEM_GETTERS[name](props);
}

function getDownloadMenuItem({
  image: { url },
}: MenuItemGetterProps): MenuItem {
  return {
    id: ImageMenuItemName.DOWNLOAD,
    children: "Download",
    action: () => {
      window.open(absoluteURL(url), "_blank");
    },
  };
}

function getCopyURLMenuItem({ image: { url } }: MenuItemGetterProps): MenuItem {
  return {
    id: ImageMenuItemName.COPY_URL,
    children: "Copy URL",
    action: () => {
      copyText(absoluteURL(url), () => toast("Copied"));
    },
  };
}

function getDeleteMenuItem({ onDelete }: MenuItemGetterProps): MenuItem {
  return {
    id: ImageMenuItemName.DELETE,
    children: "Delete",
    action: () => {
      onDelete?.().then(() => {
        toast("Deleted");
      });
    },
  };
}
export function useImageItemMenu(
  image: RenderingImageItem,
  config?: ImageItemMenuConfig
): MenuSections | null {
  const { execute: executeDelete } = useDeleteImage(image.id);
  const menuSections = React.useMemo(() => {
    if (!config) return null;
    return config.sections.map((section) => {
      const items = section.names.map((name) =>
        getMenuItemByName(name, {
          image,
          onDelete: executeDelete,
        })
      );
      return { id: section.id, items };
    });
  }, [image, config, executeDelete]);
  if (!menuSections) return null;
  return {
    children: menuSections,
  };
}
