import { createStyles, WithStyles, withStyles } from '@material-ui/core';
import Divider from '@material-ui/core/Divider';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import React from 'react';

import { Account, AccountsApi } from '../../api/client';
import { AuthContext } from '../../components/auth/authContext';
import { Shell } from '../../components/shell';
import { AccountSummary } from './accountSummary';
import { CashAccounts } from './cashAccounts';
import { CreditAccounts } from './creditAccounts';
import { InvestmentAccounts } from './investmentAccounts';
import { TotalAccountValue } from './totalAccountValue';
import { VersionContext } from '../../context/versionProvider';

const styles = () =>
  createStyles({
    divider: {
      backgroundColor: 'black',
      height: '0.25vh'
    },
    fillerPaper: {
      backgroundColor: 'rgba(255,255,255,0.95)',
      height: '100%'
    },
    fillerPaperGridItem: {
      height: '100%'
    },
    gridContainer: {
      flexWrap: 'nowrap',
      height: '100%'
    },
    subheader: {
      backgroundColor: 'rgba(172,37,45, 1)'
    },
    subheaderText: {
      color: 'white',
      marginLeft: '30px'
    }
  });

interface IProps extends WithStyles<typeof styles> {}

const accountsApi = new AccountsApi();

const Component: React.FunctionComponent<IProps> = (props: IProps) => {
  const [accounts, setAccounts] = React.useState<Account[]>([]);
  const authContext = React.useContext(AuthContext);
  const { setVersion } = React.useContext(VersionContext);

  const fetchAccounts = async () => {
    const owner = authContext.user!.username;
    const resp = await accountsApi.listAccountsRaw({ owner });

    setVersion(resp.raw.headers.get('version'));
    setAccounts(await resp.value());
  };

  React.useEffect(() => {
    fetchAccounts();
  }, []);

  // The API doesn't support account type. Fake account type using the last digit of the account number.
  const filterByLastDigit = (start: number, end: number): Account[] => {
    return accounts.filter(
      acc => start <= acc.number % 10 && acc.number % 10 <= end
    );
  };

  const cashAccounts = filterByLastDigit(0, 3);
  const investmentAccounts = filterByLastDigit(4, 6);
  const creditAccounts = filterByLastDigit(7, 9);

  return (
    <Grid
      container={true}
      alignItems={'stretch'}
      direction={'column'}
      justify={'flex-start'}
      className={props.classes.gridContainer}
    >
      <Grid item={true}>
        <div className={props.classes.subheader}>
          <Typography variant="h6" className={props.classes.subheaderText}>
            Account summary / Checking account
          </Typography>
        </div>
      </Grid>

      <Grid item={true}>
        <AccountSummary
          plusAccounts={cashAccounts.concat(investmentAccounts)}
          minusAccounts={creditAccounts}
        />
      </Grid>
      <Grid item={true}>
        <CashAccounts accounts={cashAccounts} />
      </Grid>
      <Grid item={true}>
        <InvestmentAccounts accounts={investmentAccounts} />
      </Grid>
      <Grid item={true}>
        <CreditAccounts accounts={creditAccounts} />
      </Grid>
      <Grid item={true}>
        <Divider className={props.classes.divider} />
      </Grid>
      <Grid item={true}>
        <TotalAccountValue
          plusAccounts={cashAccounts.concat(investmentAccounts)}
          minusAccounts={creditAccounts}
        />
      </Grid>
      <Grid item={true} className={props.classes.fillerPaperGridItem}>
        <Paper square={true} className={props.classes.fillerPaper} />
      </Grid>
    </Grid>
  );
};

const StyledComponent = withStyles(styles)(Component);

export const AccountsView: React.FunctionComponent<IProps> = (
  props: IProps
) => (
  <Shell showRightPanel={true}>
    <StyledComponent {...props} />
  </Shell>
);
