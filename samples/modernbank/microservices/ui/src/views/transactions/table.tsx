import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Paper from "@material-ui/core/Paper";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {Newtransaction, Transaction, TransactionsApi} from "../../api/client";
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
    table: {

    }
});

interface IUrlParms {
    accountNumber: string;
}

interface IProps extends WithStyles<typeof styles>, RouteComponentProps<IUrlParms> {
}

export const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    const accountNumber = parseInt(props.match.params.accountNumber, 10);

    const [receivedTransactions, setReceivedTxs] = React.useState<Newtransaction[]>([]);
    const [sentTransactions, setSentTxs] = React.useState<Newtransaction[]>([]);

    const [doFetch, setDoFetch] = React.useState<boolean>(true);

    const transactionsApi = new TransactionsApi({basePath: "http://35.192.59.252/v1"});

    const fetchTransactionsReceived = async () => {
        const received: any = await transactionsApi.listTransactionsReceived(accountNumber);
        setReceivedTxs(received);
    };

    const fetchTransactionsSent = async () => {
        const sent: any = await transactionsApi.listTransactionsSent(accountNumber);
        setSentTxs(sent);
    };

    React.useEffect(() => {
        if (doFetch) {
            setDoFetch(false);
            fetchTransactionsReceived();
            fetchTransactionsSent();
        }
    });

    return (
        <Paper square={true} className={props.classes.paper}>
            <Table className={props.classes.table}>
                <TableHead>
                    <TableRow>
                        {/*<TableCell>Date</TableCell>*/}
                        <TableCell align="left">From</TableCell>
                        <TableCell align="left">Deposits</TableCell>
                        <TableCell align="left">To</TableCell>
                        <TableCell align="left">Withdrawals</TableCell>
                        <TableCell align="left">Balance</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {receivedTransactions.map((transaction) => (
                        <TableRow>
                            <TableCell
                                component="th"
                                scope="row"
                                align="left"
                            >
                                <Typography>{transaction.sender}</Typography>
                            </TableCell>
                            <TableCell align="right">{transaction.amount}</TableCell>
                            <TableCell></TableCell>
                            <TableCell></TableCell>
                            <TableCell align="right"></TableCell>
                        </TableRow>
                    ))}
                    {sentTransactions.map((transaction) => (
                        <TableRow>
                            <TableCell></TableCell>
                            <TableCell></TableCell>
                            <TableCell
                                component="th"
                                scope="row"
                                align="left"
                            >
                                <Typography>{transaction.receiver}</Typography>
                            </TableCell>
                            <TableCell align="right">{transaction.amount}</TableCell>
                            <TableCell align="right"></TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </Paper>
    );
};

const RoutingComponent = withRouter(Component);

export const TransactionsTable = withStyles(styles)(RoutingComponent);
