import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Divider from "@material-ui/core/Divider";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {AuthContext} from "../../components/auth/authContext";
import {Shell} from "../../components/shell";
import {AccountSummary} from "./accountSummary";
import {CashAccounts} from "./cashAccounts";
import {CreditAccounts} from "./creditAccounts";
import {InvestmentAccounts} from "./investmentAccounts";
import {TotalAccountValue} from "./totalAccountValue";
import {Account, AccountsApi} from "../../api/client";

const styles = (theme: Theme) => createStyles({
    paper: {
        opacity: 0.5,
    },
    subheader: {
        backgroundColor: "rgba(172,37,45, 1)",
    },
    subheaderText: {
        color: "white",
        marginLeft: "30px",
    },
});

interface IProps extends WithStyles<typeof styles> {
}

const accountsApi = new AccountsApi({basePath: "http://35.192.59.252/v1"});

const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    const [accounts, setAccounts] = React.useState<Account[]>([]);
    const authContext = React.useContext(AuthContext);
    const [doFetch, setDoFetch] = React.useState<boolean>(true);

    const fetchAccounts = async () => {
        const resp: Account[] = await accountsApi.listAccounts(authContext.user!.username);
        setAccounts(resp);
        console.log("jiajesse - accounts");
    };

    React.useEffect(() => {
        setDoFetch(false);
        fetchAccounts();
    }, []);

    // The API doesn't support account type. Fake account type using the last digit of the account number.
    const filterByLastDigit = (start: number, end: number): Account[] => {
        return accounts.filter((acc) => (start <= (acc.number % 10) && (acc.number % 10) <= end));
    };

    const cashAccounts = filterByLastDigit(0, 3);
    const investmentAccounts = filterByLastDigit(4, 6);
    const creditAccounts = filterByLastDigit(7, 9);

    return (
        <>
            <div className={props.classes.subheader}>
                <Typography variant="h6" className={props.classes.subheaderText}>
                    Account summary / Checking account
                </Typography>
            </div>
            <AccountSummary accounts={accounts}/>
            <CashAccounts accounts={cashAccounts}/>
            <InvestmentAccounts accounts={investmentAccounts}/>
            <CreditAccounts accounts={creditAccounts}/>
            <Divider/>
            <TotalAccountValue plusAccounts={cashAccounts.concat(investmentAccounts)} minusAccounts={creditAccounts}/>
        </>
    );
};

const StyledComponent = withStyles(styles)(Component);

export const AccountsView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell showRightPanel={true}>
        <StyledComponent/>
    </Shell>
);
