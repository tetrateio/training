export interface IUser {
    username: string,
    firstName: string,
    lastName: string;
    email: string;
    password: string;
}

export const fakeUsers: IUser[] = [
    {
        username: "a",
        firstName: "Jesse",
        lastName: "Jiang",
        email: "jesse@tetrate.io",
        password: "a",
    },
    {
        username: "jiajesse",
        firstName: "Jesse",
        lastName: "Jiang",
        email: "jesse@tetrate.io",
        password: "tetrate",
    },
    {
        username: "liam",
        firstName: "Liam",
        lastName: "White",
        email: "liam@tetrate.io",
        password: "tetrate",
    },
];
