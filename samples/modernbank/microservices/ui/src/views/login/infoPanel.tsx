import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import React from "react";
import Grid from "@material-ui/core/Grid";
import {RegisterPageLink} from "../../routes";

const styles = (theme: Theme) => createStyles({
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
        paddingRight: 3 * theme.spacing.unit,
        width: "100%",
    },
    headerText: {
        color: "white",
        marginBottom: 5 * theme.spacing.unit,
        marginTop: 10 * theme.spacing.unit,
        textAlign: "right",
    },
    joinNowButton: {
        backgroundColor: "rgb(172,235,252)",
    },
    root: {
        display: "flex",
    },
    subheaderText: {
        color: "white",
        marginBottom: 3 * theme.spacing.unit,
        textAlign: "right",
    },
    subtitleTextGridItem: {
        width: "50%",
    },
});

interface IProps extends WithStyles<typeof styles> {
}

const Component: React.FunctionComponent<IProps> = (props: IProps) => (
    <Grid
        container={true}
        alignItems={"flex-end"}
        direction={"column"}
        className={props.classes.gridContainer}
    >
        <Grid item={true}>
            <Typography variant="h3" className={props.classes.headerText}>
                Connecting dreams with cash
            </Typography>
        </Grid>
        <Grid item={true} className={props.classes.subtitleTextGridItem}>
            <Typography variant="body1" className={props.classes.subheaderText}>
                Open a checking account and get fake cash for free with participating online demo accounts
            </Typography>
        </Grid>
        <Grid item={true}>
            <Button
                component={RegisterPageLink}
                variant={"outlined"} className={props.classes.joinNowButton}>
                Join now
            </Button>
        </Grid>
    </Grid>
);

export const InfoPanel = withStyles(styles)(Component);
