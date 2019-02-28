import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {RightPanel} from "../../components/shell/rightPanel";

const styles = (theme: Theme) => createStyles({
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    headerText: {},
    paper: {
        padding: "5px 40px",
    },
    table: {},
});

interface IProps extends WithStyles<typeof styles> {
}

interface IRow {
    name: string;
    value: string;
}

const rows: IRow[] = [
    {
        name: "Investments",
        value: "$5000",
    },
    {
        name: "Debt",
        value: "$1000",
    },
    {
        name: "Total value",
        value: "$4000",
    },
];

// TODO: need to be consistent with API.
interface IAccount {
    number: number;
    balance: number;
    owner: string;
}



export const component: React.FunctionComponent<IProps> = (props: IProps) => {
    // const accounts = React.useState();

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
                        Total available cash
                    </Typography>
                </Grid>
            </Grid>
            <Table className={props.classes.table}>
                <TableBody>
                    {rows.map((row) => (
                        <TableRow key={row.name}>
                            <TableCell component="th" scope="row">
                                {row.name}
                            </TableCell>
                            <TableCell align="right">{row.value}</TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </Paper>
    );
}

export const AccountSummary = withStyles(styles)(component);
