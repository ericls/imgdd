import React from "react";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { HiOutlinePencilSquare as EditIcon } from "react-icons/hi2";
import { ListStorageDefQuery } from "~src/__generated__/graphql";
import classNames from "classnames";
import {
  PRIMARY_TEXT_COLOR,
  SECONDARY_TEXT_COLOR_DIM,
  SECOND_LAYER,
} from "~src/ui/classNames";

import { Button } from "~src/ui/button";
import { useTranslation } from "react-i18next";
import { StorageDefTableConnectivityCell } from "./connectivityCell";

type StorageDef = ListStorageDefQuery["viewer"]["storageDefinitions"][number];

const columnHelper = createColumnHelper<StorageDef>();

export function DumbStorageDefTable({
  data,
  onAddNew,
  onEdit,
}: {
  data: StorageDef[];
  onAddNew?: () => void;
  onEdit?: (id: string) => void;
}) {
  const { t } = useTranslation();
  const columns = React.useMemo(
    () => [
      columnHelper.accessor("identifier", {
        header: t("storageDefTable.identifier"),
        cell: (info) => info.getValue(),
      }),
      columnHelper.accessor((row) => row.config.__typename, {
        id: "storageType",
        cell: (info) => {
          const value = info.getValue();
          const storageTypeStr = (() => {
            if (value === "S3StorageConfig") {
              return "S3";
            } else {
              return "Unknown";
            }
          })();
          return <span>{storageTypeStr}</span>;
        },
        header: () => <span>{t("storageDefTable.storageType")}</span>,
      }),
      columnHelper.accessor("isEnabled", {
        header: t("storageDefTable.enabled"),
        cell: (info) => {
          const value = info.getValue();
          return <span>{value ? "Yes" : "No"}</span>;
        },
      }),
      columnHelper.accessor("priority", {
        header: t("storageDefTable.priority"),
      }),
      columnHelper.accessor("id", {
        header: t("storageDefTable.connectivity"),
        cell: (info) => {
          return <StorageDefTableConnectivityCell id={info.getValue()} />;
        },
      }),
      columnHelper.accessor(
        (row) => {
          return row;
        },
        {
          header: t("storageDefTable.actions"),
          id: "actions",
          cell: (row) => {
            // Edit button
            return (
              <div className="flex items-center gap-4">
                {onEdit ? (
                  <Button
                    variant="transparent"
                    noPadding
                    className={classNames(
                      "py-2 flex items-center gap-1",
                      PRIMARY_TEXT_COLOR,
                    )}
                    onClick={() => {
                      onEdit(row.getValue().id);
                    }}
                  >
                    <EditIcon />
                    {t("common.buttonLabel.edit")}
                  </Button>
                ) : null}
              </div>
            );
          },
        },
      ),
    ],
    [onEdit, t],
  );
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });
  return (
    <div className="relative overflow-x-auto">
      <div className="flex justify-end items-center mb-2 p-1">
        {onAddNew ? (
          <Button variant="secondary" onClick={onAddNew} className={"text-sm"}>
            {t("storageDefTable.addNew")}
          </Button>
        ) : null}
      </div>
      <table className="w-full text-sm text-left rounded-md overflow-hidden">
        <thead
          className={classNames(
            "bg-neutral-50 dark:bg-neutral-700",
            SECONDARY_TEXT_COLOR_DIM,
            "uppercase",
          )}
        >
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id} className="px-6 py-3">
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row) => (
            <tr key={row.id} className={classNames(SECOND_LAYER)}>
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id} className={classNames("px-6 py-4")}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
