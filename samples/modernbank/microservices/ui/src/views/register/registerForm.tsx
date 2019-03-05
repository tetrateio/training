import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import Divider from "@material-ui/core/Divider";
import TextField from "@material-ui/core/TextField";
import React, {Dispatch, SetStateAction} from "react";
import {AccountsPageLink} from "../../routes";

const styles = (theme: Theme) => createStyles({
    button: {
        margin: theme.spacing.unit,
        width: "120px",
    },
    formControl: {
        margin: theme.spacing.unit,
        minWidth: 120,
    },
    selectEmpty: {
        marginTop: 2 * theme.spacing.unit,
    },
    textField: {
    },
});

interface IProps extends WithStyles<typeof styles> {
}

const fetchedAccounts: number[] = [
    1001, 1002, 1003,
];

interface IFormState {
    username: string;
    firstName: string;
    lastName: string;
    email: string;
    password: string;
    passwordConfirmation: string;
}

function disableSubmitButton(s: IFormState): boolean {
    return !s.username || !s.firstName || !s.lastName || !s.email || !s.password || !s.passwordConfirmation;
}

export const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    const [username, setUsername] = React.useState<string>("");
    const [firstName, setFirstName] = React.useState<string>("");
    const [lastName, setLastName] = React.useState<string>("");
    const [email, setEmail] = React.useState<string>("");
    const [password, setPassword] = React.useState<string>("");
    const [passwordConfirmation, setPasswordConfirmation] = React.useState<string>("");

    // Helper method to package all form inputs into a typed object.
    function formState(): IFormState {
        return {
            username,
            firstName,
            lastName,
            email,
            password,
            passwordConfirmation,
        };
    }

    return (
        <form>
            <TextField
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                id="username"
                label="Username"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <TextField
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                id="first-name"
                label="First Name"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <TextField
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                id="last-name"
                label="Last Name"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <TextField
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                id="email"
                label="Email"
                margin="normal"
                variant="outlined"
                type="email"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <TextField
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                id="password"
                label="Password"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <TextField
                value={passwordConfirmation}
                onChange={(e) => setPasswordConfirmation(e.target.value)}
                error={passwordConfirmation !== "" && password !== passwordConfirmation}
                id="password-confirmation"
                label="Confirm Password"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <Divider/>
            <div>
                <Button
                    variant="contained"
                    color="primary"
                    disabled={disableSubmitButton(formState())}
                    onClick={() => {console.log("Submit")}}
                    className={props.classes.button}
                >
                    Submit
                </Button>
                <Button
                    variant="contained"
                    component={AccountsPageLink}
                    className={props.classes.button}
                >
                    Cancel
                </Button>
            </div>
            <div>
                <p>fromAccount = {username}</p>
                <p>toAccount = {firstName}</p>
                <p>routingNumber = {lastName}</p>
                <p>date = {email}</p>
                <p>amount = {password}</p>
            </div>
        </form>
    );
};

export const RegisterForm = withStyles(styles)(Component);
