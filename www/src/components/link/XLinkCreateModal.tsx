import { useMutation } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import moment from "moment";
import * as React from "react";
import { useState } from "react";
import {
  Button,
  Dropdown,
  Form,
  Grid,
  Input,
  Message,
  Modal
} from "semantic-ui-react";
import { MULTI_FILE_QUERY } from "../../views/XMultiFileView";
import { useModal } from "../form";

export const CREATE_LINK_MUTATION = gql`
  mutation createLink(
    $alias: String!
    $expirationTime: Time
    $clicks: Int
    $file: ID!
  ) {
    createLink(
      input: {
        alias: $alias
        expirationTime: $expirationTime
        clicks: $clicks
        file: $file
      }
    ) {
      id
    }
  }
`;

type CreateLinkModalParams = {
  openOnStart?: boolean;
  file: String;
};

const XCreateLinkModal = ({ openOnStart, file }: CreateLinkModalParams) => {
  const [openModal, closeModal, isOpen] = useModal();
  const [error, setError] = useState<ApolloError>(null);

  // Form params
  const [alias, setAlias] = useState<string>("");
  const [expirationTime, setExpirationTime] = useState<string>(null);
  const [clicks, setClicks] = useState<number>(null);

  const [createLink, { called, loading }] = useMutation(CREATE_LINK_MUTATION, {
    refetchQueries: [{ query: MULTI_FILE_QUERY }]
  });

  const handleSubmit = () => {
    let vars = {
      alias: alias,
      expirationTime: expirationTime,
      clicks: clicks,
      file: file
    };

    createLink({
      variables: vars
    })
      .then(({ data, errors }) => {
        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          setError(e);
          return;
        }
        closeModal();
      })
      .catch(err => setError(err));
  };

  if (openOnStart) {
    openModal();
  }

  const expiryOptions = [
    {
      key: null,
      text: "Never",
      value: null
    },
    {
      key: "5 minutes",
      text: "5 minutes",
      value: moment()
        .add(5, "minutes")
        .format("YYYY-MM-DDTHH:mm:ssZ")
    },
    {
      key: "10 minutes",
      text: "10 minutes",
      value: moment()
        .add(10, "minutes")
        .toString()
    },
    {
      key: "15 minutes",
      text: "15 minutes",
      value: moment()
        .utc()
        .add(15, "minutes")
        .format("YYYY-MM-DDTHH:mm:ssZ")
    },
    {
      key: "30 minutes",
      text: "30 minutes",
      value: moment()
        .add(30, "minutes")
        .format("YYYY-MM-DDTHH:mm:ssZ")
    },
    {
      key: "1 hour",
      text: "1 hour",
      value: moment()
        .add(1, "hour")
        .format("YYYY-MM-DDTHH:mm:ssZ")
    },
    {
      key: "2 hours",
      text: "2 hours",
      value: moment()
        .add(2, "hour")
        .format("YYYY-MM-DDTHH:mm:ssZ")
    }
  ];

  return (
    <Modal
      open={isOpen}
      onClose={closeModal}
      trigger={<Button basic color="blue" icon="linkify" onClick={openModal} />}
      size="large"
      // Form properties
      as={Form}
      onSubmit={handleSubmit}
      error={called && error}
      loading={called && loading}
    >
      <Modal.Header>{"Create a Link"}</Modal.Header>
      <Modal.Content>
        <Grid verticalAlign="middle" stackable container columns={"equal"}>
          <Grid.Column>
            <Input
              fluid
              placeholder="Enter link alias(what will be after '/l/' in the URL)"
              name="alias"
              value={alias}
              onChange={(e, { value }) => setAlias(value)}
            />
          </Grid.Column>

          <Grid.Column>
            <Dropdown
              label="Link Expiration Time"
              placeholder="Never"
              value={expirationTime === null ? null : expirationTime.toString()}
              fluid
              selection
              options={expiryOptions}
              onChange={(e, selection) =>
                setExpirationTime(selection.value.toString())
              }
            />
          </Grid.Column>

          <Grid.Column>
            <Input
              fluid
              placeholder="Unlimited Clicks"
              name="clicks"
              value={clicks}
              onChange={(e, { value }) => setClicks(Number(value))}
            />
          </Grid.Column>

          <Grid.Column>
            {/* <XTargetTypeahead
              labeled
              onChange={(e, { value }) => setTargets(value)}
            /> */}
          </Grid.Column>
        </Grid>

        <Message
          error
          icon="warning"
          header={"Failed to Queue Job"}
          onDismiss={(e, data) => setError(null)}
          content={error ? error.message : "Unknown Error"}
        />
      </Modal.Content>
      <Modal.Actions>
        <Form.Button style={{ marginBottom: "10px" }} positive floated="right">
          Link!
        </Form.Button>
      </Modal.Actions>
    </Modal>
  );
};

export default XCreateLinkModal;
