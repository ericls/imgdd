import { useMutation } from "@apollo/client";
import React from "react";
import { gql } from "~src/__generated__";
import { HiOutlineQuestionMarkCircle as UnknownIcon } from "react-icons/hi2";
import { HiOutlineCheck as SuccessIcon } from "react-icons/hi2";
import { HiOutlineExclamationCircle as ErrorIcon } from "react-icons/hi2";
import { useTranslation } from "react-i18next";
import { Button } from "~src/ui/button";
import { notice } from "~src/ui/prompt";

const connectivityMutation = gql(`
mutation StorageDefTableConnectivityCellMutation(
    $input: checkStorageDefinitionConnectivityInput!
  ) {
    checkStorageDefinitionConnectivity(input: $input) {
      ok
      error
    }
  }
`);

export function StorageDefTableConnectivityCell({ id }: { id: string }) {
  const { t } = useTranslation();
  const [mutate, result] = useMutation(connectivityMutation);
  const resultOk = result.data?.checkStorageDefinitionConnectivity?.ok;
  const buttonLabel =
    resultOk === undefined
      ? t("storageDef.connection.unknown")
      : resultOk
        ? t("storageDef.connection.success")
        : t("storageDef.connection.error");
  const onClick = React.useCallback(() => {
    mutate({
      variables: {
        input: {
          id,
        },
      },
    }).then((res) => {
      const errorStr = res.data?.checkStorageDefinitionConnectivity?.error;
      if (errorStr) {
        notice(t("storageDef.connection.error"), <>{errorStr}</>);
      }
    });
  }, [id, mutate, t]);
  const buttonIcon =
    resultOk === undefined ? (
      <UnknownIcon />
    ) : resultOk ? (
      <SuccessIcon />
    ) : (
      <ErrorIcon />
    );
  return (
    <span>
      <Button
        variant={
          resultOk === undefined
            ? "secondary"
            : resultOk
              ? "transparent"
              : "secondary"
        }
        onClick={onClick}
        disabled={result.loading || resultOk === true}
        title={buttonLabel}
      >
        {buttonIcon}
      </Button>
    </span>
  );
}
