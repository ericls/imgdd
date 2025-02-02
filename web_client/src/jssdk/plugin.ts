export type IMGDDPlugin = {
  textSlots: {
    [key: string]: JSX.Element | (() => JSX.Element);
  };
};
