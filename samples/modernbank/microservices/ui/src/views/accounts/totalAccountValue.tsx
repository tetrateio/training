import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Paper from "@material-ui/core/Paper";
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
    <Paper square={true} className={props.classes.contentPaper}>
        <AccountCard
            accountName="Total account value"
            accountNumber={0}
            accountBalance={413453.55}
        />
    </Paper>
);

export const TotalAccountValue = withStyles(styles)(component);
