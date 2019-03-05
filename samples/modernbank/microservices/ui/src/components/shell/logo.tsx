import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import React from "react";

const styles = (theme: Theme) => createStyles({
    img: {
        filter: "invert(100%)",
        width: "8vw",
    },
});

interface IProps extends WithStyles<typeof styles> {
}

export const component: React.FunctionComponent<IProps> = (props: IProps) => {
    return (
        <img
            className={props.classes.img}
            src={"/logo.png"}
        />
    );
};

export const Logo = withStyles(styles)(component);
