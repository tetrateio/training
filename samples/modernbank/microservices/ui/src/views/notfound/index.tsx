import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {MoodBad} from "@material-ui/icons";
import Grid from "@material-ui/core/Grid";
import {Shell} from "../../components/shell";

const styles = (theme: Theme) => createStyles({
    gridContainer: {
        // height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    paper: {
        backgroundColor: "rgba(255,255,255,0.96)",
        display: "flex",
        height: "100%",
        justifyContent: "center",
        paddingTop: "50px",
    },
    text: {},
});

interface IProps extends WithStyles<typeof styles> {
}

const Component: React.FunctionComponent<IProps> = (props: IProps) => (
    <Paper className={props.classes.paper}>
        <Grid
            container={true}
            direction={"column"}
            alignItems={"center"}
            justify={"space-between"}
            className={props.classes.gridContainer}
        >
            <Grid item={true}>
                <Typography variant={"h3"} className={props.classes.text}>Page not found</Typography>
            </Grid>
            <Grid item={true}>
                <MoodBad/>
            </Grid>
        </Grid>
    </Paper>
);

const StyledComponent = withStyles(styles)(Component);

export const NotFoundView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell>
        <StyledComponent/>
    </Shell>
);
