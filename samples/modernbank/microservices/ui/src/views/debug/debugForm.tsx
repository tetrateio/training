import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import Divider from "@material-ui/core/Divider";
import Grid from "@material-ui/core/Grid";
import MenuItem from "@material-ui/core/MenuItem";
import Paper from "@material-ui/core/Paper";
import Select from "@material-ui/core/Select";
import TextField from "@material-ui/core/TextField";
import Typography from "@material-ui/core/Typography";
import React from "react";
import FormControl from "@material-ui/core/FormControl";
import InputLabel from "@material-ui/core/InputLabel";
import Input from "@material-ui/core/Input";

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
        marginTop: theme.spacing.unit * 2,
    },
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    paper: {
        paddingLeft: "40px",
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

export const component: React.FunctionComponent<IProps> = (props: IProps) => {
    const [fromAccount, setFromAccount] = React.useState<number>(0);
    const [toAccount, setToAccount] = React.useState<number>(0);
    const [routingNumber, setRoutingNumber] = React.useState<number>(0);
    const [date, setDate] = React.useState<number>(0);
    const [amount, setAmount] = React.useState<number>(0);

    return (
        <div>
            <TextField
                id="to-account-input"
                label="To account"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <TextField
                id="routing-number-input"
                label="Routing number"
                margin="normal"
                variant="outlined"
                fullWidth={true}
                required={true}
                className={props.classes.textField}
            />
            <TextField
                id="date-input"
                label="Date"
                margin="normal"
                variant="outlined"
                required={true}
                className={props.classes.textField}
            />
            <TextField
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
                    className={props.classes.button}
                >
                    Submit
                </Button>
                <Button
                    variant="contained"
                    className={props.classes.button}
                >
                    Cancel
                </Button>
            </div>
        </div>
    );
}

export const DebugForm = withStyles(styles)(component);
