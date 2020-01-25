import * as React from "react";
import { Divider, Feed, Header, Icon, List } from "semantic-ui-react";
import { XCredentialStatus } from ".";
import { Credential } from "../../graphql/models";

type CredentialSummaryParams = {
  credentials: Credential[];
};

const XCredentialSummary = (credentials: Credential[]) => (
  <Feed>
    <Header sub>Credentials</Header>
    {credentials ? (
      credentials.map((credential, index) => (
        <Feed.Event key={index}>
          <Feed.Label>
            <Icon
              fitted
              size="big"
              {...new XCredentialStatus().getStatus(credential).icon}
            />
          </Feed.Label>
          <Feed.Content>
            <Feed.Summary>
              <List.Header>
                Principal: {credential.principal} <br />
                Secret: {credential.secret}
              </List.Header>
            </Feed.Summary>
            <Feed.Extra text>
              {new XCredentialStatus().getStatus(credential).text}
            </Feed.Extra>
            <Feed.Meta>Fails: {credential.fails}</Feed.Meta>
            <Divider />
          </Feed.Content>
        </Feed.Event>
      ))
    ) : (
      <Header sub disabled>
        No credentials assigned to Target
      </Header>
    )}
  </Feed>
);

export default XCredentialSummary;
