import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Hidden from "@material-ui/core/Hidden";
import React from "react";
import {Shell} from "../../components/shell";
import {InfoPanel} from "./infoPanel";
import {LoginForm} from "./loginForm";

const styles = () => createStyles({
    gridContainer: {
        borderTopColor: "rgb(172,37,45)",
        borderTopStyle: "solid",
        borderTopWidth: "0.6vh",
        height: "100%",
        width: "100%",
    },
    mdUpLoginForm: {
        paddingLeft: "2vw",
    },
});

interface IProps extends WithStyles<typeof styles> {
}

const Component: React.FunctionComponent<IProps> = (props: IProps) => (
    <Grid
        container={true}
        alignItems={"stretch"}
        justify={"space-between"}
        className={props.classes.gridContainer}
    >
        <Hidden smDown={true}>
            <Grid item={true} xs={5} className={props.classes.mdUpLoginForm}>
                <LoginForm/>
            </Grid>
            <Grid item={true} xs={6}>
                <InfoPanel/>
            </Grid>
        </Hidden>
        <Hidden mdUp={true}>
            <Grid item={true} xs={12}>
                <LoginForm/>
            </Grid>
        </Hidden>
    </Grid>
);

const StyledComponent = withStyles(styles)(Component);

export const LoginView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell>
        <StyledComponent/>
    </Shell>
);
