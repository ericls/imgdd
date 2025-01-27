import React from "react";
import cx from "classnames";
import { HiLockClosed as LockClosedIcon } from "react-icons/hi2";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { Button } from "~src/ui/button";
import { HEADING_1, LINK_COLOR, TEXT_COLOR } from "~src/ui/classNames";
import { Input } from "~src/ui/input";
import { useForm } from "react-hook-form";
import { gql } from "~src/__generated__";
import { useMutation } from "@apollo/client";
import { toast } from "react-toastify";
import { Loader } from "~src/ui/loader";

type ForgotPasswordFormData = {
  email: string;
};

const sendResetPasswordEmailMutationDocument = gql(`
mutation sendResetPasswordEmail($input: SendResetPasswordEmailInput!) {
  sendResetPasswordEmail(input: $input) {
    success
  }
}
`);

export function ForgotPasswordPage() {
  const { t } = useTranslation();
  const { register, getValues, reset } = useForm<ForgotPasswordFormData>();
  const [sendResetPasswordEmail, { loading }] = useMutation(
    sendResetPasswordEmailMutationDocument,
  );
  const onSubmit = React.useCallback(() => {
    const { email } = getValues();
    sendResetPasswordEmail({
      variables: {
        input: {
          email,
        },
      },
    })
      .then(() => {
        toast.success(t("auth.forgotPassword.success"));
      })
      .finally(() => {
        reset();
      });
  }, [getValues, reset, sendResetPasswordEmail, t]);
  return (
    <div className="max-w-md flex flex-col mx-auto">
      <div>
        <h2 className={cx("mt-6 text-center", HEADING_1)}>
          {t("auth.forgotPassword.title")}
        </h2>
        <p className={cx("mt-2 text-center text-sm ", TEXT_COLOR)}>
          {t("auth.forgotPassword.description")}
        </p>
      </div>
      <form className="mt-8 space-y-6">
        <div className="-space-y-px rounded-md shadow-sm">
          <div>
            <label
              htmlFor="user-identifier"
              className="sr-only dark:text-white"
            >
              {t("auth.fieldEmail")}
            </label>
            <Input
              id="user-identifier"
              type="email"
              required
              className="relative block w-full mb-2"
              placeholder={t("auth.fieldEmailPlaceholder")}
              {...register("email")}
            />
          </div>
        </div>

        <div className="flex items-center justify-end">
          <div className="text-sm">
            <Link to="/auth" className={cx("font-medium", LINK_COLOR)}>
              {t("auth.forgotPassword.backToSignInLink")}
            </Link>
          </div>
        </div>

        <div>
          <Button
            type="button"
            className="relative flex w-full"
            onClick={onSubmit}
            disabled={loading}
          >
            <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-indigo-500 group-hover:text-indigo-400">
              {loading ? (
                <Loader size={20} />
              ) : (
                <LockClosedIcon
                  size={20}
                  className="h-5 w-5"
                  aria-hidden="true"
                />
              )}
            </span>
            {t("auth.buttonSubmit")}
          </Button>
        </div>
      </form>
    </div>
  );
}
