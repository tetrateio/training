import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Divider from "@material-ui/core/Divider";
import Typography from "@material-ui/core/Typography";
import React from "react";
import {Shell} from "../../components/shell";
import {AccountSummary} from "./accountSummary";
import {CashAccounts} from "./cashAccounts";
import {CreditAccounts} from "./creditAccounts";
import {InvestmentAccounts} from "./investmentAccounts";
import {TotalAccountValue} from "./totalAccountValue";

const styles = (theme: Theme) => createStyles({
    paper: {
        opacity: 0.5,
    },
    subheader: {
        backgroundColor: "rgba(172,37,45, 1)",
    },
    subheaderText: {
        color: "white",
        marginLeft: "30px",
    },
});

interface IProps extends WithStyles<typeof styles> {
}

const Component: React.FunctionComponent<IProps> = (props: IProps) => (
    <>
        <div className={props.classes.subheader}>
            <Typography variant="h6" className={props.classes.subheaderText}>
                Account summary / Checking account
            </Typography>
        </div>
        <AccountSummary/>
        <CashAccounts/>
        <InvestmentAccounts/>
        <CreditAccounts/>
        <Divider/>
        <TotalAccountValue/>
    </>
);

const StyledComponent = withStyles(styles)(Component);

export const AccountsView: React.FunctionComponent<IProps> = (props: IProps) => (
    <Shell showRightPanel={true}>
        <StyledComponent/>
    </Shell>
);
