import type { PropsWithChildren } from "react";
import { TopNavigation } from "./TopNavigation";

export const PageLayout = ({ children }: PropsWithChildren) => {
  return (
    <div className="w-screen h-screen flex flex-col items-center">
      <TopNavigation />
      <div className="mt-17 w-full max-w-7xl justify-center">{children}</div>
    </div>
  );
};
