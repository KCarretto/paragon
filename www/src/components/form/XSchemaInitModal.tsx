import * as React from "react";
import { useState } from "react";
import {
  Accordion,
  Button,
  ButtonProps,
  Form,
  Icon,
  Message,
  Modal,
  Tab,
  Table,
  TextArea
} from "semantic-ui-react";
import { useModal } from ".";
import { FETCH_CORS, HTTP_URL } from "../../config";
import { XCredentialSummary } from "../credential";

/* EXAMPLE SCHEMA:
[
       {
            "name": "TestTarget1",
            "os": "LINUX",
            "primaryIP": "127.0.0.1",
            "tags":
                [
                    "windows2",
                    "bop"
                ],
            "credentials":
                [
                    {
                        "kind": "password",
                        "principal": "root",
                        "secret": "changeme"
                    }
                ]
       },
       {
            "name": "TestTarget2",
            "os": "LINUX",
            "primaryIP": "127.0.0.2",
            "tags":
                [
                    "windows2",
                    "bop",
                    "shmoop"
                ],
            "credentials":
                [
                    {
                        "kind": "password",
                        "principal": "root",
                        "secret": "changeme"
                    },
                    {
                        "kind": "password",
                        "principal": "root",
                        "secret": "changeme"
                    }
                ]
       }
]
*/

const VALID_OS = ["LINUX", "WINDOWS", "BSD", "MACOS"];
const VALID_CRED_KIND = ["password"];

const validateSchema: (schema: string) => string = (schema) => {
  let targets = JSON.parse(schema);
  if (!(targets instanceof Array)) {
    throw new TypeError("Provided JSON is not an array of objects");
  }

  let newTargets = targets.map((target, index) => {
    // Name
    if (!target.name || target.name === "") {
      throw new TypeError(
        `Target ${index}: No value for 'name' specified, must provide non-empty string.`
      );
    }

    // OS
    if (!VALID_OS.includes(target.os)) {
      throw new TypeError(
        `Target ${index} (${target.name}): Invalid value for 'os': '${
          target.os
        }', must provide one of ${VALID_OS.join(", ")}.`
      );
    }

    // Primary IP
    if (!target.primaryIP || target.primaryIP === "") {
      throw new TypeError(
        `Target ${index} (${target.name}): Invalid value for 'primaryIP', must provide non-empty string.`
      );
    }

    // Credentials
    if (target.credentials) {
      if (!(target.credentials instanceof Array)) {
        throw new TypeError(
          `Target ${index} (${target.name}): Invalid value for 'credentials', must provide array of JSON objects.`
        );
      }

      target.credentials = target.credentials.map((cred, credIndex) => {
        if (!VALID_CRED_KIND.includes(cred.kind)) {
          throw new TypeError(
            `Target ${index} (${
              target.name
            }): Invalid credential at index ${credIndex}: Invalid value for 'kind', must provide one of ${VALID_CRED_KIND.join(
              ", "
            )}.`
          );
        }
        return cred;
      });
    }
    return target;
  });

  return JSON.stringify(newTargets, null, 2);
};

const renderOSIcon = (os: string) => {
  switch (os) {
    case "WINDOWS":
      return <Icon name="windows" />;
    case "LINUX":
      return <Icon name="linux" />;
    case "BSD":
      return <Icon name="freebsd" />;
    case "MACOS":
      return <Icon name="apple" />;
  }

  return <Icon name="question circle outline" />;
};

const renderTable = (schema: string) => {
  try {
    let targets = JSON.parse(schema);
    if (!targets || !(targets instanceof Array) || targets.length < 1) {
      return <div />;
    }

    return (
      <Table>
        <Table.Header>
          <Table.HeaderCell>Name</Table.HeaderCell>
          <Table.HeaderCell>PrimaryIP</Table.HeaderCell>
          <Table.HeaderCell>PublicIP</Table.HeaderCell>
          <Table.HeaderCell>OS</Table.HeaderCell>
          <Table.HeaderCell>Tags</Table.HeaderCell>
          <Table.HeaderCell>Credentials</Table.HeaderCell>
        </Table.Header>
        {targets
          .filter((target) => target)
          .map((target) => (
            <Table.Row>
              <Table.Cell>{target.name}</Table.Cell>
              <Table.Cell>{target.primaryIP}</Table.Cell>
              <Table.Cell>{target.publicIP}</Table.Cell>
              <Table.Cell>{renderOSIcon(target.os)}</Table.Cell>
              <Table.Cell>
                {target.tags && target.tags.length > 0
                  ? target.tags.join(", ")
                  : "No tags"}
              </Table.Cell>
              <Table.Cell>
                <XCredentialSummary credentials={target.credentials} />
              </Table.Cell>
            </Table.Row>
          ))}
      </Table>
    );
  } catch (err) {
    console.error("Failed to render schema preview table", err);
  }
  return <p>No Preview Available ¯\_(ツ)_/¯</p>;
};

const XSchemaInitModal: React.FC<{
  button: ButtonProps;
}> = ({ button }) => {
  const [openModal, closeModal, isOpen] = useModal();
  const [activeIndex, setActiveIndex] = useState<number>(null);

  const [error, setError] = useState<string>(null);
  const [schema, setSchema] = useState<string>("");

  const handleSubmit = () => {
    fetch(HTTP_URL + "/init/", {
      method: "POST",
      mode: FETCH_CORS,
      headers: {
        "Content-Type": "application/json",
      },
      body: schema,
    })
      .then((resp) => {
        if (resp.status === 200) {
          closeModal();
          return;
        }

        setError(
          `Failed to Initialize Teamserver: HTTP Error ${resp.statusText}`
        );
      })
      .catch((err) => setError("Failed to Initialize Teamserver: " + err));
  };

  return (
    <Modal
      open={isOpen}
      onClose={closeModal}
      trigger={<Button onClick={openModal} {...button} />}
      size="small"
      // Form properties
      as={Form}
      onSubmit={handleSubmit}
      error={!!error}
    >
      <Modal.Header>Initialize Teamserver</Modal.Header>
      <Modal.Content>
        <Tab
          panes={[
            {
              menuItem: "JSON",
              render: () => (
                <Tab.Pane>
                  <TextArea
                    rows={10}
                    placeholder={"Enter JSON Schema..."}
                    value={schema}
                    onChange={(e, { value }) => {
                      let val = value.toString();
                      try {
                        val = validateSchema(val);
                        setError(null);
                      } catch (err) {
                        setError("Invalid JSON array: " + err);
                      }
                      setSchema(val);
                    }}
                  />
                </Tab.Pane>
              ),
            },
            {
              menuItem: "Preview",
              render: () => <Tab.Pane>{renderTable(schema)}</Tab.Pane>,
            },
          ]}
        />
        <Message
          error
          icon="warning"
          header={"Schema Error"}
          onDismiss={(e, data) => setError(null)}
          content={error ? error : "Unknown Error"}
        />

        <Accordion>
          <Accordion.Title
            index={0}
            active={activeIndex === 0}
            onClick={() => setActiveIndex(activeIndex === 0 ? null : 0)}
          >
            <Icon name="dropdown" />
            <a>JSON Format Help</a>
          </Accordion.Title>
          <Accordion.Content active={activeIndex === 0}>
            <i>
              You may initialize the teamserver's database with a list of JSON
              objects (one per target) containing required target information,
              associated tags and credentials. The provided tags and credentials
              will be created for you if they do not already exist.
            </i>
            <br />
            <br />
            <b>Example:</b>
            <br />
            <pre style={{ backgroundColor: "#f5f5f5" }}>
              {"[\n\t{\n"}
              {'\t\t"name": "Human-Friendly Target Name",\n'}
              {'\t\t"primaryIP": "10.0.0.1",\n'}
              {'\t\t"os": "LINUX",\n'}
              {'\t\t"tags": ["Linux", "Web", "HTTP"],\n'}
              {'\t\t"credentials": [\n\t\t\t{\n'}
              {'\t\t\t\t"kind": "password",\n'}
              {'\t\t\t\t"principal": "root",\n'}
              {'\t\t\t\t"secret": "changeme",\n'}
              {"\t\t\t},\n\t\t],"}
              {"\n\t}\n]"}
            </pre>
          </Accordion.Content>
        </Accordion>
      </Modal.Content>
      <Modal.Actions>
        <Form.Button style={{ marginBottom: "10px" }} positive floated="right">
          Create
        </Form.Button>
      </Modal.Actions>
    </Modal>
  );
};

export default XSchemaInitModal;
