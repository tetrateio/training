import { AccountsApi, User, UsersApi } from '../client';

// TODO(dio): probably we can just declare this globally.
export const usersApi = new UsersApi();
export const accountsApi = new AccountsApi();

export const authenticationCheck = async (
  username: string,
  password: string,
  setVersion: Function
): Promise<User> => {
  const resp = await usersApi.getUserByUserNameRaw({ username });
  setVersion(resp.raw.headers.get('version'));
  const user = await resp.value();
  return username === user.username && password === user.password
    ? await user
    : null;
};
