import { Link, Outlet } from "react-router-dom";
import { ModeToggle } from "@/components/mode-toggle.tsx";
import { UserDataProvider, useUserData } from "@/provider/user-provider.tsx";

const NavBar = () => {
  const user = useUserData();

  return (
    <nav>
      <ul className="flex items-center p-2 dark:bg-slate-900 bg-gray-200">
        <li>
          <Link to="/books" className="text-center px-16">
            All Books
          </Link>
        </li>
        <li>
          <Link to="/books/myBooks" className="text-center px-16">
            My Books
          </Link>
        </li>
        <li>
          <Link to="/books/boughtBooks" className="text-center px-16">
            My bought Books
          </Link>
        </li>
        <li>
          <Link to="/transactions" className="text-center px-16">
            My Transactions
          </Link>
        </li>
        <li className="ml-auto">
          <div className="px-16">{user.profileName}</div>
          <div className="text-sm px-16">{user.balance} VV-Coins</div>
        </li>
        <li>
          <ModeToggle />
        </li>
      </ul>
    </nav>
  );
};

export const RootLayout = () => {
  return (
    <UserDataProvider>
      <header>
        <NavBar />
      </header>
      <main>
        <div className="flex w-full justify-center">
          <div className="w-1/2">
            <Outlet />
          </div>
        </div>
      </main>
    </UserDataProvider>
  );
};
