import type { ReactNode } from "react";
import { useLocation } from "react-router-dom";
import { TopNavigation } from "./TopNavigation";

export const PageLayout = ({ children }: { children: ReactNode }) => {
  const loc = useLocation();

  return (
    <div className="flex-center w-dvw h-dvh">
      {loc.pathname != "/login" && <TopNavigation />}
      {children}
    </div>
  );
};
