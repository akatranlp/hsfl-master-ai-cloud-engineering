import { ApiClientProvider } from "@/provider/api-client-provider";
import { RepositoryProvider } from "@/provider/repository-provider";
import { Outlet } from "react-router-dom";

export const RootLayout = () => {
  return (
    <ApiClientProvider>
      <RepositoryProvider>
        <Outlet />
      </RepositoryProvider>
    </ApiClientProvider>
  );
};
