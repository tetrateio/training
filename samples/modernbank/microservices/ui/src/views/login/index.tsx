import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import React from "react";
import {Shell} from "../../components/shell";
import {InfoPanel} from "./infoPanel";
import {LoginForm} from "./loginForm";
import Hidden from "@material-ui/core/Hidden";

const styles = () => createStyles({
    gridContainer: {
        borderTopColor: "rgb(172,37,45)",
        borderTopStyle: "solid",
        borderTopWidth: "0.5vh",
        height: "100%", /* Force the grid to be same size as parent Paper component. */
        width: "100%",
    },
    mdUpLoginForm: {
        paddingLeft: "2vw",
    },
    paper: {
        backgroundColor: "white",
        height: "50%",
        // opacity: 0.5,
        paddingLeft: "2vw",
        paddingRight: "2vw",
        paddingTop: "3vh",
        width: "100%",
    },
    subheader: {
        backgroundColor: "rgba(172,37,45, 1)",
    },
    subheaderText: {
        color: "white",
        marginLeft: "30px",
    },
    textField: {
        // marginLeft: theme.spacing.unit,
        // marginRight: theme.spacing.unit,
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
            <Grid item={true} xs={6} className={props.classes.mdUpLoginForm}>
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
