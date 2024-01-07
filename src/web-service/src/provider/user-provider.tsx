import React, { createContext, useContext } from "react";
import { useQuery } from "@tanstack/react-query";
import { Navigate } from "react-router-dom";
import { useRepository } from "./repository-provider";

type UserDataProviderProps = {
  children: React.ReactNode;
};

const initialState: User = {
  email: "",
  id: 0,
  balance: 0,
  profileName: "",
};

const UserDataContext = createContext<User>(initialState);
export const UserDataProvider: React.FC<UserDataProviderProps> = ({ children }) => {
  const { userRepo } = useRepository();

  const { data: user, isLoading } = useQuery({ queryKey: ["me"], queryFn: () => userRepo.getMe(), retry: false, refetchInterval: 10000 });
  if (isLoading) return <div>Page is currently Loading...</div>;
  if (!user || user.id === 0) return <Navigate to="/login" />;
  return <UserDataContext.Provider value={user}>{children}</UserDataContext.Provider>;
};

export const useUserData = () => {
  const context = useContext(UserDataContext);

  if (context === undefined) throw new Error("useUserData must be used within a UserDataProvider");

  return context;
};
