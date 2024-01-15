import { DBResetButton } from "@/components/db-reset-button";
import { ApiClientProvider } from "@/provider/api-client-provider";
import { RepositoryProvider } from "@/provider/repository-provider";
import { Outlet, useLocation } from "react-router-dom";

export const RootLayout = () => {
  const location = useLocation();

  return (
    <ApiClientProvider>
      <RepositoryProvider>
        {(location.pathname === "/login" || location.pathname === "/register") && (
          <div className="flex items-center justify-center mt-2">
            <DBResetButton />
          </div>
        )}
        <Outlet />
      </RepositoryProvider>
    </ApiClientProvider>
  );
};
