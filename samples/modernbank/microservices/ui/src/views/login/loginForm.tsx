import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import Paper from "@material-ui/core/Paper";
import Snackbar from "@material-ui/core/Snackbar";
import SnackbarContent from "@material-ui/core/SnackbarContent";
import TextField from "@material-ui/core/TextField";
import React, {useContext} from "react";
import {Redirect, RouteComponentProps, withRouter} from "react-router";
import {fakeSignInCheck} from "../../components/auth/fakeAuth";
import {AccountsPath, TransferPath} from "../../routes";
import {AuthContext} from "../../components/auth/authContext";

const styles = (theme: Theme) => createStyles({
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
        paddingLeft: "50px",
        width: "100%",
    },
    headerText: {},
    paper: {
        backgroundColor: "rgba(255,255,255,0.95)",
        // filter: "invert(100%)",
        height: "50vh",
        width: "100%",
        paddingLeft: 3 * theme.spacing.unit,
        paddingRight: 3 * theme.spacing.unit,
        paddingTop: 5 * theme.spacing.unit,
    },
    signOnButton: {
        backgroundColor: "rgb(233,121,51)",
        marginTop: 3 * theme.spacing.unit,
    },
    subheader: {
        backgroundColor: "rgba(172,37,45, 1)",
    },
    subheaderText: {
        color: "white",
        marginLeft: "30px",
    },
    textField: {
        // marginLeft: theme.spacing.unit,
        // marginRight: theme.spacing.unit,
        borderBottom: "10px",
    },
    inputField: {},
});

interface IProps extends WithStyles<typeof styles>, RouteComponentProps<{}> {
}

const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    const [username, setUsername] = React.useState<string>("");
    const [password, setPassword] = React.useState<string>("");

    const [loginFailed, setLoginFailed] = React.useState<boolean>(false);

    const authContext = useContext(AuthContext);

    const signInHandler = () => {
        if (fakeSignInCheck(username, password)) {
            setLoginFailed(false);
            authContext.isAuthenticated = true;
            props.history.push(AccountsPath);
        } else {
            setLoginFailed(true);
        }
    };

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === "Enter") {
            signInHandler();
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
                    InputProps={{
                        classes: {
                            root: props.classes.inputField,
                        },
                    }}
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
                    className={props.classes.textField}
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

const RoutingComponent = withRouter(Component);

export const LoginForm = withStyles(styles)(RoutingComponent);
