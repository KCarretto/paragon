import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Card, Image, Loader, Radio } from "semantic-ui-react";
import { XCardGroup } from "../components/layout";
import { XErrorMessage } from "../components/messages";
import XUserChangeNameModal from "../components/user/XUserChangeNameModel";
import { User } from "../graphql/models";

export const PROFILE_QUERY = gql`
  {
    me {
      id
      name
      photoURL
      isActivated
      isAdmin
    }
  }
`;

export type ProfileResponse = {
  me: User;
};

const XProfileView = () => {
  const { called, loading, data, error } = useQuery<ProfileResponse>(
    PROFILE_QUERY
  );

  return (
    <XCardGroup>
      <Loader disabled={!called || !loading} />
      {!data || !data.me ? (
        "hiii" + <span />
      ) : (
        <Card>
          {data.me.photoURL !== "" ? (
            <Image src={data.me.photoURL} wrapped ui={false} />
          ) : (
            <Image src="/app/default_profile.gif" wrapped ui={false} />
          )}
          <Card.Content>
            <Card.Header>{data.me.name}</Card.Header>
            <Card.Description>
              <Radio
                label={"Activated"}
                checked={data.me.isActivated}
                toggle
                type="radio"
              />
              <br />
              <Radio
                label={"Is Admin"}
                toggle
                checked={data.me.isAdmin}
                type="radio"
              />
              <XUserChangeNameModal />
            </Card.Description>
          </Card.Content>
        </Card>
      )}
      <XErrorMessage title="Error Loading Profile" err={error} />
    </XCardGroup>
  );
};

export default XProfileView;
