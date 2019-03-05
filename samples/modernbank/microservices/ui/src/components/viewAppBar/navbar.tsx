import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import React from "react";
import {RouteComponentProps, withRouter} from "react-router";
import {AccountsPageLink, transactionsPageLink, transferPageLink} from "../../routes";

const height = 35;

const styles = (theme: Theme) => createStyles({
    button: {
        color: "white",
        textTransform: "none", /* Material button text defaults to upper case; disable it. */
    },
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    paper: {
        backgroundColor: "rgb(233,121,51)",
        height: `${height}px`,
    },
    searchIcon: {
        height: "100%",
    },
});

interface IUrlParams {
    accountNumber: string;
}

interface IProps extends WithStyles<typeof styles>, RouteComponentProps<IUrlParams> {
}

export const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    return (
        <Paper square={true} className={props.classes.paper}>
            <Grid
                container={true}
                alignItems="center"
                justify="space-around"
                className={props.classes.gridContainer}
            >
                <Grid item={true}>
                    <Button component={AccountsPageLink} className={props.classes.button}>
                        Accounts
                    </Button>
                </Grid>
                <Grid item={true}>
                    <Button
                        component={transferPageLink(props.match.params.accountNumber)}
                        className={props.classes.button}
                    >
                        Send money
                    </Button>
                </Grid>
                <Grid item={true}>
                    <Button
                        component={transactionsPageLink(props.match.params.accountNumber)}
                        className={props.classes.button}
                    >
                        Transactions
                    </Button>
                </Grid>
            </Grid>
        </Paper>
    );
};

const RoutingComponent = withRouter(Component);

export const Navbar = withStyles(styles)(RoutingComponent);
