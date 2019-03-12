import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import Paper from "@material-ui/core/Paper";
import Snackbar from "@material-ui/core/Snackbar";
import SnackbarContent from "@material-ui/core/SnackbarContent";
import TextField from "@material-ui/core/TextField";
import React from "react";
import {RouteComponentProps, withRouter} from "react-router";
import {authenticationCheck} from "../../api/client-utils";
import {AuthContext} from "../../components/auth/authContext";
import {AccountsPath} from "../../routes";

const styles = () => createStyles({
    paper: {
        backgroundColor: "rgba(255,255,255,0.95)",
        boxSizing: "border-box",
        height: "40vh",
        paddingLeft: "3vw",
        paddingRight: "3vw",
        paddingTop: "4vw",
        width: "100%",
    },
    passwordTextField: {
    },
    signOnButton: {
        backgroundColor: "rgb(233,121,51)",
        marginTop: "3vh",
        textTransform: "none", /* Material button text defaults to all-caps; disable this. */
    },
    usernameTextField: {
    },
});

interface IProps extends WithStyles<typeof styles>, RouteComponentProps<{}> {
}

const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    const [username, setUsername] = React.useState<string>("");
    const [password, setPassword] = React.useState<string>("");
    const [loginFailed, setLoginFailed] = React.useState<boolean>(false);
    const authContext = React.useContext(AuthContext);

    const signInHandler = async () => {
        const authenticatedUser = await authenticationCheck(username, password);
        if (authenticatedUser) {
            setLoginFailed(false);
            authContext.isAuthenticated = true;
            authContext.user = authenticatedUser;
            props.history.push(AccountsPath);
        } else {
            setLoginFailed(true);
        }
    };

    const handleKeyPress = async (e: React.KeyboardEvent) => {
        if (e.key === "Enter") {
            await signInHandler();
        }
    };

    return (
        <Paper square={true} className={props.classes.paper}>
            <form onKeyPress={handleKeyPress}>
                <TextField
                    value={username}
                    onChange={(event) => setUsername(event.target.value)}
                    id="username-input"
                    label="Username"
                    margin="normal"
                    variant="outlined"
                    fullWidth={true}
                    required={true}
                    className={props.classes.usernameTextField}
                />
                <TextField
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    id="password-input"
                    label="Password"
                    margin="normal"
                    variant="outlined"
                    fullWidth={true}
                    required={true}
                    type="password"
                    className={props.classes.passwordTextField}
                />
                <Button
                    variant={"outlined"}
                    fullWidth={true}
                    size={"large"}
                    onClick={signInHandler}
                    className={props.classes.signOnButton}
                >
                    Sign on
                </Button>
            </form>
            <Snackbar
                open={loginFailed}
                autoHideDuration={3000}
                onClose={() => setLoginFailed(false)}
            >
                <SnackbarContent
                    message="Login failed"
                />
            </Snackbar>
        </Paper>
    );
};

const RoutingAwareComponent = withRouter(Component);

export const LoginForm = withStyles(styles)(RoutingAwareComponent);
