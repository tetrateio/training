import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {RouteComponentProps, withRouter} from "react-router";
import {Account, AccountsApi} from "../../api/client";
import {AuthContext} from "../../components/auth/authContext";
import {Shell} from "../../components/shell";
import {Navbar} from "../../components/viewAppBar/navbar";
import {Subheader} from "../../components/viewAppBar/subheader";
import {TransactionsTable} from "./table";

const styles = () => createStyles({
    fillerPaper: {
        backgroundColor: "rgba(255,255,255,0.95)",
        boxShadow: "none",
        height: "100%",
    },
    fillerPaperGridItem: {
        height: "100%",
    },
    gridContainer: {
        flexWrap: "nowrap",
        height: "100%",
    },
    subheader: {
        backgroundColor: "rgb(172,37,45)",
    },
    subheaderText: {
        color: "white",
        marginLeft: "5vh",
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
        <Grid
            container={true}
            alignItems={"stretch"}
            direction={"column"}
            justify={"flex-start"}
            className={props.classes.gridContainer}
        >
            <Grid item={true}>
                <div className={props.classes.subheader}>
                    <Typography variant="h6" className={props.classes.subheaderText}>
                        Account summary / Transactions
                    </Typography>
                </div>
            </Grid>
            <Grid item={true}>
                <Navbar/>
            </Grid>
            <Grid item={true}>
                <Subheader accountBalance={account.balance} accountNumber={account.number}/>
            </Grid>
            <Grid item={true}>
                <TransactionsTable/>
            </Grid>
            <Grid item={true} className={props.classes.fillerPaperGridItem}>
                <Paper square={true} className={props.classes.fillerPaper}/>
            </Grid>
        </Grid>
    );
};

const RoutingComponent = withRouter(Component);

const StyledComponent = withStyles(styles)(RoutingComponent);

export const TransactionsView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell showRightPanel={true}>
        <StyledComponent {...props}/>
    </Shell>
);
