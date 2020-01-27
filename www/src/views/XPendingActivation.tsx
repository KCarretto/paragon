import Cookies from "js-cookie";
import * as React from "react";
import { Redirect } from "react-router-dom";
import { Grid, Header, Image, Segment } from "semantic-ui-react";
import { XLoadingMessage } from "../components/messages";

export type XPendingActivationParams = {
  userID?: number;
  isActivated: boolean;
  isAdmin: boolean;
};

const XPendingActivation = (props: XPendingActivationParams) => {
  let userIDCookie = Cookies.get("pg-userid");

  if (!props.userID && !userIDCookie) {
    return <Redirect to="/login" />;
  }

  if (props.isActivated || props.isAdmin) {
    return <Redirect to="/" />;
  }

  return (
    <Grid textAlign="center" style={{ height: "100vh" }} verticalAlign="middle">
      <Grid.Column style={{ maxWidth: 450 }}>
        <Header as="h2" color="blue" textAlign="center">
          <Image src="/app/logo512.png" />
          Registration Requested
        </Header>
        <Segment stacked>
          <XLoadingMessage
            title="Pending Account Activation"
            msg="Waiting for an admin to approve your request for access..."
            hidden={false}
          />
        </Segment>
      </Grid.Column>
    </Grid>
  );
};
export default XPendingActivation;
