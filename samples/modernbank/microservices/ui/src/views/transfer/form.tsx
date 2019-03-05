import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import Divider from "@material-ui/core/Divider";
import TextField from "@material-ui/core/TextField";
import React, {Dispatch, SetStateAction} from "react";
import {RouteComponentProps, withRouter} from "react-router";
import {Account, AccountsApi, TransactionsApi} from "../../api/client";
import {AccountsPageLink, AccountsPath} from "../../routes";
import MenuItem from "@material-ui/core/MenuItem";
import Select from "@material-ui/core/Select";
import InputLabel from "@material-ui/core/InputLabel";
import FormControl from "@material-ui/core/FormControl";
import {AuthContext} from "../../components/auth/authContext";

const styles = () => createStyles({
    button: {
        margin: "1vh",
        width: "20vh",
    },
    formControl: {
        paddingBottom: "1vh",
        width: "100%",
    },
    textField: {
        paddingBottom: "1vh",
    },
});

interface IUrlParams {
    accountNumber: string;
}

interface IProps extends WithStyles<typeof styles>, RouteComponentProps<IUrlParams> {
}

interface IFormState {
    fromAccount: string;
    toAccount: string;
    amount: string;
}

function isPositiveInteger(s: string): boolean {
    const n = Number(s);
    return !isNaN(n) && Number.isInteger(n) && n > 0;
}

// Input validation functions.
function isValidPositiveIntegerInput(input: string): boolean {
    if (input === "") {
        // Don't invalidate empty inputs.
        return true;
    }
    return isPositiveInteger(input);
}

function disableSubmitButton(s: IFormState): boolean {
    return !s.fromAccount || !s.toAccount || !s.amount ||
        !isValidPositiveIntegerInput(s.fromAccount) ||
        !isValidPositiveIntegerInput(s.toAccount) ||
        !isValidPositiveIntegerInput(s.amount);
}

const accountsApi = new AccountsApi({basePath: "http://35.192.59.252/v1"});

export const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    const [toAccount, setToAccount] = React.useState<string>("");
    const [amount, setAmount] = React.useState<string>("");
    const authContext = React.useContext(AuthContext);

    const fromAccount = props.match.params.accountNumber;

    const [accounts, setAccounts] = React.useState<Account[]>([]);
    const fetchAccounts = async () => {
        const resp: Account[] = await accountsApi.listAccounts(authContext.user!.username);
        setAccounts(resp);
    };

    React.useEffect(() => {
        fetchAccounts();
    }, []);

    // Helper method to package all form inputs into a typed object.
    function formState(): IFormState {
        return {
            amount,
            fromAccount,
            toAccount,
        };
    }

    // onChange event handlers for every form input.
    function intFieldInputHandler(setStateFunc: Dispatch<SetStateAction<string>>) {
        return (event: React.ChangeEvent<HTMLSelectElement>) => {
            setStateFunc(event.target.value);
        };
    }

    const toAccountInputHandler = intFieldInputHandler(setToAccount);
    const amountInputHandler = intFieldInputHandler(setAmount);

    const submitTransfer = async () => {
        const transactionsApi = new TransactionsApi({basePath: "http://35.192.59.252/v1"});
        const newTransaction = await transactionsApi.createTransaction({
            amount: Number(amount),
            receiver: parseInt(toAccount, 10),
            sender: parseInt(fromAccount, 10),
        });
        // TODO: block on await
        props.history.push(AccountsPath);
    };

    return (
        <form>
            <TextField
                value={fromAccount}
                id="from-account-read-only"
                label="From account"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
                InputProps={{
                    readOnly: true,
                }}
            />
            <FormControl className={props.classes.formControl}>
                <InputLabel>To account</InputLabel>
                <Select
                    value={toAccount}
                    onChange={toAccountInputHandler}
                >
                    {accounts
                        .filter((account) => account.number !== parseInt(props.match.params.accountNumber, 10))
                        .map((account: Account) => (
                            <MenuItem value={account.number} key={account.number}>
                                {account.number}
                            </MenuItem>
                        ))}
                </Select>
            </FormControl>
            <TextField
                value={amount}
                onChange={amountInputHandler}
                error={!isValidPositiveIntegerInput(amount)}
                id="amount-input"
                label="Amount"
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
                    onClick={() => submitTransfer()}
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

const RoutingComponent = withRouter(Component);

export const Form = withStyles(styles)(RoutingComponent);
