import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {Navbar} from "../../components/viewAppBar/navbar";
import {Subheader} from "../../components/viewAppBar/subheader";
import {DebugForm} from "./debugForm";

const styles = (theme: Theme) => createStyles({
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    paper: {
        paddingLeft: "40px",
    },
    subheader: {
        backgroundColor: "rgba(172,37,45, 1)",
    },
    subheaderText: {
        color: "white",
        marginLeft: "30px",
    },
});

interface IProps extends WithStyles<typeof styles> {
}

export const component: React.FunctionComponent<IProps> = (props: IProps) => (
    <>
        <DebugForm/>
    </>
);

export const DebugView = withStyles(styles)(component);
