import {fakeUsers, IUser} from "../../api/fake/users";

export const fakeSignInCheck = (username: string, password: string): boolean => {
    for (const user of fakeUsers) {
        if (user.username === username && user.password === password) {
            return true;
        }
    }
    return false;
}
