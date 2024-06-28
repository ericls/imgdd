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
  const columns = React.useMemo(
    () => [
      columnHelper.accessor("identifier", {
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
        header: () => <span>Storage Type</span>,
      }),
      columnHelper.accessor("isEnabled", {
        header: "Enabled",
      }),
      columnHelper.accessor("priority", {
        header: "Priority",
      }),
      columnHelper.accessor(
        (row) => {
          return row;
        },
        {
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
                      PRIMARY_TEXT_COLOR
                    )}
                    onClick={() => {
                      onEdit(row.getValue().id);
                    }}
                  >
                    <EditIcon />
                    Edit
                  </Button>
                ) : null}
              </div>
            );
          },
        }
      ),
    ],
    [onEdit]
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
            Add New
          </Button>
        ) : null}
      </div>
      <table className="w-full text-sm text-left">
        <thead
          className={classNames(
            "bg-neutral-50 dark:bg-neutral-700",
            SECONDARY_TEXT_COLOR_DIM,
            "uppercase"
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
                        header.getContext()
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
