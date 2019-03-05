import React from "react";
import {User} from "../../api/client";

export interface IAuthContext {
    isAuthenticated: boolean;
    user?: User;
}

const initialAuthContext: IAuthContext = {
    isAuthenticated: false,
    user: null,
};

export const AuthContext = React.createContext(initialAuthContext);
