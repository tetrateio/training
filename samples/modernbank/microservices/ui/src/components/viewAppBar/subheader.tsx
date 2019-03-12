import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import {AccountBalance} from "@material-ui/icons";
import React from "react";

const borderBottomWidth = 3;
const height = 50;

const styles = () => createStyles({
    accountBalanceText: {
        paddingRight: "2vw",
        textAlign: "right",
    },
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    paper: {
        backgroundColor: "rgba(255,255,255,0.95)",
        borderBottom: `${borderBottomWidth}px solid rgb(196,196,196)`,
        height: `${height}px`,
        paddingLeft: "1vw",
        paddingRight: "1vw",
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
                        className={props.classes.accountBalanceText}
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
