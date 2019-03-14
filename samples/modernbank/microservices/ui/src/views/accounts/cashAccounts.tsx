import { createStyles, WithStyles, withStyles } from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import React from 'react';
import { Account } from '../../api/client';
import { AccountCard } from '../../components/accounts/accountCard';

const styles = () =>
  createStyles({
    contentPaper: {
      backgroundColor: 'rgba(255,255,255,0.95)',
      paddingBottom: '1vh',
      paddingTop: '1vh'
    },
    headerPaper: {
      backgroundColor: 'rgb(173,187,202)',
      paddingLeft: '5vh'
    },
    headerText: {
      color: 'white'
    }
  });

interface IProps extends WithStyles<typeof styles> {
  accounts: Account[];
}

export const component: React.FunctionComponent<IProps> = (props: IProps) => (
  <>
    <Paper square={true} className={props.classes.headerPaper}>
      <Typography variant="h6" className={props.classes.headerText}>
        Cash accounts
      </Typography>
    </Paper>
    <Paper square={true} className={props.classes.contentPaper}>
      {props.accounts.map((account: Account) => (
        <AccountCard
          accountName="Checking account"
          accountNumber={account.number}
          accountBalance={account.balance}
          key={account.number}
        />
      ))}
    </Paper>
  </>
);

export const CashAccounts = withStyles(styles)(component);
