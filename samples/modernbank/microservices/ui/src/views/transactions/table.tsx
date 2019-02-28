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
import {fakeTransactions} from "../../api/fake/transactions";

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

interface IProps extends WithStyles<typeof styles> {
}

const fetchedAccounts: number[] = [
    1001, 1002, 1003,
];

interface IFormState {
    fromAccount: string;
    toAccount: string;
    routingNumber: string;
    date: string;
    amount: string;
}

export const component: React.FunctionComponent<IProps> = (props: IProps) => {
    return (
        <Paper square={true} className={props.classes.paper}>
            <Table className={props.classes.table}>
                <TableHead>
                    <TableRow>
                        <TableCell>Date</TableCell>
                        <TableCell align="right">Description</TableCell>
                        <TableCell align="right">Deposits</TableCell>
                        <TableCell align="right">Withdrawals</TableCell>
                        <TableCell align="right">Balance</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {fakeTransactions.map((transaction) => (
                        <TableRow key={transaction.id}>
                            <TableCell component="th" scope="row">
                            </TableCell>
                            <TableCell align="left">
                                <Typography>transaction.name</Typography>
                            </TableCell>
                            <TableCell align="right">{transaction.amount}</TableCell>
                            <TableCell align="right"></TableCell>
                            <TableCell align="right">"10.00"</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </Paper>
    );
};

export const TransactionsTable = withStyles(styles)(component);
