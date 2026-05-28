export const routeSegments = {
  siteAdmin: "site-admin",
  profile: "profile",
  users: "users",
  images: "images",
  storage: "storage",
} as const;

export const routes = {
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
  },
};
