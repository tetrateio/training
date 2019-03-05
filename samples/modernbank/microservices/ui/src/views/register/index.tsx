import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Paper from "@material-ui/core/Paper";
import React from "react";
import {Shell} from "../../components/shell";
import {RegisterForm} from "./registerForm";

const styles = () => createStyles({
    paper: {
        backgroundColor: "rgba(255,255,255,0.97)",
        paddingBottom: "8vh",
        paddingLeft: "20vw",
        paddingRight: "20vw",
        paddingTop: "5vh",
    },
});

interface IProps extends WithStyles<typeof styles> {
}

const Component: React.FunctionComponent<IProps> = (props: IProps) => (
    <Paper square={true} className={props.classes.paper}>
        <RegisterForm/>
    </Paper>
);

const StyledComponent = withStyles(styles)(Component);

export const RegisterView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell>
        <StyledComponent {...props}/>
    </Shell>
);
