import { TokenGenerator, createApiClient } from "@/lib/api-client";
import axios, { AxiosInstance } from "axios";
import React, { createContext, useCallback, useContext, useState } from "react";
import { useNavigate } from "react-router-dom";

type ApiClientProviderProps = {
  children: React.ReactNode;
};

type ApiClientProviderState = {
  apiClient: AxiosInstance;
};

const ApiClientContext = createContext<ApiClientProviderState | undefined>(undefined);

export const ApiClientProvider: React.FC<ApiClientProviderProps> = ({ children }) => {
  const navigate = useNavigate();
  const generateToken = useCallback(async () => {
    try {
      console.log("Get token");
      const resp = await axios.post("/api/v1/refresh-token");
      return resp.data.token_type === "bearer" ? resp.data.access_token : null;
    } catch (e) {
      navigate("/login");
      return null;
    }
  }, [navigate]);
  const [apiClient] = useState(() => createApiClient("/api/v1", new TokenGenerator(generateToken)));

  return <ApiClientContext.Provider value={{ apiClient }}>{children}</ApiClientContext.Provider>;
};

export const useApiClient = () => {
  const context = useContext(ApiClientContext);

  if (context === undefined) throw new Error("useApiClient must be used within a ApiClientProvider");

  return context;
};
