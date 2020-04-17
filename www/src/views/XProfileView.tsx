import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Card, Image, Radio } from "semantic-ui-react";
import { XBoundary, XCardGroup } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import XUserChangeNameModal from "../components/user/XUserChangeNameModel";
import { User } from "../graphql/models";

export const PROFILE_QUERY = gql`
  {
    me {
      ... on User {
        id
        name
        photoURL
        isActivated
        isAdmin
      }
    }
  }
`;

export type ProfileResponse = {
  me: User;
};

const XProfileView = () => {
  const {
    loading,
    error,
    data: {
      me: {
        isAdmin = false,
        isActivated = false,
        photoURL = "/app/default_profile.gif",
        name = "anonymous hippo"
      } = {}
    } = {}
  } = useQuery<ProfileResponse>(PROFILE_QUERY);

  const whenLoading = (
    <XLoadingMessage title="Loading User" msg="Fetching profile info" />
  );

  return (
    <React.Fragment>
      <XErrorMessage title="Error Loading Profile" err={error} />
      <XBoundary boundary={whenLoading} show={!loading}>
        <XCardGroup>
          <Card>
            <Image src={photoURL} wrapped ui={false} />
            <Card.Content>
              <Card.Header>{name}</Card.Header>
              <Card.Description>
                <Radio
                  label={"Activated"}
                  checked={isActivated}
                  toggle
                  type="radio"
                />
                <br />
                <Radio
                  label={"Is Admin"}
                  toggle
                  checked={isAdmin}
                  type="radio"
                  style={{ marginTop: "5px", marginBottom: "5px" }}
                />
                <XUserChangeNameModal />
              </Card.Description>
            </Card.Content>
          </Card>
        </XCardGroup>
      </XBoundary>
    </React.Fragment>
  );
};

export default XProfileView;
