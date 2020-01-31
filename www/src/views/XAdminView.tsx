import { useMutation, useQuery } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
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
  Loader,
  Radio
} from "semantic-ui-react";
import XClipboard from "../components/form/XClipboard";
import { XErrorMessage } from "../components/messages";
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
  const [servicesError, setServicesError] = useState<ApolloError>(null);
  const [usersError, setUsersError] = useState<ApolloError>(null);

  const usersQuery = useQuery<UsersResponse>(ADMIN_USERS_QUERY, {
    onError: err => setUsersError(err)
  });
  const servicesQuery = useQuery<ServicesResponse>(ADMIN_SERVICES_QUERY, {
    onError: err => setServicesError(err)
  });

  const [activateUser] = useMutation(ACTIVATE_USER_MUTATION, {
    refetchQueries: [{ query: ADMIN_USERS_QUERY }]
  });
  const [deactivateUser] = useMutation(DEACTIVATE_USER_MUTATION, {
    refetchQueries: [{ query: ADMIN_USERS_QUERY }]
  });
  const [activateService] = useMutation(ACTIVATE_SERVICE_MUTATION, {
    refetchQueries: [{ query: ADMIN_SERVICES_QUERY }]
  });
  const [deactivateService] = useMutation(DEACTIVATE_SERVICE_MUTATION, {
    refetchQueries: [{ query: ADMIN_SERVICES_QUERY }]
  });
  const [makeAdmin] = useMutation(MAKE_ADMIN_MUTATION, {
    refetchQueries: [{ query: ADMIN_USERS_QUERY }]
  });
  const [removeAdmin] = useMutation(REMOVE_ADMIN_MUTATION, {
    refetchQueries: [{ query: ADMIN_USERS_QUERY }]
  });

  const handleClick = (e, titleProps) => {
    const { index } = titleProps;
    if (index === 0) {
      setOpenUsers(!openUsers);
    } else {
      setOpenServices(!openServices);
    }
  };

  const handleUserActivate = (e, p) => {
    if (p.checked) {
      activateUser({
        variables: { id: p.value }
      }).then(({ data, errors }) => {
        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          setUsersError(e);
          return;
        }
      });
    } else {
      deactivateUser({
        variables: { id: p.value }
      }).then(({ data, errors }) => {
        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          setUsersError(e);
          return;
        }
      });
    }
  };

  const handleUserAdmin = (e, p) => {
    if (p.checked) {
      makeAdmin({ variables: { id: p.value } }).then(({ data, errors }) => {
        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          setUsersError(e);
          return;
        }
      });
    } else {
      removeAdmin({ variables: { id: p.value } }).then(({ data, errors }) => {
        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          setUsersError(e);
          return;
        }
      });
    }
  };

  const handleServiceActivate = (e, p) => {
    if (p.checked) {
      activateService({
        variables: { id: p.value }
      }).then(({ data, errors }) => {
        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          setServicesError(e);
          return;
        }
      });
    } else {
      deactivateService({
        variables: { id: p.value }
      }).then(({ data, errors }) => {
        if (errors && errors.length > 0) {
          let s = errors.map(e => e.message);
          let e = new ApolloError({
            graphQLErrors: errors,
            errorMessage: s.join("\n")
          });
          setServicesError(e);
          return;
        }
      });
    }
  };

  const showUsers = () => {
    if (
      !usersQuery.data ||
      !usersQuery.data.users ||
      usersQuery.data.users.length < 1
    ) {
      return <span />;
    }
    return (
      <React.Fragment>
        <Loader disabled={!usersQuery.called || !usersQuery.loading} />
        <Accordion.Title
          inverted
          active={openUsers}
          index={0}
          onClick={handleClick}
        >
          <Header icon as="h2" textAlign="center">
            <Icon name="dropdown" />
            <Header.Content>Users</Header.Content>
          </Header>
        </Accordion.Title>
        <Accordion.Content active={openUsers}>
          <Card.Group centered itemsPerRow={4}>
            {usersQuery.data.users.map(u => {
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
                      />
                    </Card.Description>
                  </Card.Content>
                </Card>
              );
            })}
          </Card.Group>
        </Accordion.Content>
        <XErrorMessage title="Error Loading Users" err={usersError} />
        <Divider />
      </React.Fragment>
    );
  };

  const showServices = () => {
    if (
      !servicesQuery.data ||
      !servicesQuery.data.services ||
      servicesQuery.data.services.length < 1
    ) {
      return <span />;
    }
    return (
      <React.Fragment>
        <Loader disabled={!servicesQuery.called || !servicesQuery.loading} />
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
          <Card.Group centered itemsPerRow={4}>
            {servicesQuery.data.services.map(s => {
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
          </Card.Group>
        </Accordion.Content>
        <XErrorMessage title="Error Loading Users" err={servicesError} />
        <Divider />
      </React.Fragment>
    );
  };

  return (
    <Accordion>
      {showUsers()}
      {showServices()}
    </Accordion>
  );
};

export default XAdminView;
