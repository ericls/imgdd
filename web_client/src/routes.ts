export const routeSegments = {
  siteAdmin: "site-admin",
  profile: "profile",
  users: "users",
  images: "images",
  storage: "storage",
  apiKeys: "api-keys",
} as const;

export const routes = {
  image: (imageId: string) => `/${routeSegments.images}/${imageId}`,
  apiDocs: "/docs/api",
  siteAdmin: {
    root: `/${routeSegments.siteAdmin}`,
    images: `/${routeSegments.siteAdmin}/${routeSegments.images}`,
    users: `/${routeSegments.siteAdmin}/${routeSegments.users}`,
    userImages: (orgUserId: string) =>
      `/${routeSegments.siteAdmin}/${routeSegments.users}/${orgUserId}/${routeSegments.images}`,
  },
  profile: {
    root: `/${routeSegments.profile}`,
    images: `/${routeSegments.profile}/${routeSegments.images}`,
    apiKeys: `/${routeSegments.profile}/${routeSegments.apiKeys}`,
    image: (imageId: string) =>
      `/${routeSegments.profile}/${routeSegments.images}/${imageId}`,
    editImage: (imageId: string) =>
      `/${routeSegments.profile}/${routeSegments.images}/${imageId}/edit`,
  },
};
