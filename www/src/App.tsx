import * as React from "react";
import { useEffect, useState } from "react";
import {
  BrowserRouter as Router,
  Redirect,
  Route,
  Switch
} from "react-router-dom";
import "semantic-ui-css/semantic.min.css";
import "./App.css";
import { XLayout, XUnimplemented } from "./components/layout";
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
    fetch(window.location.origin + "/status").then(
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

  let showLogin = !authenticated || (!activated && !admin);

  return (
    <XGraphProvider>
      <Router>
        <Switch>
          <Route
            path="/login"
            render={routeProps =>
              showLogin ? (
                <XLogin pending={authenticated} />
              ) : (
                <Redirect to="/" />
              )
            }
          />
          {showLogin ? (
            <Redirect to="/login" />
          ) : (
            <XLayout
              routeMap={Routes}
              userID={userID}
              isAdmin={admin}
              className="App"
            >
              <Route exact path="/" component={XUnimplemented} />
              <Route exact path="/news_feed" component={XUnimplemented} />
              <Route exact path="/profile" component={XUnimplemented} />

              <Route exact path="/targets" component={XMultiTargetView} />
              <Route exact path="/jobs" component={XMultiJobView} />
              <Route exact path="/tags" component={XMultiTagView} />
              <Route exact path="/files" component={XMultiFileView} />

              <Route path="/targets/:id" component={XTargetView} />
              <Route path="/tasks/:id" component={XTaskView} />
              <Route path="/jobs/:id" component={XJobView} />
            </XLayout>
          )}
        </Switch>
      </Router>
    </XGraphProvider>
  );
};

export default App;
