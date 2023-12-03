import { Link, Outlet, useNavigate } from "react-router-dom";
import { ModeToggle } from "@/components/mode-toggle.tsx";
import { UserDataProvider, useUserData } from "@/provider/user-provider.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { logout } from "@/repository/user.ts";
import { toast } from "react-hot-toast";

const LogoutButton = () => {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const { mutate } = useMutation({
    mutationFn: async () => {
      logout();
      queryClient.clear();
      await queryClient.invalidateQueries();
      await queryClient.invalidateQueries({ queryKey: ["me"] });
      await queryClient.resetQueries();
    },
    onSuccess: async () => {
      navigate("/login");
    },
    onError: () => {
      toast.error("An error occurred. Please try again.");
    },
  });

  const handleLogout = () => {
    mutate();
  };

  return (
    <Button variant="link" onClick={handleLogout}>
      Logout
    </Button>
  );
};

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
          <DropdownMenu>
            <DropdownMenuTrigger>
              <div className="px-16">{user.profileName}</div>
              <div className="text-sm px-16">{user.balance} VV-Coins</div>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuLabel>My Account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <Button variant="link">
                  <Link to="/user/addCoins">Add VV-Coins</Link>
                </Button>
              </DropdownMenuItem>
              <DropdownMenuItem>
                <LogoutButton />
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </li>
        <li>
          <ModeToggle />
        </li>
      </ul>
    </nav>
  );
};

export const MainLayout = () => {
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
