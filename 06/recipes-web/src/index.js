import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import './index.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.min.js';
import { Auth0Provider } from "@auth0/auth0-react";

ReactDOM.render(
    <Auth0Provider
        domain="dev-3r7xgdez.us.auth0.com"
        clientId="42mdjceC3s0YXyTob2inFlelaIAL9t7S"
        redirectUri={window.location.origin}
    >
        <App />
    </Auth0Provider>,
    document.getElementById("root")
);
