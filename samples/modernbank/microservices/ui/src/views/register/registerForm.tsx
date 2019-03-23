import {
  createStyles,
  WithStyles,
  withStyles,
  Select,
  MenuItem,
  Input,
  Checkbox,
  FormControlLabel
} from '@material-ui/core';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';
import { CreateAccountTypeEnum } from '../../api/client';
import { usersApi, accountsApi } from '../../api/client-utils';
import { AuthContext } from '../../components/auth/authContext';
import { AccountsPageLink, AccountsPath } from '../../routes';
import { VersionContext } from '../../context/versionProvider';
import useSessionstorage from '@rooks/use-sessionstorage';

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
  accountType: boolean;
}

const disableSubmitButton = (s: IFormState): boolean => {
  return (
    !s.username ||
    !s.firstName ||
    !s.lastName ||
    !s.email ||
    !s.password ||
    !s.passwordConfirmation ||
    !s.accountType
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
  const [cash, setCash] = React.useState<boolean>(false);
  const [saving, setSaving] = React.useState<boolean>(false);

  const authContext = React.useContext(AuthContext);
  const { setVersion } = React.useContext(VersionContext);
  const { set } = useSessionstorage('user');

  const createAccount = async (owner: string, type: CreateAccountTypeEnum) => {
    await accountsApi.createAccount({ owner, type });
  };

  const submitNewUserForm = async () => {
    try {
      const resp = await usersApi.createUserRaw({
        body: { email, firstName, lastName, password, username }
      });

      const newUser = await resp.value();
      setVersion(resp.raw.headers.get('version'));

      if (cash) {
        await createAccount(username, CreateAccountTypeEnum.Cash);
      }

      if (saving) {
        await createAccount(username, CreateAccountTypeEnum.Saving);
      }

      authContext.isAuthenticated = true;
      authContext.user = newUser;

      // We set the authenticated user to the session storage.
      set(JSON.stringify(newUser));
      props.history.push(AccountsPath);
    } catch (err) {
      // TODO(dio): show the error.
      console.log(err);
    }
  };

  // Helper method to package all form inputs into a typed object.
  const formState = (): IFormState => ({
    email,
    firstName,
    lastName,
    password,
    passwordConfirmation,
    username,
    accountType: cash || saving
  });

  const handleChange = (name: string) => (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    // TODO(dio): I'm lazy!
    if (name === 'cash') {
      setCash(event.target.checked);
    } else if (name === 'saving') {
      setSaving(event.target.checked);
    }
  };

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

      <FormControlLabel
        control={<Checkbox checked={cash} onChange={handleChange('cash')} />}
        label="Cash"
      />

      <FormControlLabel
        control={
          <Checkbox checked={saving} onChange={handleChange('saving')} />
        }
        label="Saving"
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
