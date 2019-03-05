import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import React from "react";
import {bannerBorderBottomWidth, Shell} from "../../components/shell";
import {RegisterForm} from "./registerForm";
import Paper from "@material-ui/core/Paper";

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
        // opacity: 0.5,
        // width: "100%",
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
    <Paper square={true} className={props.classes.paper}>
        <RegisterForm/>
    </Paper>
);

const StyledComponent = withStyles(styles)(Component);

export const RegisterView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell>
        <StyledComponent/>
    </Shell>
);
