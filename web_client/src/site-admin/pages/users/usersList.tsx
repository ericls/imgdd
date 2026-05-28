import React, { useState } from "react";
import { gql } from "~src/__generated__";
import { HEADING_2 } from "~src/ui/classNames";
import { BlockLoader } from "~src/ui/loader";
import { useTranslation } from "react-i18next";
import { UsersTable } from "./usersTable";
import { Input } from "~src/ui/input";
import { useDebounce } from "~src/lib/hooks";
import { useQuery } from "@apollo/client/react";
import { useSearchParams } from "react-router";
import classNames from "classnames";
import { Select } from "~src/ui/select";

const listUsersQuery = gql(`
  query ListUsers($limit: Int, $offset: Int, $search: String) {
    viewer {
      id
      paginatedAllUsers(limit: $limit, offset: $offset, search: $search) {
        nodes {
          id
          email
          organizationUsers {
            id
            organization {
              id
              name
            }
          }
        }
        pageInfo {
          totalCount
          hasNextPage
          hasPreviousPage
        }
      }
    }
  }
`);

export function UsersList() {
  const { t } = useTranslation();
  const [searchTerm, setSearchTerm] = useState("");
  const [searchParams, setSearchParams] = useSearchParams();
  const currentPage = parseInt(searchParams.get("page") ?? "0", 10);
  const pageSize = parseInt(searchParams.get("pageSize") ?? "10", 10);
  const debouncedSearchTerm = useDebounce(searchTerm, 500);
  const [prevDebouncedSearchTerm, setPrevDebouncedSearchTerm] =
    useState(debouncedSearchTerm);
  if (prevDebouncedSearchTerm !== debouncedSearchTerm) {
    setPrevDebouncedSearchTerm(debouncedSearchTerm);
    setSearchParams(
      (p) => {
        p.set("page", "0");
        return p;
      },
      { replace: true },
    );
  }

  const { data: usersData, loading } = useQuery(listUsersQuery, {
    fetchPolicy: "cache-and-network",
    variables: {
      limit: pageSize,
      offset: currentPage * pageSize,
      search: debouncedSearchTerm || undefined,
    },
  });

  const users = usersData?.viewer.paginatedAllUsers.nodes || [];
  const pageInfo = usersData?.viewer.paginatedAllUsers.pageInfo;
  const totalPages = pageInfo ? Math.ceil(pageInfo.totalCount / pageSize) : 0;

  return (
    <div className="m-auto max-w-5xl">
      <h1 className={classNames(HEADING_2, "font-poppins")}>
        {t("usersList.title", "Users")}
      </h1>

      <div className="my-4 flex items-center gap-4 mt-2">
        <Input
          type="text"
          placeholder={t("usersList.searchPlaceholder", "Search by email...")}
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="w-full max-w-md"
        />
        <div className="flex items-center gap-2 shrink-0">
          <label className="text-sm text-gray-600 dark:text-gray-400 whitespace-nowrap">
            {t("usersList.pageSize", "Per page")}
          </label>
          <Select
            value={pageSize}
            onChange={(e) =>
              setSearchParams((p) => {
                p.set("pageSize", e.target.value);
                p.set("page", "0");
                return p;
              })
            }
          >
            {[10, 25, 50, 100, 200].map((n) => (
              <option key={n} value={n}>
                {n}
              </option>
            ))}
          </Select>
        </div>
      </div>

      {!usersData && loading && <BlockLoader />}

      <div className="mt-4">
        <UsersTable data={users} />
      </div>

      {/* Pagination controls */}
      {pageInfo && totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <div className="text-sm text-gray-700 dark:text-gray-300">
            {t(
              "usersList.showingResults",
              "Showing {{from}} to {{to}} of {{total}} results",
              {
                from: currentPage * pageSize + 1,
                to: Math.min((currentPage + 1) * pageSize, pageInfo.totalCount),
                total: pageInfo.totalCount,
              },
            )}
          </div>
          <div className="flex gap-2">
            <button
              onClick={() =>
                setSearchParams((p) => {
                  p.set("page", String(currentPage - 1));
                  return p;
                })
              }
              disabled={!pageInfo.hasPreviousPage}
              className="px-3 py-1 text-sm rounded-sm border border-gray-300 dark:border-gray-600 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              {t("usersList.previous", "Previous")}
            </button>
            <span className="px-3 py-1 text-sm">
              {t("usersList.pageInfo", "Page {{current}} of {{total}}", {
                current: currentPage + 1,
                total: totalPages,
              })}
            </span>
            <button
              onClick={() =>
                setSearchParams((p) => {
                  p.set("page", String(currentPage + 1));
                  return p;
                })
              }
              disabled={!pageInfo.hasNextPage}
              className="px-3 py-1 text-sm rounded-sm border border-gray-300 dark:border-gray-600 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              {t("usersList.next", "Next")}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
