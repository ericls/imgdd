import React from "react";
import { useParams } from "react-router";
import { DEFAULT_MENU_CONFIG } from "~src/common/ImageGallery/menu";
import { ImageGallery } from "~src/common/ImageGallery/render";
import { useTranslation } from "react-i18next";
import { useQuery } from "@apollo/client/react";
import { gql } from "~src/__generated__/gql";

const OrgUserByIdDoc = gql(`
  query OrgUserById($id: ID!) {
    viewer {
      id
      organizationUserById(id: $id) {
        id
        user {
          id
          name
          avatarUrl
        }
      }
    }
  }
`);

export function UserImageGallery() {
  const { userId } = useParams();
  const { t } = useTranslation();
  const { data } = useQuery(OrgUserByIdDoc, {
    variables: { id: userId! },
    skip: !userId,
  });
  const orgUser = data?.viewer.organizationUserById;

  return (
    <>
      <div className="flex items-center gap-3 mb-4">
        {orgUser && (
          <img
            src={orgUser.user.avatarUrl}
            alt=""
            className="w-10 h-10 rounded-full object-cover flex-shrink-0"
          />
        )}
        <h2 className="text-2xl font-bold">
          {orgUser ? orgUser.user.name : t("userImageGallery.title")}
        </h2>
      </div>
      <div className="p-4">
        <ImageGallery createdById={userId} menuConfig={DEFAULT_MENU_CONFIG} />
      </div>
    </>
  );
}
