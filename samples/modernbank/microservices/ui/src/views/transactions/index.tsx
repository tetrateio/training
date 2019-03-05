import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {Shell} from "../../components/shell";
import {Navbar} from "../../components/viewAppBar/navbar";
import {Subheader} from "../../components/viewAppBar/subheader";
import {TransactionsTable} from "./table";
import {Account, AccountsApi} from "../../api/client";
import {AuthContext} from "../../components/auth/authContext";
import {RouteComponentProps, withRouter} from "react-router";

const styles = (theme: Theme) => createStyles({
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    paper: {
        backgroundColor: "rgba(255,255,255,0.97)",
        height: "100vh",
        paddingLeft: 2 * theme.spacing.unit,
        paddingRight: 2 * theme.spacing.unit,
    },
    subheader: {
        backgroundColor: "rgb(172,37,45)",
    },
    subheaderText: {
        color: "white",
        marginLeft: "30px",
    },
});

interface IUrlParams {
    accountNumber: string;
}

const accountsApi = new AccountsApi({basePath: "http://35.192.59.252/v1"});

interface IProps extends WithStyles<typeof styles>, RouteComponentProps<IUrlParams> {
}

const initialAccount: Account = {
    balance: 0,
    number: 0,
    owner: "",
};

export const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    const authContext = React.useContext(AuthContext);
    const [account, setAccount] = React.useState<Account>(initialAccount);

    const accountNumber = parseInt(props.match.params.accountNumber, 10);

    // TODO(jiajesse): GET .../accounts/{number} doesn't work. Use GET .../accounts and filter for now.
    // const fetchAccount = async () => {
    //     const resp: Account = await accountsApi.getAccountByNumber(authContext.user!.username, accountNumber);
    //     setAccount(resp);
    // };
    // React.useEffect(() => {
    //     fetchAccount();
    // }, []);

    const fetchAccounts = async () => {
        const resp: Account[] = await accountsApi.listAccounts(authContext.user!.username);
        const acc = resp.find((a) => a.number === accountNumber);
        setAccount(acc);
    };

    React.useEffect(() => {
        fetchAccounts();
    }, []);

    return (
        <>
            <div className={props.classes.subheader}>
                <Typography variant="h6" className={props.classes.subheaderText}>
                    Account summary / Checking account
                </Typography>
            </div>
            <Navbar/>
            <Subheader accountBalance={account.balance} accountNumber={account.number}/>
            <TransactionsTable/>
        </>
    );
};

const RoutingComponent = withRouter(Component);

const StyledComponent = withStyles(styles)(RoutingComponent);

export const TransactionsView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell showRightPanel={true}>
        <StyledComponent {...props}/>
    </Shell>
);
