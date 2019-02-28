import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import React from "react";
import Paper from "@material-ui/core/Paper";
import {AccountBalance} from "@material-ui/icons";
import Typography from "@material-ui/core/Typography";

const borderBottomWidth = 3;
const height = 50;

const styles = (theme: Theme) => createStyles({
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    leftSide: {
        display: "flex",
        alignItems: "center",
    },
    paper: {
        borderBottom: `${borderBottomWidth}px solid rgb(196,196,196)`,
        height: `${height}px`,
        padding: "0px 20px",
    },
    spacer: {
        width: 2 * theme.spacing.unit,
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
                    <div className={props.classes.leftSide}>
                        <AccountBalance/>
                        <span className={props.classes.spacer}/>
                        <Typography variant={"subtitle2"} inline={true}>Checking account</Typography>
                    </div>
                </Grid>
                <Grid item={true}>
                    <Typography variant={"subtitle2"}>Available balance</Typography>
                </Grid>
            </Grid>
        </Paper>
    );
};

export const Subheader = withStyles(styles)(component);
