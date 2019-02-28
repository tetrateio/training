import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {AccountCard} from "../../components/accounts/accountCard";

const styles = (theme: Theme) => createStyles({
    contentPaper: {
        paddingBottom: "10px",
        paddingTop: "10px",
    },
    headerPaper: {
        backgroundColor: "rgb(173,187,202)",
        paddingLeft: "40px",
    },
    headerText: {
        color: "white",
    },
});

interface IProps extends WithStyles<typeof styles> {
}

export const component: React.FunctionComponent<IProps> = (props: IProps) => (
    <>
        <Paper square={true} className={props.classes.headerPaper}>
            <Typography variant="h6" className={props.classes.headerText}>
                Credit
            </Typography>
        </Paper>
        <Paper square={true} className={props.classes.contentPaper}>
            <AccountCard
                accountName="Line of credit"
                accountNumber={8234}
                accountBalance={4393.55}
            />
        </Paper>
    </>
);

export const CreditAccounts = withStyles(styles)(component);
