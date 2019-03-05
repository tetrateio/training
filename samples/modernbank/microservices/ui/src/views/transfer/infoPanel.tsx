import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Typography from "@material-ui/core/Typography";
import React from "react";

const styles = () => createStyles({
    headerText: {
        paddingBottom: "3vh",
    },
    root: {
        paddingTop: "4vh",
    },
    subheaderText: {
    },
});

interface IProps extends WithStyles<typeof styles> {
}

export const Component: React.FunctionComponent<IProps> = (props: IProps) => (
    <div className={props.classes.root}>
        <Typography variant="h6" className={props.classes.headerText}>
            Send money fast!
        </Typography>
        <Typography variant="body1" className={props.classes.subheaderText}>
            When the money is just pretend, the transfer can be immediate and risk free. It is easy to
            play and fun to try. Enter some amounts and watch the expansion of the money supply.
        </Typography>
    </div>
);

export const InfoPanel = withStyles(styles)(Component);
