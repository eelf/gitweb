import * as React from 'react';
import * as ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {grpc} from "grpc-web-client";
import {Gitweb} from "./proto/gitweb_pb_service";
import {AuthRequest, AuthResponse} from "./proto/gitweb_pb";

const h = new URLSearchParams(location.hash.substr(1))
const at = h.get('access_token');
if (at) {
    const req = new AuthRequest();
    req.setGoogleAccessToken(at);
    grpc.invoke(Gitweb.Auth, {
        host: location.origin,
        request: req,
        onMessage: (res: AuthResponse) => {
            console.log('auth resp', res.getText());
        },
        onEnd: (code, msg, trailers) => {
            console.log('auth code', code, msg);
        },
    });
}

ReactDOM.render(
  <App />,
  document.getElementById('root')
);
