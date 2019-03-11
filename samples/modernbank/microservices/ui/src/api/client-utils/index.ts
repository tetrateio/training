import {AccountsApi, User, UsersApi} from "../client";

// TODO(jiajesse): Figure out what to set this to.
const basePath = "http://35.197.239.230/v1";

export const usersApi = new UsersApi({basePath});
export const accountsApi = new AccountsApi({basePath});

export const authenticationCheck = async (username: string, password: string): Promise<User> => {
    const options = {
        method: "GET",
    };
    const user: User = await usersApi.getUserByUserName(username, options);
    return (username === user.username && password === user.password) ? user : null;
};


