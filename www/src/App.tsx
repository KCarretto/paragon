import * as React from "react";
import { useEffect, useState } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import "semantic-ui-css/semantic.min.css";
import "./App.css";
import { XLayout, XPrivateRoute, XUnimplemented } from "./components/layout";
import { HTTP_URL } from "./config";
import { Routes } from "./config/routes";
import { XGraphProvider } from "./graphql";
import {
  XJobView,
  XLogin,
  XMultiFileView,
  XMultiJobView,
  XMultiTagView,
  XMultiTargetView,
  XTargetView,
  XTaskView
} from "./views";

type StatusResult = {
  activated: boolean;
  userid?: number;
  is_authenticated: boolean;
  is_activated: boolean;
  is_admin: boolean;
};

// localStorage

// localStorage.setItem("user.IsActivated", "false")
// localStorage.setItem("user.IsAdmin", "false")
// localStorage.setItem("user.IsAdmin", "false")

const App = () => {
  const [userID, setUserID] = useState<string>(null);
  const [authenticated, setAuthenticated] = useState(false);
  const [activated, setActivated] = useState(false);
  const [admin, setAdmin] = useState(false);

  const fetchUserInfo = () => {
    fetch(HTTP_URL + "/status").then(
      resp =>
        resp.json().then(
          (data: StatusResult) => {
            setUserID(data.userid ? String(data.userid) : null);
            setAuthenticated(data.is_authenticated || false);
            setActivated(data.is_activated || false);
            setAdmin(data.is_admin || false);
          },
          err => console.error("failed to parse response json", err)
        ),
      err => console.error("failed to request status", err)
    );
  };

  useEffect(fetchUserInfo, []);

  let authz = authenticated && (activated || admin);

  return (
    <XGraphProvider>
      <Router>
        <Switch>
          <Route path="/login">
            <XLogin authorized={authz} pending={authenticated} />
          </Route>
          <XLayout
            routeMap={Routes}
            userID={userID}
            isAdmin={admin}
            className="App"
          >
            <XPrivateRoute
              authorized={authz}
              exact
              path="/"
              component={XUnimplemented}
            />
            <XPrivateRoute
              authorized={authz}
              exact
              path="/news_feed"
              component={XUnimplemented}
            />
            <XPrivateRoute
              authorized={authz}
              exact
              path="/profile"
              component={XUnimplemented}
            />

            <XPrivateRoute
              authorized={authz}
              exact
              path="/targets"
              component={XMultiTargetView}
            />
            <XPrivateRoute
              authorized={authz}
              exact
              path="/jobs"
              component={XMultiJobView}
            />
            <XPrivateRoute
              authorized={authz}
              exact
              path="/tags"
              component={XMultiTagView}
            />
            <XPrivateRoute
              authorized={authz}
              exact
              path="/files"
              component={XMultiFileView}
            />

            <XPrivateRoute
              authorized={authz}
              path="/targets/:id"
              component={XTargetView}
            />
            <XPrivateRoute
              authorized={authz}
              path="/tasks/:id"
              component={XTaskView}
            />
            <XPrivateRoute
              authorized={authz}
              path="/jobs/:id"
              component={XJobView}
            />
          </XLayout>
        </Switch>
      </Router>
    </XGraphProvider>
  );
};

export default App;
