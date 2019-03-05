import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import React from "react";
import Paper from "@material-ui/core/Paper";
import {AccountBalance} from "@material-ui/icons";
import Typography from "@material-ui/core/Typography";

const borderBottomWidth = 3;
const height = 50;

const styles = (theme: Theme) => createStyles({
    accountBalanceLabel: {
        textAlign: "right",
    },
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    leftSide: {
        display: "flex",
        alignItems: "center",
    },
    paper: {
        borderBottom: `${borderBottomWidth}px solid rgb(196,196,196)`,
        height: `${height}px`,
        paddingLeft: 4 * theme.spacing.unit,
        paddingRight: 4 * theme.spacing.unit,
    },
    spacer: {
        width: 2 * theme.spacing.unit,
    },
});

interface IProps extends WithStyles<typeof styles> {
    accountBalance: number;
    accountNumber: number;
}

export const component: React.FunctionComponent<IProps> = (props: IProps) => {
    return (
        <Paper square={true} className={props.classes.paper}>
            <Grid
                container={true}
                alignItems={"center"}
                justify={"space-between"}
                className={props.classes.gridContainer}
            >
                <Grid item={true} xs={1}>
                    <AccountBalance/>
                </Grid>
                <Grid item={true} xs={3}>
                    <Typography variant={"subtitle2"}>Account</Typography>
                    <Typography variant={"subtitle2"}>{props.accountNumber}</Typography>
                </Grid>
                <Grid item={true} xs={6}/>
                <Grid item={true} xs={2}>
                    <Typography
                        variant={"subtitle2"}
                        className={props.classes.accountBalanceLabel}
                    >
                        {`$${props.accountBalance.toFixed(2)}`}
                    </Typography>
                    <Typography variant={"subtitle2"}>Available balance</Typography>
                </Grid>
            </Grid>
        </Paper>
    );
};

export const Subheader = withStyles(styles)(component);
