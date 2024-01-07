import axios, { AxiosInstance } from "axios";

export class UserRepository {
  constructor(private apiClient: AxiosInstance) {}

  //This is how it should look (apiClient instead of axios)
  async getMe() {
    const response = await this.apiClient.get<User>("/users/me");
    return response.data;
  }

  async getUserById(userId: number) {
    const response = await this.apiClient.get<User>(`/users/${userId}`);
    return response.data;
  }

  async addCoins(user: UpdateUser) {
    const response = await this.apiClient.patch<void>("/users/me", user);
    return response.data;
  }

  async login(email: string, password: string) {
    const response = await axios.post<void>("/api/v1/login", { email, password }, { withCredentials: true });
    return response.data;
  }

  async register(email: string, password: string, profileName: string) {
    const response = await axios.post<void>("/api/v1/register", { email, password, profileName }, { withCredentials: true });
    return response.data;
  }

  async logout() {
    const response = await this.apiClient.post<void>("/logout");
    return response.data;
  }
}
