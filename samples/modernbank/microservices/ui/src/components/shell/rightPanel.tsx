import { createStyles, WithStyles, withStyles } from '@material-ui/core';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import { AttachMoney } from '@material-ui/icons';
import React from 'react';

const styles = () =>
  createStyles({
    gridContainer: {
      height: '100%',
      paddingBottom: '5vh',
      paddingLeft: '2vw',
      paddingRight: '2vw',
      paddingTop: '5vh'
    },
    headerText: {
      color: 'rgb(246, 193, 118)'
    },
    paper: {
      backgroundColor: 'rgba(0,0,0,0.7)',
      borderTopColor: 'rgb(233,121,51)',
      borderTopStyle: 'solid',
      borderTopWidth: '0.5vh',
      boxSizing: 'border-box',
      height: '100%',
      textAlign: 'center'
    },
    subheaderText: {
      color: 'white'
    }
  });

interface IProps extends WithStyles<typeof styles> {}

export const component: React.FunctionComponent<IProps> = (props: IProps) => {
  return (
    <Paper square={true} className={props.classes.paper}>
      <Grid
        container={true}
        alignItems={'center'}
        direction={'column'}
        justify={'space-between'}
        className={props.classes.gridContainer}
      >
        <Grid item={true}>
          <AttachMoney color={'primary'} fontSize="large" />
        </Grid>
        <Grid item={true}>
          <Typography variant={'h4'} className={props.classes.headerText}>
            Need cash?
          </Typography>
        </Grid>
        <Grid item={true}>
          <Typography
            variant={'body1'}
            align={'center'}
            className={props.classes.subheaderText}
          >
            Why not transfer yourself some through our completely functional
            online demo application? It is fast and easy and completely
            uninsured.
          </Typography>
        </Grid>
        <Grid item={true} xs={6} />
      </Grid>
    </Paper>
  );
};

export const RightPanel = withStyles(styles)(component);
