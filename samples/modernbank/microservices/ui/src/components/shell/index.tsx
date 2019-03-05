import {createStyles, Theme, withStyles, WithStyles} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Hidden from "@material-ui/core/Hidden";
import React from "react";
import "typeface-roboto";
import {Header} from "./header";
import "./index.css";
import {RightPanel} from "./rightPanel";

const styles = (theme: Theme) => createStyles({
    banner: {
        backgroundColor: "rgba(130,138,161, 0.99)",
        borderBottomColor: "rgb(172,235,252)",
        borderBottomStyle: "solid",
        borderBottomWidth: "0.5vh",
        height: "15vh",
        width: "100vw",
    },
    content: {
        bottom: "0",
        left: "0",
        margin: "auto",
        position: "absolute",
        right: "0",
        top: "15vh",
        width: "75vw",
    },
    gridContainer: {
        width: "100%",
    },
    header: {
        height: "15vh",
        margin: "auto",
        position: "relative",
        top: "-15.5vh",
        width: "75vw",
    },
    mdUpMainPanelWithRightPanel: {
        width: "55vw",
    },
    mdUpMainPanelWithoutRightPanel: {
        width: "75vw",
    },
    mdUpRightPanel: {
        width: "20vw",
    },
    smDownMainPanel: {
        width: "75vw",
    },
});

interface IProps extends WithStyles<typeof styles> {
    showRightPanel?: boolean;
    children: JSX.Element;
}

const component: React.FunctionComponent<IProps> = (props: IProps) => {
    const mdUpElement =
        (props.showRightPanel)
            ? (
                <>
                    <Grid item={true} className={props.classes.mdUpMainPanelWithRightPanel}>
                        {props.children}
                    </Grid>
                    <Grid item={true} className={props.classes.mdUpRightPanel}>
                        <RightPanel/>
                    </Grid>
                </>
            ) :
            (
                <Grid item={true} className={props.classes.mdUpMainPanelWithoutRightPanel}>
                    {props.children}
                </Grid>
            );

    return (
        <>
            <div className={props.classes.banner}/>
            <div className={props.classes.header}>
                <Header/>
            </div>
            <div className={props.classes.content}>
                <Grid
                    container={true}
                    className={props.classes.gridContainer}
                >
                    <Hidden smDown={true}>
                        {mdUpElement}
                    </Hidden>
                    <Hidden mdUp={true}>
                        <Grid item={true} className={props.classes.smDownMainPanel}>
                            {props.children}
                        </Grid>
                    </Hidden>
                </Grid>
            </div>
        </>
    );
};

export const Shell = withStyles(styles)(component);
