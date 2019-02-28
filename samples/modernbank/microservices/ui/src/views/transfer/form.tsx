import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import Divider from "@material-ui/core/Divider";
import FormControl from "@material-ui/core/FormControl";
import InputLabel from "@material-ui/core/InputLabel";
import MenuItem from "@material-ui/core/MenuItem";
import Select from "@material-ui/core/Select";
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
        marginLeft: theme.spacing.unit,
        marginRight: theme.spacing.unit,
    },
});

interface IProps extends WithStyles<typeof styles> {
}

const fetchedAccounts: number[] = [
    1001, 1002, 1003,
];

interface IFormState {
    fromAccount: string;
    toAccount: string;
    routingNumber: string;
    date: string;
    amount: string;
}

function isPositiveInteger(s: string): boolean {
    const n = Number(s);
    return !isNaN(n) && Number.isInteger(n) && n > 0;
}

function disableSubmitButton(s: IFormState): boolean {
    return !s.fromAccount || !s.toAccount || !s.routingNumber || !s.date || !s.amount;
}

export const component: React.FunctionComponent<IProps> = (props: IProps) => {
    const [fromAccount, setFromAccount] = React.useState<string>("");
    const [toAccount, setToAccount] = React.useState<string>("");
    const [routingNumber, setRoutingNumber] = React.useState<string>("");
    const [date, setDate] = React.useState<string>("");
    const [amount, setAmount] = React.useState<string>("");

    // Helper method to package all form inputs into a typed object.
    function formState(): IFormState {
        return {
            amount,
            date,
            fromAccount,
            routingNumber,
            toAccount,
        };
    }

    // onChange event handlers for every form input.
    function intFieldInputHandler(setStateFunc: Dispatch<SetStateAction<string>>) {
        return (event: React.ChangeEvent<HTMLSelectElement>) => {
            setStateFunc(event.target.value);
        };
    }

    const fromAccountInputHandler = intFieldInputHandler(setFromAccount);
    const toAccountInputHandler = intFieldInputHandler(setToAccount);
    const routingNumberInputHandler = intFieldInputHandler(setRoutingNumber);
    const dateInputHandler = intFieldInputHandler(setDate);
    const amountInputHandler = intFieldInputHandler(setAmount);

    // Input validation functions.
    function isValidPositiveIntegerInput(input: string): boolean {
        if (input === "") {
            // Don't invalidate for empty inputs.
            return true;
        }
        return isPositiveInteger(input);
    }

    return (
        <form>
            <FormControl className={props.classes.formControl}>
                <InputLabel>From account</InputLabel>
                <Select
                    value={fromAccount}
                    onChange={fromAccountInputHandler}
                >
                    {fetchedAccounts.map((accountNumber: number) => (
                        <MenuItem value={accountNumber} key={accountNumber}>
                            {accountNumber}
                        </MenuItem>
                    ))}
                </Select>
            </FormControl>
            <TextField
                value={toAccount}
                onChange={toAccountInputHandler}
                error={!isValidPositiveIntegerInput(toAccount)}
                id="to-account-input"
                label="To account"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <TextField
                value={routingNumber}
                onChange={routingNumberInputHandler}
                error={!isValidPositiveIntegerInput(routingNumber)}
                id="routing-number-input"
                label="Routing number"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <TextField
                value={date}
                onChange={dateInputHandler}
                error={!isValidPositiveIntegerInput(date)}
                id="date-input"
                label="Date"
                margin="normal"
                variant="outlined"
                required={true}
                className={props.classes.textField}
            />
            <TextField
                value={amount}
                onChange={amountInputHandler}
                error={!isValidPositiveIntegerInput(amount)}
                id="amount-input"
                label="Amount"
                margin="normal"
                variant="outlined"
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
                <p>fromAccount = {fromAccount}</p>
                <p>toAccount = {toAccount}</p>
                <p>routingNumber = {routingNumber}</p>
                <p>date = {date}</p>
                <p>amount = {amount}</p>
            </div>
        </form>
    );
};

export const Form = withStyles(styles)(component);
