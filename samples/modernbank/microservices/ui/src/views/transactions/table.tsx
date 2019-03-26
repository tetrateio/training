import { createStyles, WithStyles, withStyles } from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';
import React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';
import {
  Account,
  AccountsApi,
  Newtransaction,
  TransactionsApi
} from '../../api/client';
import useSessionstorage from '@rooks/use-sessionstorage';

const styles = () =>
  createStyles({
    paper: {
      backgroundColor: 'rgba(255,255,255,0.97)'
    },
    tableHead: {
      backgroundColor: 'rgb(173,187,202)'
    },
    tableHeadCell: {
      color: 'white'
    }
  });

interface IUrlParams {
  accountNumber: string;
}

interface IProps
  extends WithStyles<typeof styles>,
    RouteComponentProps<IUrlParams> {}

export const Component: React.FunctionComponent<IProps> = (props: IProps) => {
  const accountNumber = parseInt(props.match.params.accountNumber, 10);

  const [transactions, setTrxs] = React.useState<Newtransaction[]>([]);

  const [doFetch, setDoFetch] = React.useState<boolean>(true);

  const transactionsApi = new TransactionsApi();
  const accountsApi = new AccountsApi();
  const { value } = useSessionstorage('user', '');

  function merge(a: Newtransaction[], b: Newtransaction[]) {
    return a.concat(b).sort((i: Newtransaction, j: Newtransaction) => {
      return i.timestamp - j.timestamp;
    });
  }

  const fetchTransactions = async () => {
    const accounts: any = await accountsApi.listAccounts({
      username: JSON.parse(value).username
    });

    const received: any = await transactionsApi.listTransactionsReceived({
      receiver: accountNumber
    });

    const sent: any = await transactionsApi.listTransactionsSent({
      sender: accountNumber
    });

    const currentAccount: any = accounts
      .filter(account => {
        return account.number === accountNumber;
      })
      .pop();

    const merged = merge(received, sent).map((trx, index) => {
      if (trx.receiver === accountNumber) {
        currentAccount.balance -= trx.amount;
      } else {
        currentAccount.balance += trx.amount;
      }
      trx.balance = currentAccount.balance;
      return trx;
    });

    setTrxs(merged);
  };

  React.useEffect(() => {
    if (doFetch) {
      setDoFetch(false);
      fetchTransactions();
    }
  });

  return (
    <Paper square={true} className={props.classes.paper}>
      <Table>
        <TableHead className={props.classes.tableHead}>
          <TableRow>
            <TableCell className={props.classes.tableHeadCell}>From</TableCell>
            <TableCell className={props.classes.tableHeadCell}>
              Deposits
            </TableCell>
            <TableCell className={props.classes.tableHeadCell}>To</TableCell>
            <TableCell className={props.classes.tableHeadCell}>
              Withdrawals
            </TableCell>
            <TableCell className={props.classes.tableHeadCell}>
              Balance
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {transactions.map((transaction, index) => (
            <TableRow key={index}>
              <TableCell component="th" scope="row" align="left">
                <Typography>{transaction.sender}</Typography>
              </TableCell>
              <TableCell align="right">
                {transaction.sender !== accountNumber ? transaction.amount : ''}
              </TableCell>
              <TableCell component="th" scope="row" align="left">
                <Typography>{transaction.receiver}</Typography>
              </TableCell>
              <TableCell align="right">
                {transaction.receiver !== accountNumber
                  ? transaction.amount
                  : ''}
              </TableCell>
              <TableCell align="right">{transaction.balance}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Paper>
  );
};

const RoutingComponent = withRouter(Component);

export const TransactionsTable = withStyles(styles)(RoutingComponent);
