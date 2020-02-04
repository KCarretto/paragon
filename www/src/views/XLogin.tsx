import * as React from "react";
import { FunctionComponent } from "react";
import { useHistory, useLocation } from "react-router-dom";
import { Button, Grid, Header, Icon, Image } from "semantic-ui-react";
import { XLoadingMessage } from "../components/messages";
import { HTTP_URL } from "../config";

const XLogin: FunctionComponent<{ authorized: boolean; pending: boolean }> = ({
  authorized,
  pending
}) => {
  let location: { state?: { from: { pathname: string } } } = useLocation();
  let history = useHistory();
  let { from } = location.state || { from: { pathname: "/" } };

  if (authorized) {
    history.replace(from);
  }

  return (
    <Grid textAlign="center" style={{ height: "100vh" }} verticalAlign="middle">
      <Grid.Column style={{ maxWidth: 600 }}>
        <Header as="h2" color="blue" textAlign="center">
          <Image src="/app/logo512.png" /> Paragon Login
        </Header>
        {pending ? (
          <XLoadingMessage
            title="Pending Account Activation"
            msg="Waiting for an admin to approve your request for access..."
            hidden={false}
          />
        ) : (
          <Button
            href={HTTP_URL + "/oauth/login"}
            icon
            basic
            labelPosition="right"
            size="large"
          >
            <Icon name="google" />
            Login With Google
          </Button>
        )}
      </Grid.Column>
    </Grid>
  );
};
export default XLogin;
