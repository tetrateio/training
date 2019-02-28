import React from "react";

export interface IAuthContext {
    isAuthenticated: boolean;
    username?: string;
}

const initialAuthContext = {
    isAuthenticated: false,
};

export const AuthContext = React.createContext(initialAuthContext);