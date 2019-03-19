import { AccountsApi, User, UsersApi } from '../client';

// TODO(dio): probably we can just declare this globally.
export const usersApi = new UsersApi();
export const accountsApi = new AccountsApi();

export const authenticationCheck = async (
  username: string,
  password: string
): Promise<User> => {
  const user: User = await usersApi.getUserByUserName({ username });
  return username === user.username && password === user.password ? user : null;
};
