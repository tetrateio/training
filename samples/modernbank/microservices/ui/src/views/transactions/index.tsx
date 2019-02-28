import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {Shell} from "../../components/shell";
import {Navbar} from "../../components/viewAppBar/navbar";
import {Subheader} from "../../components/viewAppBar/subheader";
import {TransactionsTable} from "./table";

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

interface IProps extends WithStyles<typeof styles> {
}

export const Component: React.FunctionComponent<IProps> = (props: IProps) => (
    <>
        <div className={props.classes.subheader}>
            <Typography variant="h6" className={props.classes.subheaderText}>
                Account summary / Checking account
            </Typography>
        </div>
        <Navbar/>
        <Subheader/>
        <TransactionsTable/>
    </>
);

const StyledComponent = withStyles(styles)(Component);

export const TransactionsView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell showRightPanel={true}>
        <StyledComponent/>
    </Shell>
);
