import axios from "axios";

export const getMe = async () => {
  const response = await axios.get<User>("/api/v1/users/me");
  return response.data;
};
