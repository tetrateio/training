import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import React from "react";
import {bannerBorderBottomWidth, Shell} from "../../components/shell";
import {LoginForm} from "./loginForm";
import {InfoPanel} from "./infoPanel";

const borderTopWidth = bannerBorderBottomWidth;

const styles = (theme: Theme) => createStyles({
    gridContainer: {
        borderTop: `${borderTopWidth}px solid rgb(172,37,45)`,
        height: "100%", /* Force the grid to be same size as parent Paper component. */
        paddingLeft: "50px",
        width: "100%",
    },
    paper: {
        backgroundColor: "white",
        height: "50%",
        // opacity: 0.5,
        width: "100%",
        paddingLeft: 3 * theme.spacing.unit,
        paddingRight: 3 * theme.spacing.unit,
        paddingTop: 5 * theme.spacing.unit,
    },
    signOnButton: {
        color: "rgb(233,121,51)",
        marginTop: 3 * theme.spacing.unit,
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
        <Grid item={true} xs={4}>
            <LoginForm/>
        </Grid>
        <Grid item={true} xs={6}>
            <InfoPanel/>
        </Grid>
    </Grid>
);

const StyledComponent = withStyles(styles)(Component);

export const LoginView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell>
        <StyledComponent/>
    </Shell>
);
