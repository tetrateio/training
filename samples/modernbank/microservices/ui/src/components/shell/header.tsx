import { createStyles, WithStyles, withStyles } from '@material-ui/core';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import { AccountCircle } from '@material-ui/icons';
import React from 'react';
import { RouteComponentProps, withRouter } from 'react-router';
import { AccountsPageLink, LoginPath } from '../../routes';
import { AuthContext } from '../auth/authContext';
import { Logo } from './logo';

const styles = () =>
  createStyles({
    companyLogoButton: {
      textTransform:
        'none' /* Material button text defaults to all-caps; disable this. */
    },
    companyText: {
      color: 'white',
      fontSize: '4vw',
      fontStyle: 'italic',
      paddingRight: '2vw'
    },
    gridContainer: {
      height: '100%'
    },
    paper: {
      background:
        'linear-gradient(90deg, rgba(60,79,112,1) 0%, rgba(128,121,141,1) 50%, rgba(213,173,177,1) 100%)',
      height: '100%',
      paddingLeft: '2vw',
      paddingRight: '2vw'
    }
  });

interface IProps extends WithStyles<typeof styles>, RouteComponentProps<{}> {}

export const Component: React.FunctionComponent<IProps> = (props: IProps) => {
  const [anchorEl, setAnchorEl] = React.useState<HTMLElement | null>(null);
  const authContext = React.useContext(AuthContext);

  const handleLogout = () => {
    setAnchorEl(null);
    authContext.isAuthenticated = false;
    props.history.push(LoginPath);
  };

  // Don't show the account icon when the user is not signed in.
  const accountIcon = authContext.isAuthenticated ? (
    <>
      <Button
        onClick={e => {
          setAnchorEl(e.currentTarget);
        }}
      >
        <Avatar>
          <AccountCircle />
        </Avatar>
      </Button>
      <Menu
        id="simple-menu"
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={() => setAnchorEl(null)}
      >
        <MenuItem onClick={() => handleLogout()}>Logout</MenuItem>
      </Menu>
    </>
  ) : (
    <></>
  );

  return (
    <Paper square={true} className={props.classes.paper}>
      <Grid
        container={true}
        alignItems={'center'}
        justify={'space-between'}
        className={props.classes.gridContainer}
      >
        <Grid item={true}>
          <Button
            component={AccountsPageLink}
            className={props.classes.companyLogoButton}
          >
            <Typography
              variant="h4"
              inline={true}
              className={props.classes.companyText}
            >
              BridgeNational
            </Typography>
            <Logo />
          </Button>
        </Grid>
        <Grid item={true}>{accountIcon}</Grid>
      </Grid>
    </Paper>
  );
};

const RoutingAwareComponent = withRouter(Component);

export const Header = withStyles(styles)(RoutingAwareComponent);
