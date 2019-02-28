import React from "react";
import {Redirect, Route, RouteComponentProps, RouteProps, Switch} from "react-router";
import {BrowserRouter, Link} from "react-router-dom";
import {AuthContext, IAuthContext} from "../components/auth/authContext";
import {AccountsView} from "../views/accounts";
import {LoginView} from "../views/login";
import {NotFoundView} from "../views/notfound";
import {TransferView} from "../views/transfer";
import {TransactionsView} from "../views/transactions";

// Single location for all paths.
export const AccountsPath = "/accounts";
export const LoginPath = "/login";
export const TransferPath = "/transfer";
export const TransactionsPath = "/transactions";


// Convenient Link components to pass as the ButtonBase#component prop to navigation Buttons, e.g.
//     <Button component={AccountsPageLink}>Button text</Button>
export const AccountsPageLink: React.FunctionComponent<{}> = (props) => <Link to={AccountsPath} {...props}/>;
export const TransferPageLink: React.FunctionComponent<{}> = (props) => <Link to={TransferPath} {...props}/>;
export const TransactionsPageLink: React.FunctionComponent<{}> = (props) => <Link to={TransactionsPath} {...props}/>;

interface IPrivateRouteProps extends RouteProps {
    component: React.FunctionComponent<any>;
}

// Helper component for auto-redirecting unauthenticated users to the login page.
const PrivateRoute: React.FunctionComponent<IPrivateRouteProps> = ({component: Component, ...rest}) => {
    const authContext: IAuthContext = React.useContext(AuthContext);
    const renderFunc = (props: RouteComponentProps<any>): React.ReactNode => (
        authContext.isAuthenticated
            ? (<Component {...props}/>)
            : (<Redirect to={{pathname: LoginPath, state: {from: props.location}}}/>));
    return (
        <Route {...rest} render={renderFunc}/>
    );
};

// Route-to-component table. Should only be used as the root component in src/index.tsx.
export const ViewsRouter: React.FunctionComponent<{}> = () => {
    return (
        <BrowserRouter>
            <Switch>
                <PrivateRoute
                    path={AccountsPath}
                    exact={true}
                    component={AccountsView}
                />
                <PrivateRoute
                    path={TransferPath}
                    component={TransferView}
                />
                <PrivateRoute
                    path={`${TransferPath}/:accountId`}
                    component={TransferView}
                />
                <PrivateRoute
                    path={TransactionsPath}
                    component={TransactionsView}
                />
                <Route path={LoginPath} exact={true} component={LoginView}/>
                <Route path={"/"} exact={true} component={LoginView}/>
                <Route component={NotFoundView}/>
            </Switch>
        </BrowserRouter>
    );
}
