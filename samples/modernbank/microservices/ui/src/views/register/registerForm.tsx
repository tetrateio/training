import { createStyles, WithStyles, withStyles } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';
import { User } from '../../api/client';
import { usersApi } from '../../api/client-utils';
import { AuthContext } from '../../components/auth/authContext';
import { AccountsPageLink, AccountsPath } from '../../routes';

const styles = () =>
  createStyles({
    button: {
      margin: '1vw',
      width: '15vw'
    },
    buttons: {
      paddingTop: '2vh'
    },
    textField: {}
  });

interface IProps extends WithStyles<typeof styles>, RouteComponentProps<{}> {}

interface IFormState {
  username: string;
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  passwordConfirmation: string;
}

const disableSubmitButton = (s: IFormState): boolean => {
  return (
    !s.username ||
    !s.firstName ||
    !s.lastName ||
    !s.email ||
    !s.password ||
    !s.passwordConfirmation
  );
};

export const Component: React.FunctionComponent<IProps> = (props: IProps) => {
  const [username, setUsername] = React.useState<string>('');
  const [firstName, setFirstName] = React.useState<string>('');
  const [lastName, setLastName] = React.useState<string>('');
  const [email, setEmail] = React.useState<string>('');
  const [password, setPassword] = React.useState<string>('');
  const [passwordConfirmation, setPasswordConfirmation] = React.useState<
    string
  >('');

  const authContext = React.useContext(AuthContext);

  const submitNewUserForm = async () => {
    const newUser: User = await usersApi.createUser({
      email,
      firstName,
      lastName,
      password,
      username
    });
    authContext.isAuthenticated = true;
    authContext.user = newUser;
    props.history.push(AccountsPath);
  };

  // Helper method to package all form inputs into a typed object.
  const formState = (): IFormState => ({
    email,
    firstName,
    lastName,
    password,
    passwordConfirmation,
    username
  });

  const handleKeyPress = async (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      await submitNewUserForm();
    }
  };

  return (
    <form onKeyPress={handleKeyPress}>
      <TextField
        value={username}
        onChange={e => setUsername(e.target.value)}
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
        onChange={e => setFirstName(e.target.value)}
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
        onChange={e => setLastName(e.target.value)}
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
        onChange={e => setEmail(e.target.value)}
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
        onChange={e => setPassword(e.target.value)}
        id="password"
        label="Password"
        margin="normal"
        variant="outlined"
        type="password"
        fullWidth={true}
        required={true}
        className={props.classes.textField}
      />
      <TextField
        value={passwordConfirmation}
        onChange={e => setPasswordConfirmation(e.target.value)}
        error={passwordConfirmation !== '' && password !== passwordConfirmation}
        id="password-confirmation"
        label="Confirm Password"
        margin="normal"
        variant="outlined"
        type="password"
        fullWidth={true}
        required={true}
        className={props.classes.textField}
      />
      <div className={props.classes.buttons}>
        <Button
          variant="contained"
          color="primary"
          disabled={disableSubmitButton(formState())}
          onClick={() => submitNewUserForm()}
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
    </form>
  );
};

const RoutingAwareComponent = withRouter(Component);

export const RegisterForm = withStyles(styles)(RoutingAwareComponent);
