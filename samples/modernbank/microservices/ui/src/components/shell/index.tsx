import {install} from "@material-ui/styles";
import {createStyles, Theme, withStyles, WithStyles} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import makeStyles from "@material-ui/styles/makeStyles";
import React from "react";
import "typeface-roboto";
import {Header} from "./header";
import "./index.css";
import {RightPanel} from "./rightPanel";

// CSS pixel constants.
const bannerHeight = 100;
export const bannerBorderBottomWidth = 7;
const headerWidth = 1000;
const contentWidth = headerWidth;

// CSS constants that vary on the component props.
const mainPanelWidth = (props: IProps): number => {
    return (props.showRightPanel) ? 750 : headerWidth;
};
const rightPanelWidth = (props: IProps): number => {
    return (props.showRightPanel) ? (contentWidth - mainPanelWidth(props)) : 0;
};

const useStyles = makeStyles({
    mainPanel: (props: IProps) => ({
        width: `${mainPanelWidth(props)}px`,
    }),
    rightPanel: (props: IProps) => ({
        width: `${rightPanelWidth(props)}px`,
    }),
});

const styles = (theme: Theme) => createStyles({
    banner: {
        backgroundColor: "rgba(130,138,161, 0.99)",
        borderBottom: `${bannerBorderBottomWidth}px solid rgb(172,235,252)`,
        height: `${bannerHeight}px`,
        width: "100vw",
    },
    content: {
        bottom: "0",
        left: "0",
        margin: "auto",
        position: "absolute",
        right: "0",
        top: `${bannerHeight}px`,
        width: `${contentWidth}px`,
    },
    gridContainer: {
        height: "100vh",
        width: "100%", /* Force the grid to be same size as parent div. */
    },
    header: {
        height: "100px",
        margin: "auto",
        position: "relative",
        top: `${-1 * (bannerHeight + bannerBorderBottomWidth)}px`,
        width: `${headerWidth}px`,
    },
});

interface IProps extends WithStyles<typeof styles> {
    showRightPanel?: boolean;
    children: JSX.Element;
}

const component: React.FunctionComponent<IProps> = (props: IProps) => {
    const dynamicClasses = useStyles(props);

    const rightPanel = (props.showRightPanel)
        ? (
            <Grid item={true} className={dynamicClasses.rightPanel}>
                <RightPanel/>
            </Grid>
        )
        : (<></>);

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
                    <Grid item={true} className={dynamicClasses.mainPanel}>
                        {props.children}
                    </Grid>
                    {rightPanel}
                </Grid>
            </div>
        </>
    );
}

export const Shell = withStyles(styles)(component);
