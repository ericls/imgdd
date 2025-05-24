import React, { useState } from "react";
import { gql } from "~src/__generated__";
import { HEADING_2 } from "~src/ui/classNames";
import { BlockLoader } from "~src/ui/loader";
import { useTranslation } from "react-i18next";
import { UsersTable } from "./usersTable";
import { Input } from "~src/ui/input";
import { useDebounce } from "~src/lib/hooks";
import { useQuery } from "@apollo/client";
import classNames from "classnames";

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
  const [currentPage, setCurrentPage] = useState(0);
  const pageSize = 10;
  const debouncedSearchTerm = useDebounce(searchTerm, 500);

  const { data: usersData, loading } = useQuery(listUsersQuery, {
    fetchPolicy: "cache-and-network",
    variables: {
      limit: pageSize,
      offset: currentPage * pageSize,
      search: debouncedSearchTerm || undefined,
    },
  });

  // Reset to first page when search term changes
  React.useEffect(() => {
    setCurrentPage(0);
  }, [debouncedSearchTerm]);

  const users = usersData?.viewer.paginatedAllUsers.nodes || [];
  const pageInfo = usersData?.viewer.paginatedAllUsers.pageInfo;
  const totalPages = pageInfo ? Math.ceil(pageInfo.totalCount / pageSize) : 0;

  return (
    <div className="m-auto max-w-5xl">
      <h1 className={classNames(HEADING_2, "font-poppins")}>
        {t("usersList.title", "Users")}
      </h1>

      <div className="my-4 w-full max-w-md mt-2">
        <Input
          type="text"
          placeholder={t("usersList.searchPlaceholder", "Search by email...")}
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
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
              onClick={() => setCurrentPage(currentPage - 1)}
              disabled={!pageInfo.hasPreviousPage}
              className="px-3 py-1 text-sm rounded border border-gray-300 dark:border-gray-600 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-700"
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
              onClick={() => setCurrentPage(currentPage + 1)}
              disabled={!pageInfo.hasNextPage}
              className="px-3 py-1 text-sm rounded border border-gray-300 dark:border-gray-600 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              {t("usersList.next", "Next")}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
