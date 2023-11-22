import axios from "axios";

export const getMe = async () => {
  const response = await axios.get<User>("/api/v1/users/me");
  return response.data;
};

export const addCoins = async (user: UpdateUser) => {
  const response = await axios.patch<void>("/api/v1/users/me", user);
  return response.data;
};
