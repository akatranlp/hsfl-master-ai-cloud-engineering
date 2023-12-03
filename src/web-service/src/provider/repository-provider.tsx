import { UserRepository } from "@/repository/user";
import React, { createContext, useContext, useState } from "react";
import { useApiClient } from "./api-client-provider";

type RepositoryProviderProps = {
  children: React.ReactNode;
};

type RepositoryProviderState = {
  userRepo: UserRepository;
};

const RepositoryContext = createContext<RepositoryProviderState | undefined>(undefined);

export const RepositoryProvider: React.FC<RepositoryProviderProps> = ({ children }) => {
  const { apiClient } = useApiClient();

  const [userRepo] = useState(() => new UserRepository(apiClient));

  return <RepositoryContext.Provider value={{ userRepo }}>{children}</RepositoryContext.Provider>;
};

export const useRepository = () => {
  const context = useContext(RepositoryContext);

  if (context === undefined) throw new Error("useUserRepository must be used within a UserRepositoryProvider");

  return context;
};
