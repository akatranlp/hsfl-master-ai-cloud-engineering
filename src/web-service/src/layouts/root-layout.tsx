import { Outlet } from "react-router-dom";
import { ModeToggle } from "@/components/mode-toggle.tsx";

export const RootLayout = () => {
  return (
    <div>
      <ModeToggle />
      <header>Hallo Layout</header>
      <main>
        <Outlet />
      </main>
    </div>
  );
};
