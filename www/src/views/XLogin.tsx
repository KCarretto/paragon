import * as React from "react";
import { Redirect } from "react-router-dom";
import {
  Button,
  Form,
  Grid,
  Header,
  Icon,
  Image,
  Segment
} from "semantic-ui-react";

export type XLoginParams = {
  userID?: string;
  isActivated: boolean;
  isAdmin: boolean;
};

const XLogin = (props: XLoginParams) => {
  if (props.userID) {
    if (props.isActivated || props.isAdmin) {
      return <Redirect to="/" />;
    }
    return <Redirect to="/login/pending" />;
  }

  return (
    <Grid textAlign="center" style={{ height: "100vh" }} verticalAlign="middle">
      <Grid.Column style={{ maxWidth: 450 }}>
        <Header as="h2" color="blue" textAlign="center">
          <Image src="/app/logo512.png" /> Paragon Login
        </Header>
        <Form size="large">
          <Segment stacked>
            <Button
              href="/oauth/login"
              icon
              basic
              labelPosition="right"
              size="large"
            >
              <Icon name="google" />
              Login With Google
            </Button>
          </Segment>
        </Form>
      </Grid.Column>
    </Grid>
  );
};
export default XLogin;
