import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import {MoodBad} from "@material-ui/icons";
import React from "react";
import {Shell} from "../../components/shell";

const styles = () => createStyles({
    gridContainer: {
        height: "100%",
    },
    paper: {
        backgroundColor: "rgba(255,255,255,0.95)",
        height: "40vh",
        paddingBottom: "20vh",
        paddingTop: "20vh",
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
            justify={"space-around"}
            className={props.classes.gridContainer}
        >
            <Grid item={true}>
                <Typography variant={"h3"} className={props.classes.text}>Page not found</Typography>
            </Grid>
            <Grid item={true}>
                <MoodBad fontSize="large"/>
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
