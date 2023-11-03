type User = {
  id: number;
  email: string;
  password: string;
  username: string;
  profileName: string;
  balance: number;
};

type UpdateUser = {
  profileName: string;
  balance: number;
};
