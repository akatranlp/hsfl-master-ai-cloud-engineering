import React, { createContext, useContext, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { getMe } from "@/repository/user.ts";
import { Navigate } from "react-router-dom";

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
  // const [userr, setUser] = useState<User | null>(null);
  useState;

  const { data: user, isLoading } = useQuery({ queryKey: ["me"], queryFn: getMe });
  if (isLoading) return <div>Page is currently Loading...</div>;
  if (!user || user.id === 0) return <Navigate to="/login" />;
  return <UserDataContext.Provider value={user}>{children}</UserDataContext.Provider>;
};

export const useUserData = () => useContext(UserDataContext);
