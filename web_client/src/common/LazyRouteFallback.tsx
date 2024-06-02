import React from "react";
import { useNavigation } from "react-router-dom";
import { FullScreenLoader } from "~src/ui/fullscreenLoader";

export function LazyRouteFallback() {
  // Fallback component to show a loading spinner while a lazy route is loading
  // Intended to be used in a layout component
  const navigation = useNavigation();
  if (navigation.state === "loading") {
    return <FullScreenLoader />;
  }
  return null;
}
