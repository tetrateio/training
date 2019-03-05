import {fakeUsers, IUser} from "../../api/fake/users";
import {User, UsersApi} from "../../api/client";

export const fakeSignInCheck = (username: string, password: string): boolean => {
    for (const user of fakeUsers) {
        if (user.username === username && user.password === password) {
            return true;
        }
    }
    return false;
}

export const authenticationCheck = async (username: string, password: string) => {
    const options = {
        method: "GET",
        // TODO(jiajesse): Delete this after manual testing.
        headers: {
            "Access-Control-Allow-Origin": "*",
        },
    };
    const usersApi = new UsersApi({basePath: "http://35.192.59.252/v1"});
    const user: User = await usersApi.getUserByUserName(username, options);
    return (username === user.username && password === user.password) ? user : null;

};
