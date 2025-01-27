import React from "react";
import cx from "classnames";
import { HiLockClosed as LockClosedIcon } from "react-icons/hi2";
import { useTranslation } from "react-i18next";
import { Button } from "~src/ui/button";
import { HEADING_1 } from "~src/ui/classNames";
import { Input } from "~src/ui/input";
import { useForm } from "react-hook-form";
import { gql } from "~src/__generated__";
import { useMutation } from "@apollo/client";
import { toast } from "react-toastify";
import { Loader } from "~src/ui/loader";
import { useNavigate } from "react-router-dom";

type ResetPasswordFormData = {
  message: string;
  newPassword: string;
  confirmPassword: string;
};

const resetPasswordDocument = gql(`
mutation resetPassword($input: ResetPasswordInput!) {
  resetPassword(input: $input) {
    success
  }
}
`);

export function ResetPasswordPage() {
  const { t } = useTranslation();
  const { register, getValues, reset, setValue } =
    useForm<ResetPasswordFormData>();
  const navigate = useNavigate();
  React.useEffect(() => {
    const searchParams = new URLSearchParams(window.location.search);
    const message = searchParams.get("message");
    if (message) {
      setValue("message", message, { shouldDirty: false });
    }
    window.history.replaceState({}, "", window.location.pathname);
  }, [setValue]);
  const [resetPassword, { loading }] = useMutation(resetPasswordDocument, {
    errorPolicy: "all",
  });
  const onSubmit = React.useCallback(() => {
    const { message, newPassword, confirmPassword } = getValues();
    if (newPassword !== confirmPassword) {
      toast.error(t("auth.passwordsDoNotMatch"));
      return;
    }
    if (newPassword === "") {
      toast.error(t("auth.passwordCannotBeEmpty"));
      return;
    }
    resetPassword({
      variables: {
        input: {
          message,
          password: newPassword,
        },
      },
    }).then((data) => {
      if (!data.data?.resetPassword.success) {
        toast.error(t("auth.resetPassword.error"));
        return;
      }
      toast.success(t("auth.resetPassword.success"));
      reset();
      navigate("/auth");
    });
  }, [getValues, navigate, reset, resetPassword, t]);
  return (
    <div className="max-w-md flex flex-col mx-auto">
      <div>
        <h2 className={cx("mt-6 text-center", HEADING_1)}>
          {t("auth.resetPassword.title")}
        </h2>
      </div>
      <form className="mt-8 space-y-6">
        <div className="-space-y-px rounded-md shadow-sm">
          <div className="hidden">
            <label htmlFor="reset-message" className="sr-only dark:text-white">
              {t("auth.fieldResetMessage")}
            </label>
            <Input
              id="reset-message"
              type="hidden"
              required
              className="relative block w-full mb-2"
              {...register("message")}
            />
          </div>
          <div>
            <label htmlFor="reset-np" className="sr-only dark:text-white">
              {t("auth.resetPassword.newPasswordField")}
            </label>
            <Input
              id="reset-np"
              type="password"
              required
              className="relative block w-full mb-2"
              placeholder={t("auth.resetPassword.newPasswordPlaceholder")}
              {...register("newPassword")}
            />
          </div>
          <div>
            <label htmlFor="reset-npc" className="sr-only dark:text-white">
              {t("auth.resetPassword.newPasswordConfirmField")}
            </label>
            <Input
              id="reset-npc"
              type="password"
              required
              className="relative block w-full mb-2"
              placeholder={t(
                "auth.resetPassword.newPasswordConfirmPlaceholder",
              )}
              {...register("confirmPassword")}
            />
          </div>
        </div>

        <div>
          <Button
            type="button"
            className="relative flex w-full"
            onClick={onSubmit}
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
