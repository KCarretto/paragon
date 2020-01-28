import * as React from "react";
import { Button, Grid, Header, Icon, Image } from "semantic-ui-react";
import { XLoadingMessage } from "../components/messages";
import { HTTP_URL } from "../config";

export type XLoginParams = {
  pending: boolean;
};

const XLogin = (props: XLoginParams) => (
  <Grid textAlign="center" style={{ height: "100vh" }} verticalAlign="middle">
    <Grid.Column style={{ maxWidth: 600 }}>
      <Header as="h2" color="blue" textAlign="center">
        <Image src="/app/logo512.png" /> Paragon Login
      </Header>
      {props.pending ? (
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
export default XLogin;
