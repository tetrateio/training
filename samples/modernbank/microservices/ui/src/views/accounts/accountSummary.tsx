import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {Account} from "../../api/client";

const styles = () => createStyles({
    gridContainer: {
        height: "100%",
    },
    headerText: {},
    paper: {
        backgroundColor: "rgba(255,255,255,0.95)",
        paddingBottom: "1vh",
        paddingLeft: "2vh",
        paddingRight: "2vh",
        paddingTop: "1vh",
    },
    table: {},
});

interface IProps extends WithStyles<typeof styles> {
    plusAccounts: Account[];
    minusAccounts: Account[];
}

export const component: React.FunctionComponent<IProps> = (props: IProps) => {
    const creditAmount =
        props.plusAccounts
            .map((account) => account.balance)
            .reduce((sum, b) => sum + b, 0);
    const debtAmount =
        props.minusAccounts
            .map((account) => account.balance)
            .reduce((sum, b) => sum + b, 0);
    const totalAmount = creditAmount - debtAmount;

    return (
        <Paper square={true} className={props.classes.paper}>
            <Grid
                container={true}
                alignItems={"center"}
                justify={"space-between"}
                className={props.classes.gridContainer}
            >
                <Grid item={true}>
                    <Typography variant="h6" className={props.classes.headerText}>
                        Account summary
                    </Typography>
                </Grid>
                <Grid item={true} xs={4}/>
                <Grid item={true}>
                    <Typography variant="body1">
                        {`$${totalAmount.toFixed(2)}`}
                    </Typography>
                    <Typography variant="body1">
                        Total available cash
                    </Typography>
                </Grid>
            </Grid>
            <Table className={props.classes.table}>
                <TableBody>
                    <TableRow key={"credit"}>
                        <TableCell component="th" scope="row">
                            Credit
                        </TableCell>
                        <TableCell align="right">{`$${creditAmount.toFixed(2)}`}</TableCell>
                    </TableRow>
                    <TableRow key={"debt"}>
                        <TableCell component="th" scope="row">
                            Debt
                        </TableCell>
                        <TableCell align="right">{`$${debtAmount.toFixed(2)}`}</TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        </Paper>
    );
}

export const AccountSummary = withStyles(styles)(component);
