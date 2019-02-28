import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import {Theme} from "@material-ui/core";
import Card from "@material-ui/core/Card";
import CardActionArea from "@material-ui/core/CardActionArea";
import CardContent from "@material-ui/core/CardContent";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";
import {AccountBalanceWallet} from "@material-ui/icons";
import React from "react";
import {AccountsPageLink, TransferPageLink} from "../../routes";

const styles = (theme: Theme) => createStyles({
    card: {},
    gridContainer: {
        height: "100%", /* Force the grid to be same size as parent Paper component. */
    },
    headerText: {},
    root: {
        padding: "10px 20px",
    },
});

interface IProps extends WithStyles<typeof styles> {
    accountName: string;
    accountNumber?: number;
    accountBalance: number;
}

export const component: React.FunctionComponent<IProps> = (props: IProps) => (
    <div className={props.classes.root}>
        <Card className={props.classes.card}>
            <CardActionArea component={TransferPageLink}>
                <CardContent>
                    <Grid
                        container={true}
                        alignItems={"center"}
                        justify={"space-between"}
                        className={props.classes.gridContainer}
                    >
                        <Grid item={true}>
                            <AccountBalanceWallet/>
                        </Grid>
                        <Grid item={true}>
                            <Typography variant="subtitle1" className={props.classes.headerText}>
                                {props.accountName}
                            </Typography>
                        </Grid>
                        <Grid item={true} xs={4}/>
                        <Grid item={true}>
                            <div>
                                {"$" + props.accountBalance.toFixed(2)}
                            </div>
                            <div>
                                <Typography variant="body1">
                                    Available balance
                                </Typography>
                            </div>
                        </Grid>
                    </Grid>
                </CardContent>
            </CardActionArea>
        </Card>
    </div>
);

export const AccountCard = withStyles(styles)(component);
