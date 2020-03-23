import { useMutation, useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import {
  Accordion,
  Card,
  Divider,
  Header,
  Icon,
  Image,
  Radio
} from "semantic-ui-react";
import XClipboard from "../components/form/XClipboard";
import { XBoundary, XCardGroup } from "../components/layout";
import { XErrorMessage, XLoadingMessage } from "../components/messages";
import { Service, User } from "../graphql/models";

export const ADMIN_USERS_QUERY = gql`
  {
    users {
      id
      name
      photoURL
      isActivated
      isAdmin
    }
  }
`;

export const ACTIVATE_USER_MUTATION = gql`
  mutation activateUser($id: ID!) {
    activateUser(input: { id: $id }) {
      id
      name
      photoURL
      isActivated
      isAdmin
    }
  }
`;

export const DEACTIVATE_USER_MUTATION = gql`
  mutation deactivateUser($id: ID!) {
    deactivateUser(input: { id: $id }) {
      id
      name
      photoURL
      isActivated
      isAdmin
    }
  }
`;

export const MAKE_ADMIN_MUTATION = gql`
  mutation makeAdmin($id: ID!) {
    makeAdmin(input: { id: $id }) {
      id
      name
      photoURL
      isActivated
      isAdmin
    }
  }
`;

export const REMOVE_ADMIN_MUTATION = gql`
  mutation removeAdmin($id: ID!) {
    removeAdmin(input: { id: $id }) {
      id
      name
      photoURL
      isActivated
      isAdmin
    }
  }
`;

export const ADMIN_SERVICES_QUERY = gql`
  {
    services {
      id
      name
      pubKey
      isActivated
    }
  }
`;

export const ACTIVATE_SERVICE_MUTATION = gql`
  mutation activateService($id: ID!) {
    activateService(input: { id: $id }) {
      id
      name
      isActivated
    }
  }
`;

export const DEACTIVATE_SERVICE_MUTATION = gql`
  mutation deactivateService($id: ID!) {
    deactivateService(input: { id: $id }) {
      id
      name
      isActivated
    }
  }
`;

export type UsersResponse = {
  users: User[];
};
export type ServicesResponse = {
  services: Service[];
};

const XAdminView = () => {
  const [openUsers, setOpenUsers] = useState<boolean>(false);
  const [openServices, setOpenServices] = useState<boolean>(false);

  const {
    loading: userLoading,
    error: userError,
    data: { users = [] } = {}
  } = useQuery<UsersResponse>(ADMIN_USERS_QUERY);
  const {
    loading: svcLoading,
    error: svcError,
    data: { services = [] } = {}
  } = useQuery<ServicesResponse>(ADMIN_SERVICES_QUERY);

  const [activateUser, { error: activateUserError }] = useMutation(
    ACTIVATE_USER_MUTATION,
    {
      refetchQueries: [{ query: ADMIN_USERS_QUERY }]
    }
  );
  const [deactivateUser, { error: deactivateUserError }] = useMutation(
    DEACTIVATE_USER_MUTATION,
    {
      refetchQueries: [{ query: ADMIN_USERS_QUERY }]
    }
  );
  const [activateService, { error: activateServiceError }] = useMutation(
    ACTIVATE_SERVICE_MUTATION,
    {
      refetchQueries: [{ query: ADMIN_SERVICES_QUERY }]
    }
  );
  const [deactivateService, { error: deactivateServiceError }] = useMutation(
    DEACTIVATE_SERVICE_MUTATION,
    {
      refetchQueries: [{ query: ADMIN_SERVICES_QUERY }]
    }
  );
  const [makeAdmin, { error: makeAdminError }] = useMutation(
    MAKE_ADMIN_MUTATION,
    {
      refetchQueries: [{ query: ADMIN_USERS_QUERY }]
    }
  );
  const [removeAdmin, { error: removeAdminError }] = useMutation(
    REMOVE_ADMIN_MUTATION,
    {
      refetchQueries: [{ query: ADMIN_USERS_QUERY }]
    }
  );

  const handleClick = (e, titleProps) => {
    const { index } = titleProps;
    if (index === 0) {
      setOpenUsers(!openUsers);
    } else {
      setOpenServices(!openServices);
    }
  };

  const handleUserActivate = (e, p) => {
    p.checked
      ? activateUser({
          variables: { id: p.value }
        })
      : deactivateUser({
          variables: { id: p.value }
        });
  };

  const handleUserAdmin = (e, p) => {
    p.checked
      ? makeAdmin({ variables: { id: p.value } })
      : removeAdmin({ variables: { id: p.value } });
  };

  const handleServiceActivate = (e, p) => {
    p.checked
      ? activateService({ variables: { id: p.value } })
      : deactivateService({ variables: { id: p.value } });
  };

  const showUsers = () => {
    return (
      <React.Fragment>
        <Accordion.Title active={openUsers} index={0} onClick={handleClick}>
          <Header icon as="h2" textAlign="center">
            <Icon name="dropdown" />
            <Header.Content>Users</Header.Content>
          </Header>
        </Accordion.Title>
        <Accordion.Content active={openUsers}>
          <XErrorMessage title="Error Loading Users" err={userError} />
          <XCardGroup>
            {users.map(u => {
              return (
                <Card key={u.id}>
                  {u.photoURL !== "" ? (
                    <Image src={u.photoURL} wrapped ui={false} />
                  ) : (
                    <Image src="/app/default_profile.gif" wrapped ui={false} />
                  )}
                  <Card.Content>
                    <Card.Header>{u.name}</Card.Header>
                    <Card.Description>
                      <Radio
                        label={"Activated"}
                        value={u.id}
                        onClick={handleUserActivate}
                        checked={u.isActivated}
                        toggle
                        type="radio"
                      />
                      <br />
                      <Radio
                        label={"Is Admin"}
                        toggle
                        value={u.id}
                        onClick={handleUserAdmin}
                        checked={u.isAdmin}
                        type="radio"
                        style={{ marginTop: "5px" }}
                      />
                      <XErrorMessage
                        title="Error Updating User"
                        err={activateUserError || deactivateUserError}
                      />
                      <XErrorMessage
                        title="Error Updating Admin"
                        err={makeAdminError || removeAdminError}
                      />
                    </Card.Description>
                  </Card.Content>
                </Card>
              );
            })}
          </XCardGroup>
        </Accordion.Content>
        <Divider />
      </React.Fragment>
    );
  };

  const showServices = () => {
    return (
      <React.Fragment>
        <Accordion.Title
          inverted
          active={openServices}
          index={1}
          onClick={handleClick}
        >
          <Header icon as="h2" textAlign="center">
            <Icon name="dropdown" />
            <Header.Content>Services</Header.Content>
          </Header>
        </Accordion.Title>

        <Accordion.Content active={openServices}>
          <XErrorMessage title="Error Loading Services" err={svcError} />
          <XCardGroup>
            {services.map(s => {
              return (
                <Card key={s.id}>
                  <Image src="/app/default_profile.gif" wrapped ui={false} />
                  <Card.Content>
                    <Card.Header>{s.name}</Card.Header>
                    <Card.Description>
                      <Radio
                        label={"Activated"}
                        value={s.id}
                        onClick={handleServiceActivate}
                        checked={s.isActivated}
                        toggle
                        type="radio"
                      />
                      <XErrorMessage
                        title="Error Updating Service"
                        err={activateServiceError || deactivateServiceError}
                      />
                    </Card.Description>
                  </Card.Content>
                  <Card.Content extra>
                    <a>
                      <Icon name="key" />
                      <XClipboard value={s.pubKey}>
                        Public Key: {s.pubKey.slice(0, 15)}...
                      </XClipboard>
                    </a>
                  </Card.Content>
                </Card>
              );
            })}
          </XCardGroup>
        </Accordion.Content>
        <Divider />
      </React.Fragment>
    );
  };

  const whenSVCLoading = (
    <XLoadingMessage title="Loading Services" msg="Fetching Service Info" />
  );
  const whenSVCEmpty = <h1>No Services Loaded</h1>;

  const whenUserLoading = (
    <XLoadingMessage title="Loading Users" msg="Fetching User Info" />
  );
  const whenUserEmpty = <h1>No Users Loaded</h1>;

  return (
    <React.Fragment>
      <Accordion>
        <XBoundary boundary={whenUserLoading} show={!userLoading}>
          <XBoundary boundary={whenUserEmpty} show={users && users.length > 0}>
            {users && showUsers()}
          </XBoundary>
        </XBoundary>

        <XBoundary boundary={whenSVCLoading} show={!svcLoading}>
          <XBoundary
            boundary={whenSVCEmpty}
            show={services && services.length > 0}
          >
            {services && showServices()}
          </XBoundary>
        </XBoundary>
      </Accordion>
    </React.Fragment>
  );
};

export default XAdminView;
