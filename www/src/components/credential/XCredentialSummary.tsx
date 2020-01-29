import * as React from "react";
import { FunctionComponent } from "react";
import { Icon, Label, Table } from "semantic-ui-react";
import { Credential } from "../../graphql/models";
import { XClipboard } from "../form";

type CredentialSummaryParams = {
  credentials: Credential[];
};

const XCredentialSummary: FunctionComponent<CredentialSummaryParams> = ({
  credentials
}) => (
  <Table basic="very" collapsing padded>
    <Table.Body>
      {credentials ? (
        credentials.map((credential, index) => (
          <Table.Row key={index}>
            <Table.Cell>
              <Icon name="user" />
              <a>
                <XClipboard value={credential.principal}>
                  {credential.principal}
                </XClipboard>
              </a>
            </Table.Cell>
            <Table.Cell>
              {" "}
              <Icon name="key" />
              <a>
                <XClipboard value={credential.secret}>
                  {credential.secret}
                </XClipboard>
              </a>
            </Table.Cell>
            <Table.Cell>
              <Label>
                {!credential.fails || credential.fails < 1 ? (
                  <Icon name="times circle" color="red" size="large" />
                ) : (
                  <Icon name="check circle" color="green" size="large" />
                )}
                {credential.fails || 0} Fail(s)
              </Label>
            </Table.Cell>
          </Table.Row>
        ))
      ) : (
        <Table.Footer fullWidth>
          <Table.Row>
            <Table.HeaderCell colSpan="3">
              No credentials assigned to Target
            </Table.HeaderCell>
          </Table.Row>
        </Table.Footer>
      )}
    </Table.Body>
  </Table>
);

export default XCredentialSummary;
