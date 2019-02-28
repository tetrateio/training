interface ITransaction {
    id: string;
    sender: number;
    receiver: number;
    amount: number;
}

export const fakeTransactions: ITransaction[] = [
    {
        amount: 1000.53,
        id: "1",
        receiver: 5201,
        sender: 5202,
    },
    {
        amount: 30.32,
        id: "2",
        sender: 8501,
        receiver: 5201,
    },
    {
        amount: 400.53,
        id: "3",
        sender: 5201,
        receiver: 8501,
    },
];
