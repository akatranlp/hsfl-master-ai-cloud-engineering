import axios from "axios";

export const getMe = async () => {
  const response = await axios.get<User>("/api/v1/users/me");
  return response.data;
};

export const addCoins = async (user: UpdateUser) => {
  const response = await axios.patch<void>("/api/v1/users/me", user);
  return response.data;
};

export const login = async (email: string, password: string) => {
  const response = await axios.post<void>("/api/v1/login", { email, password });
  return response.data;
};

export const register = async (email: string, password: string, profileName: string) => {
  const response = await axios.post<void>("/api/v1/register", { email, password, profileName });
  return response.data;
};

export const logout = async () => {
  const response = await axios.post<void>("/api/v1/logout");
  return response.data;
};
