import React from "react";
import { RenderingImageItem } from "./types";
import { MenuItem, MenuSections } from "~src/ui/menu";
import { useDeleteImage } from "./data";
import { copyText } from "~src/lib/copyText";
import { toast } from "react-toastify";
import { absoluteURL } from "~src/lib/url";
import { prompt } from "~src/ui/prompt";
import i18n from "~src/i18n";
import { Trans } from "~node_modules/react-i18next";

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
  props: MenuItemGetterProps,
): MenuItem {
  return MENU_ITEM_GETTERS[name](props);
}

function getDownloadMenuItem({
  image: { url },
}: MenuItemGetterProps): MenuItem {
  return {
    id: ImageMenuItemName.DOWNLOAD,
    children: i18n.t("imageItem.download"),
    action: () => {
      window.open(absoluteURL(url), "_blank");
    },
  };
}

function getCopyURLMenuItem({ image: { url } }: MenuItemGetterProps): MenuItem {
  return {
    id: ImageMenuItemName.COPY_URL,
    children: i18n.t("imageItem.copyURL"),
    action: () => {
      copyText(absoluteURL(url), () => toast(i18n.t("common.toast.copied")));
    },
  };
}

function getDeleteMenuItem({
  onDelete,
  image: { name },
}: MenuItemGetterProps): MenuItem {
  return {
    id: ImageMenuItemName.DELETE,
    children: i18n.t("common.buttonLabel.delete"),
    action: () => {
      prompt({
        content: (
          <Trans i18nKey={"imageItem.confirmDelete"} values={{ name }}></Trans>
        ),
        title: "Delete",
        yesDestructive: true,
        showCancel: true,
      }).then((confirmed) => {
        if (!confirmed) return;
        onDelete?.().then(() => {
          toast(i18n.t("common.toast.deleted"));
        });
      });
    },
  };
}
export function useImageItemMenu(
  image: RenderingImageItem,
  config?: ImageItemMenuConfig,
): MenuSections | null {
  const { execute: executeDelete } = useDeleteImage(image.id);
  const menuSections = React.useMemo(() => {
    if (!config) return null;
    return config.sections.map((section) => {
      const items = section.names.map((name) =>
        getMenuItemByName(name, {
          image,
          onDelete: executeDelete,
        }),
      );
      return { id: section.id, items };
    });
  }, [image, config, executeDelete]);
  if (!menuSections) return null;
  return {
    children: menuSections,
  };
}
