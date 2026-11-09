import React from "react";
import classNames from "classnames";
import dayjs from "dayjs";
import { Link } from "react-router";
import { useMutation, useQuery } from "@apollo/client/react";
import { useTranslation } from "react-i18next";
import { toast } from "react-toastify";
import { HiOutlineClipboardDocumentList } from "react-icons/hi2";
import { gql } from "~src/__generated__";
import { copyText } from "~src/lib/copyText";
import { routes } from "~src/routes";
import { Button } from "~src/ui/button";
import {
  HEADING_2,
  LINK_COLOR,
  SECOND_LAYER,
  SECONDARY_TEXT_COLOR_DIM,
  TEXT_COLOR,
} from "~src/ui/classNames";
import { Input } from "~src/ui/input";
import { prompt } from "~src/ui/prompt";

const API_KEYS_QUERY = gql(`
  query APIKeys {
    viewer {
      id
      apiKeys {
        id
        name
        keyPrefix
        createdAt
        lastUsedAt
      }
    }
  }
`);

const CREATE_API_KEY_MUTATION = gql(`
  mutation CreateAPIKey($input: CreateAPIKeyInput!) {
    createAPIKey(input: $input) {
      apiKey {
        id
      }
      secret
    }
  }
`);

const REVOKE_API_KEY_MUTATION = gql(`
  mutation RevokeAPIKey($input: RevokeAPIKeyInput!) {
    revokeAPIKey(input: $input) {
      id
    }
  }
`);

function formatTime(value: string | null | undefined) {
  if (!value) return null;
  return dayjs(value).format("YYYY-MM-DD HH:mm");
}

function NewSecret({ secret, onDone }: { secret: string; onDone: () => void }) {
  const { t } = useTranslation();
  const onCopy = React.useCallback(() => {
    copyText(secret, () => toast(t("common.toast.copied")));
  }, [secret, t]);
  return (
    <div
      className={classNames(
        "rounded-md border border-green-600/50 p-4 my-4",
        SECOND_LAYER,
        TEXT_COLOR,
      )}
    >
      <p className="font-medium mb-2">{t("profile.apiKeys.newKeyTitle")}</p>
      <p className={classNames("text-sm mb-3", SECONDARY_TEXT_COLOR_DIM)}>
        {t("profile.apiKeys.newKeyWarning")}
      </p>
      <div className="flex gap-2 items-center">
        <code className="grow break-all font-mono text-sm bg-neutral-100 dark:bg-neutral-900 rounded-sm p-2 select-all">
          {secret}
        </code>
        <Button
          variant="secondary"
          onClick={onCopy}
          className="flex items-center gap-1 shrink-0"
        >
          <HiOutlineClipboardDocumentList />
          {t("profile.apiKeys.copy")}
        </Button>
      </div>
      <div className="mt-3 flex justify-end">
        <Button variant="transparent" onClick={onDone}>
          {t("profile.apiKeys.done")}
        </Button>
      </div>
    </div>
  );
}

export function APIKeys() {
  const { t } = useTranslation();
  const { data, loading, refetch } = useQuery(API_KEYS_QUERY, {
    fetchPolicy: "cache-and-network",
  });
  const [createAPIKey, { loading: creating }] = useMutation(
    CREATE_API_KEY_MUTATION,
  );
  const [revokeAPIKey] = useMutation(REVOKE_API_KEY_MUTATION);
  const [name, setName] = React.useState("");
  const [newSecret, setNewSecret] = React.useState<string | null>(null);

  const onCreate = React.useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      const trimmed = name.trim();
      if (!trimmed) return;
      try {
        const result = await createAPIKey({
          variables: { input: { name: trimmed } },
        });
        const secret = result.data?.createAPIKey.secret;
        if (secret) {
          setNewSecret(secret);
          setName("");
        }
        await refetch();
      } catch (err) {
        toast.error((err as Error).message);
      }
    },
    [createAPIKey, name, refetch],
  );

  const onRevoke = React.useCallback(
    async (id: string, keyName: string) => {
      const confirmed = await prompt({
        title: t("profile.apiKeys.revokeConfirmTitle"),
        content: (
          <p>{t("profile.apiKeys.revokeConfirmBody", { name: keyName })}</p>
        ),
        yesText: t("profile.apiKeys.revoke"),
        yesDestructive: true,
      });
      if (!confirmed) return;
      try {
        await revokeAPIKey({ variables: { input: { id } } });
        toast(t("profile.apiKeys.revoked"));
        await refetch();
      } catch (err) {
        toast.error((err as Error).message);
      }
    },
    [refetch, revokeAPIKey, t],
  );

  const apiKeys = data?.viewer.apiKeys ?? [];

  return (
    <div className="m-auto max-w-3xl mx-8">
      <h1 className={classNames(HEADING_2, "font-poppins")}>
        {t("profile.apiKeys.title")}
      </h1>
      <p className={classNames("mt-2 text-sm", SECONDARY_TEXT_COLOR_DIM)}>
        {t("profile.apiKeys.description")}{" "}
        <Link to={routes.apiDocs} className={LINK_COLOR}>
          {t("profile.apiKeys.docsLink")}
        </Link>
      </p>

      <form className="mt-6 flex gap-2 max-w-md" onSubmit={onCreate}>
        <Input
          type="text"
          value={name}
          maxLength={100}
          placeholder={t("profile.apiKeys.namePlaceholder")}
          aria-label={t("profile.apiKeys.namePlaceholder")}
          onChange={(e) => setName(e.target.value)}
        />
        <Button
          type="submit"
          disabled={creating || !name.trim()}
          className="shrink-0 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {t("profile.apiKeys.create")}
        </Button>
      </form>

      {newSecret && (
        <NewSecret secret={newSecret} onDone={() => setNewSecret(null)} />
      )}

      <div className="mt-6">
        {loading && !data ? null : apiKeys.length === 0 ? (
          <p className={classNames("text-sm", SECONDARY_TEXT_COLOR_DIM)}>
            {t("profile.apiKeys.empty")}
          </p>
        ) : (
          <ul className="space-y-2">
            {apiKeys.map((key) => (
              <li
                key={key.id}
                className={classNames(
                  "rounded-md p-3 flex items-center justify-between gap-4",
                  SECOND_LAYER,
                  TEXT_COLOR,
                )}
              >
                <div className="min-w-0">
                  <div className="font-medium truncate">{key.name}</div>
                  <div
                    className={classNames(
                      "text-xs mt-1 flex flex-wrap gap-x-4",
                      SECONDARY_TEXT_COLOR_DIM,
                    )}
                  >
                    <code className="font-mono">{key.keyPrefix}…</code>
                    <span>
                      {t("profile.apiKeys.createdAt", {
                        time: formatTime(key.createdAt),
                      })}
                    </span>
                    <span>
                      {key.lastUsedAt
                        ? t("profile.apiKeys.lastUsedAt", {
                            time: formatTime(key.lastUsedAt),
                          })
                        : t("profile.apiKeys.neverUsed")}
                    </span>
                  </div>
                </div>
                <Button
                  variant="secondary"
                  className="shrink-0 text-red-600 dark:text-red-400"
                  onClick={() => onRevoke(key.id, key.name)}
                >
                  {t("profile.apiKeys.revoke")}
                </Button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
