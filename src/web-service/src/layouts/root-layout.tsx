import { ApiClientProvider } from "@/provider/api-client-provider";
import { Outlet } from "react-router-dom";

export const RootLayout = () => {
  return (
    <ApiClientProvider>
      <Outlet />
    </ApiClientProvider>
  );
};
