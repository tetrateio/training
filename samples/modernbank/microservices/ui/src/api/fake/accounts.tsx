interface IAccount {
    number: number;
    balance: number;
    owner: string;
}

export const fakeAccounts: IAccount[] = [
    {
        balance: 85000.23,
        number: 5201,
        owner: "jiajesse",
    },
    {
        balance: 5000.45,
        number: 5202,
        owner: "jiajesse",
    },
    {
        balance: 5000.45,
        number: 8501,
        owner: "liam",
    },
];
