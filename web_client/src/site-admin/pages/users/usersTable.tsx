import React from "react";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { ListUsersQuery } from "~src/__generated__/graphql";
import classNames from "classnames";
import { SECONDARY_TEXT_COLOR_DIM, SECOND_LAYER } from "~src/ui/classNames";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";

type User = ListUsersQuery["viewer"]["paginatedAllUsers"]["nodes"][number];

const columnHelper = createColumnHelper<User>();

export function UsersTable({ data }: { data: User[] }) {
  const { t } = useTranslation();
  const columns = React.useMemo(
    () => [
      columnHelper.accessor("id", {
        header: "ID",
        cell: (info) => info.getValue(),
      }),
      columnHelper.accessor("email", {
        header: t("usersTable.email", "Email"),
        cell: (info) => info.getValue(),
      }),
      columnHelper.display({
        id: "viewImages",
        header: t("usersTable.viewImages", "View Images"),
        cell: ({ row }) => (
          <Link
            to={`${row.original.id}/images`}
            className="text-blue-500 hover:underline"
          >
            {t("usersTable.viewImagesButton", "View Images")}
          </Link>
        ),
      }),
    ],
    [t],
  );
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });
  return (
    <div className="relative overflow-x-auto">
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
      {data.length === 0 && (
        <div className="text-center py-4 italic text-gray-500">
          {t("usersTable.noUsers", "No users found")}
        </div>
      )}
    </div>
  );
}
