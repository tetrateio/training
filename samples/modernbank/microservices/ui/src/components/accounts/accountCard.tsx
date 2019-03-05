import {createStyles, WithStyles, withStyles} from "@material-ui/core";
import Card from "@material-ui/core/Card";
import CardActionArea from "@material-ui/core/CardActionArea";
import CardContent from "@material-ui/core/CardContent";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";
import {AccountBalanceWallet} from "@material-ui/icons";
import React from "react";
import {AccountsPageLink, transferPageLink} from "../../routes";

const styles = () => createStyles({
    card: {},
    centerGridItem: {
        textAlign: "center",
    },
    gridContainer: {
        height: "100%",
    },
    headerText: {
    },
    root: {
        paddingBottom: "1vh",
        paddingLeft: "2vw",
        paddingRight: "2vw",
        paddingTop: "1vh",
    },
    subheaderText: {
    },
});

interface IProps extends WithStyles<typeof styles> {
    accountName: string;
    accountNumber?: number;
    accountBalance: number;
}

export const Component: React.FunctionComponent<IProps> = (props: IProps) => {
    const redirectLink =
        (!!props.accountNumber)
            ? transferPageLink(props.accountNumber.toString())
            : AccountsPageLink;
    return (
        <div className={props.classes.root}>
            <Card className={props.classes.card}>
                <CardActionArea component={redirectLink}>
                    <CardContent>
                        <Grid
                            container={true}
                            alignItems={"center"}
                            justify={"space-between"}
                            className={props.classes.gridContainer}
                        >
                            <Grid item={true} xs={1}>
                                <AccountBalanceWallet/>
                            </Grid>
                            <Grid item={true} xs={3} className={props.classes.centerGridItem}>
                                <Typography
                                    variant="subtitle1"
                                    className={props.classes.headerText}
                                >
                                    {props.accountName}
                                </Typography>
                                <Typography
                                    variant={"body1"}
                                    className={props.classes.subheaderText}
                                >
                                    {props.accountNumber}
                                </Typography>
                            </Grid>
                            <Grid item={true} xs={4}>
                            </Grid>
                            <Grid item={true} xs={2}>
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
};

export const AccountCard = withStyles(styles)(Component);
