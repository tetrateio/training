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
    gridContainer: {
        width: "100%",
    },
    mdUpContent: {
        bottom: "0",
        height: "85vh",
        left: "0",
        margin: "auto",
        position: "absolute",
        right: "0",
        top: "15vh",
        width: "75vw",
    },
    mdUpHeader: {
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
    smDownContent: {
        bottom: "0",
        height: "85vh",
        left: "0",
        margin: "auto",
        position: "absolute",
        right: "0",
        top: "15vh",
        width: "100vw",
    },
    smDownHeader: {
        height: "15vh",
        margin: "auto",
        position: "relative",
        top: "-15.5vh",
        width: "100vw",
    },
    smDownMainPanel: {
        width: "100vw",
    },
});

interface IProps extends WithStyles<typeof styles> {
    showRightPanel?: boolean;
    children: JSX.Element;
}

const SmDownComponent: React.FunctionComponent<IProps> = (props: IProps) => {
    return (
        <>
            <div className={props.classes.banner}/>
            <div className={props.classes.smDownHeader}>
                <Header/>
            </div>
            <Hidden mdUp={true}>
                <div className={props.classes.smDownContent}>
                    <Grid
                        container={true}
                        className={props.classes.gridContainer}
                    >
                        <Grid item={true} className={props.classes.smDownMainPanel}>
                            {props.children}
                        </Grid>
                    </Grid>
                </div>
            </Hidden>
        </>
    );
};

const MdUpComponent: React.FunctionComponent<IProps> = (props: IProps) => {
    const rightPanelElem =
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
                <>
                    <Grid item={true} className={props.classes.mdUpMainPanelWithoutRightPanel}>
                        {props.children}
                    </Grid>
                </>
            );
    return (
        <>
            <div className={props.classes.banner}/>
            <div className={props.classes.mdUpHeader}>
                <Header/>
            </div>
            <Hidden smDown={true}>
                <div className={props.classes.mdUpContent}>
                    <Grid
                        container={true}
                        className={props.classes.gridContainer}
                    >
                        {rightPanelElem}
                    </Grid>
                </div>
            </Hidden>
        </>
    );
};

const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    return (
        <>
            <Hidden smDown={true}>
                <MdUpComponent {...props}/>
            </Hidden>
            <Hidden mdUp={true}>
                <SmDownComponent {...props}/>
            </Hidden>
        </>
    );
};

export const Shell = withStyles(styles)(Component);
