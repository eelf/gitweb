import * as React from 'react';
import {grpc} from "grpc-web-client";
import {RepoListRequest, RepoListResponse} from './proto/gitweb_pb';
import {Gitweb} from "./proto/gitweb_pb_service";

class App extends React.Component<any, any> {
  constructor(props) {
    super(props);
    this.state = {
      repos: []
    };
  }

  componentDidMount() {
    grpc.invoke(Gitweb.RepoList, {
      debug: undefined,
      host: location.origin,
      request: new RepoListRequest(),
      metadata: undefined,
      transport: undefined,
      onHeaders: undefined,
      onMessage: (res: RepoListResponse) => {
        const repos = res.getReposList();
        this.setState({repos})
      },
      onEnd: (code, msg, trailers) => {
        console.log("RepoList", code)
        console.log(msg);
      },
    });
  }

  render() {
    let clientId = '434983029577-8h8cehvio5ocb7v97ajj99k84d58pqmu.apps.googleusercontent.com';
    return <div>
      <button onClick={() => location.href = 'https://accounts.google.com/o/oauth2/v2/auth' +
          '?scope=openid' +
          '&access_type=online' +
          '&include_granted_scopes=true' +
          '&response_type=token' +
          // '&state=state_parameter_passthrough_value' +
          '&redirect_uri=' + location.href + '' +
          '&client_id=' + clientId}>Login</button>
      hi
      <ul>
      {this.state.repos.map((e:RepoListResponse.repo) => <li key={e.getName()}>{e.getName()}</li>)}
      </ul>
    </div>;
  }
}

export default App;
