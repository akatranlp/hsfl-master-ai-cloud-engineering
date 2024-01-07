import axios from "axios";

export class TokenGenerator {
  public token: string;
  constructor(public generateToken: () => Promise<string>) {
    this.token = "";
  }
}

const parseJwt = (token: string) => {
  if (!token) return undefined;
  const base64Url = token.split(".")[1];
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map((c) => {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join("")
  );

  return JSON.parse(jsonPayload);
};

const isTokenValid = (token: string) => {
  const parsedToken = parseJwt(token);
  if (!parsedToken) return false;
  return parsedToken.exp ? parsedToken.exp > Date.now() / 1000 : true;
};

export const createApiClient = (baseURL: string, generator: TokenGenerator) => {
  console.log("create Api Client");

  const axiosInstance = axios.create({
    baseURL,
    withCredentials: true,
    headers: { Authorization: `Bearer ${generator.token}` },
  });

  axiosInstance.interceptors.request.use(async (req) => {
    if (!isTokenValid(generator.token)) {
      generator.token = await generator.generateToken();
    }
    req.headers.Authorization = `Bearer ${generator.token}`;
    return req;
  });

  axiosInstance.interceptors.response.use((res) => {
    if ((res?.request?.responseURL as string | undefined)?.endsWith("/logout")) {
      generator.token = "";
    }
    return res;
  });

  return axiosInstance;
};
