import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import Grid from "@material-ui/core/Grid";
import React from "react";
import Paper from "@material-ui/core/Paper";
import {Search} from "@material-ui/icons";
import {AccountsPageLink, TransactionsPageLink, TransferPageLink} from "../../routes";

const height = 35;

const styles = (theme: Theme) => createStyles({
    button: {
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

interface IProps extends WithStyles<typeof styles> {
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
                <Grid item={true}>
                    <Button component={AccountsPageLink} className={props.classes.button}>
                        Accounts
                    </Button>
                </Grid>
                <Grid item={true}>
                    <Button component={TransferPageLink} className={props.classes.button}>
                        Transfer funds
                    </Button>
                </Grid>
                <Grid item={true}>
                    <Button className={props.classes.button}>
                        Pay bills
                    </Button>
                </Grid>
                <Grid item={true}>
                    <Button className={props.classes.button}>
                        Send money
                    </Button>
                </Grid>
                <Grid item={true}>
                    <Button className={props.classes.button}>
                        View statements
                    </Button>
                </Grid>
                <Grid item={true}>
                    <Button component={TransactionsPageLink} className={props.classes.button}>
                        Search transactions
                    </Button>
                    <Search className={props.classes.searchIcon}/>
                </Grid>
            </Grid>
        </Paper>
    );
};

export const Navbar = withStyles(styles)(component);
