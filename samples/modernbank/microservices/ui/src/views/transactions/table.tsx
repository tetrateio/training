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
import { Newtransaction, TransactionsApi } from '../../api/client';

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

  const [receivedTransactions, setReceivedTxs] = React.useState<
    Newtransaction[]
  >([]);
  const [sentTransactions, setSentTxs] = React.useState<Newtransaction[]>([]);

  const [doFetch, setDoFetch] = React.useState<boolean>(true);

  const transactionsApi = new TransactionsApi();

  const fetchTransactionsReceived = async () => {
    const received: any = await transactionsApi.listTransactionsReceived({
      receiver: accountNumber
    });
    setReceivedTxs(received);
  };

  const fetchTransactionsSent = async () => {
    const sent: any = await transactionsApi.listTransactionsSent({
      sender: accountNumber
    });
    setSentTxs(sent);
  };

  React.useEffect(() => {
    if (doFetch) {
      setDoFetch(false);
      fetchTransactionsReceived();
      fetchTransactionsSent();
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
          {receivedTransactions.map(transaction => (
            <TableRow>
              <TableCell component="th" scope="row" align="left">
                <Typography>{transaction.sender}</Typography>
              </TableCell>
              <TableCell align="right">{transaction.amount}</TableCell>
              <TableCell />
              <TableCell />
              <TableCell />
            </TableRow>
          ))}
          {sentTransactions.map(transaction => (
            <TableRow>
              <TableCell />
              <TableCell />
              <TableCell component="th" scope="row" align="left">
                <Typography>{transaction.receiver}</Typography>
              </TableCell>
              <TableCell align="right">{transaction.amount}</TableCell>
              <TableCell />
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Paper>
  );
};

const RoutingComponent = withRouter(Component);

export const TransactionsTable = withStyles(styles)(RoutingComponent);
