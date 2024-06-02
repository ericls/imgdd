import React from "react";
import cx from "classnames";
import { HiLockClosed as LockClosedIcon } from "react-icons/hi2";

import { Input } from "~src/ui/input";
import { Button } from "~src/ui/button";
import { HEADING_1, LINK_COLOR, TEXT_COLOR } from "~src/ui/classNames";
import { useForm } from "react-hook-form";
import { useApolloClient, useMutation } from "@apollo/client";
import { toast } from "react-toastify";
import { Link, useNavigate } from "react-router-dom";
import { gql } from "~src/__generated__/gql";
import { Trans, useTranslation } from "react-i18next";

type AuthFormData = {
  email: string;
  password: string;
  password2: string;
};

const createUserMutationDocument = gql(`
mutation createUserWithOrganization($input: CreateUserWithOrganizationInput!) {
  createUserWithOrganization(
    input: $input
  ) {
    viewer {
      id
      organizationUser {
        id
        user {
          id
          email
          name
        }
      }
    }
  }
}
`);

const authenticateMutationDocument = gql(`
mutation authenticate($email: String!, $password: String!) {
  authenticate(email: $email, password: $password) {
    viewer {
      id
      organizationUser {
        id
        user {
          id
          email
          name
        }
      }
    }
  }
}
`);

export function AuthPage() {
  const { t } = useTranslation();
  const [mode, setMode] = React.useState<"login" | "reg">("login");
  const apolloClient = useApolloClient();
  const { register, getValues, reset } = useForm<AuthFormData>();
  const navigate = useNavigate();
  const setLogin = React.useCallback(() => {
    setMode("login");
    reset();
  }, [setMode, reset]);
  const setReg = React.useCallback(() => {
    setMode("reg");
    reset();
  }, [setMode, reset]);
  const [authenticate] = useMutation(authenticateMutationDocument);
  const [createUser] = useMutation(createUserMutationDocument);
  const onLogin = React.useCallback(() => {
    const { email, password } = getValues();
    authenticate({ variables: { email, password } })
      .then((res) => {
        if (res.data?.authenticate.viewer.organizationUser?.id) {
          apolloClient.refetchQueries({ include: "active" }).then(() => {
            navigate("/");
          });
        } else {
          throw new Error();
        }
      })
      .catch(() => {
        toast.error(t("auth.authenticationFailed"));
      });
  }, [getValues, authenticate, apolloClient, navigate, t]);
  const onCreateUser = React.useCallback(() => {
    const { email, password, password2 } = getValues();
    if (password !== password2) {
      toast.error(t("auth.passwordsDoNotMatch"));
      return;
    }
    createUser({
      variables: {
        input: {
          userEmail: email,
          userPassword: password,
          organizationName: "",
        },
      },
    })
      .then((res) => {
        if (res.data?.createUserWithOrganization.viewer.organizationUser?.id) {
          apolloClient.refetchQueries({ include: ["auth"] }).then(() => {
            navigate("/");
          });
        } else {
          throw new Error("User creation failed.");
        }
      })
      .catch(() => {
        toast.error(t("auth.userCreationFailed"));
      });
  }, [getValues, createUser, t, apolloClient, navigate]);
  const onSubmit = React.useCallback(() => {
    if (mode === "login") {
      onLogin();
    } else if (mode === "reg") {
      onCreateUser();
    }
  }, [mode, onLogin, onCreateUser]);
  return (
    <div className="max-w-md flex flex-col mx-auto">
      <div>
        <h2 className={cx("mt-6 text-center", HEADING_1)}>
          {mode === "login"
            ? t("auth.titleSignIn")
            : t("auth.titleCreateAccount")}
        </h2>
        {mode === "login" ? (
          <p className={cx("mt-2 text-center text-sm ", TEXT_COLOR)}>
            <Trans i18nKey={"auth.orCreateAccount"}>
              Or{" "}
              <a
                href="#"
                onClick={setReg}
                className={cx("font-medium", LINK_COLOR)}
              >
                create an account
              </a>
            </Trans>
          </p>
        ) : (
          <p className={cx("mt-2 text-center text-sm ", TEXT_COLOR)}>
            <Trans i18nKey="auth.orSignIn">
              Or{" "}
              <a
                href="#"
                onClick={setLogin}
                className={cx("font-medium", LINK_COLOR)}
              >
                sign in
              </a>
            </Trans>
          </p>
        )}
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
          <div>
            <label htmlFor="password" className="sr-only dark:text-white">
              {t("auth.fieldPassword")}
            </label>
            <Input
              id="password"
              type="password"
              required
              className="relative block w-full"
              placeholder={t("auth.fieldPasswordPlaceholder")}
              {...register("password")}
            />
          </div>
          {mode === "reg" && (
            <div>
              <label htmlFor="password2" className="sr-only dark:text-white">
                {t("auth.fieldRepeatPassword")}
              </label>
              <Input
                id="password2"
                type="password"
                required
                className="relative block w-full mt-2"
                placeholder={t("auth.fieldRepeatPasswordPlaceholder")}
                {...register("password2")}
              />
            </div>
          )}
        </div>

        <div className="flex items-center justify-end">
          <div className="text-sm">
            <Link to="/auth/forgot" className={cx("font-medium", LINK_COLOR)}>
              {t("auth.forgotPasswordLink")}
            </Link>
          </div>
        </div>

        <div>
          <Button
            type="button"
            className="relative flex w-full"
            onClick={onSubmit}
          >
            <span className="absolute inset-y-0 left-0 flex items-center pl-3">
              <LockClosedIcon
                size={20}
                className="h-5 w-5 text-indigo-500 group-hover:text-indigo-400"
                aria-hidden="true"
              />
            </span>
            {mode == "login"
              ? t("auth.buttonSignIn")
              : t("auth.buttonCreateAccount")}
          </Button>
        </div>
      </form>
    </div>
  );
}
