import React, { useState } from "react";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { ListUsersQuery } from "~src/__generated__/graphql";
import classNames from "classnames";
import {
  SECONDARY_TEXT_COLOR_DIM,
  SECOND_LAYER,
  LINK_COLOR,
} from "~src/ui/classNames";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { HiChevronRight, HiChevronDown } from "react-icons/hi";

type User = ListUsersQuery["viewer"]["paginatedAllUsers"]["nodes"][number];
type OrganizationUser = User["organizationUsers"][number];

const columnHelper = createColumnHelper<User | OrganizationUser>();

export function UsersTable({ data }: { data: User[] }) {
  const { t } = useTranslation();
  const [expandedRows, setExpandedRows] = useState<Record<string, boolean>>({});

  const toggleRow = (userId: string) => {
    setExpandedRows((prev) => ({ ...prev, [userId]: !prev[userId] }));
  };

  const columns = React.useMemo(
    () => [
      columnHelper.display({
        id: "expand",
        header: "",
        cell: ({ row }) =>
          "organizationUsers" in row.original ? (
            <button
              onClick={() => toggleRow(row.original.id)}
              className="flex items-center justify-center text-gray-700 dark:text-gray-300"
            >
              {expandedRows[row.original.id] ? (
                <HiChevronDown className="text-xl" />
              ) : (
                <HiChevronRight className="text-xl" />
              )}
            </button>
          ) : null,
      }),
      columnHelper.accessor("id", {
        header: "ID",
        cell: (info) => info.getValue(),
      }),
      columnHelper.accessor("email", {
        header: t("usersTable.email"),
        cell: (info) =>
          "organizationUsers" in info.row.original ? (
            info.getValue()
          ) : (
            <span className="pl-8">{info.getValue() as string}</span>
          ),
      }),
      columnHelper.display({
        id: "actions",
        header: "",
        cell: ({ row }) =>
          "organizationUsers" in row.original ? null : (
            <Link
              to={`${row.original.id}/images`}
              className={classNames(LINK_COLOR, "ml-2")}
            >
              {t("usersTable.viewImagesButton")}
            </Link>
          ),
      }),
    ],
    [t, expandedRows],
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
            <React.Fragment key={row.id}>
              <tr
                className={classNames(
                  SECOND_LAYER,
                  "border-t border-neutral-300 dark:border-neutral-700",
                )}
              >
                {row.getVisibleCells().map((cell) => (
                  <td key={cell.id} className={classNames("px-6 py-4")}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
              {expandedRows[row.original.id] &&
                "organizationUsers" in row.original &&
                row.original.organizationUsers.map(
                  (orgUser: OrganizationUser) => (
                    <tr
                      key={orgUser.id}
                      className="bg-neutral-200 dark:bg-neutral-800 border-t border-neutral-300 dark:border-neutral-700"
                    >
                      <td className="px-6 py-4" />
                      <td className="px-6 py-4 text-gray-600 dark:text-gray-400">
                        {orgUser.organization.name}
                      </td>
                      <td className="px-6 py-4" />
                      <td className="px-6 py-4">
                        <Link
                          to={`${orgUser.id}/images`}
                          className={classNames(LINK_COLOR, "ml-2")}
                        >
                          {t("usersTable.viewImagesButton")}
                        </Link>
                      </td>
                    </tr>
                  ),
                )}
            </React.Fragment>
          ))}
        </tbody>
      </table>
      {data.length === 0 && (
        <div className="text-center py-4 italic text-gray-500">
          {t("usersTable.noUsers")}
        </div>
      )}
    </div>
  );
}
