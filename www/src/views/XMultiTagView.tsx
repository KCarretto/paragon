import { useQuery } from "@apollo/react-hooks";
import gql from "graphql-tag";
import * as React from "react";
import { Container, Label, Loader, Menu } from "semantic-ui-react";
import XClipboard from "../components/form/XClipboard";
import { XErrorMessage } from "../components/messages";
import { XNoTagsFound, XTagCreateModal } from "../components/tag";
import { Tag } from "../graphql/models";

export const MULTI_TAG_QUERY = gql`
  {
    tags {
      id
      name
    }
  }
`;

export type MultiTagResponse = {
  tags: Tag[];
};

const XMultiTagView = () => {
  const { called, loading, error, data } = useQuery<MultiTagResponse>(
    MULTI_TAG_QUERY
  );

  const showList = () => {
    if (!data || !data.tags || data.tags.length < 1) {
      return <XNoTagsFound />;
    }
    return (
      <Label.Group tag size="big">
        {data.tags.map((tag, index) => (
          <XClipboard value={tag.name}>
            <Label key={index}>{tag.name}</Label>
          </XClipboard>
        ))}
      </Label.Group>
    );
  };

  return (
    <Container style={{ padding: "10px" }}>
      <Menu secondary>
        <Menu.Item position="right">
          <XTagCreateModal />
        </Menu.Item>
      </Menu>
      <Container fluid style={{ padding: "20px" }}>
        <Loader disabled={!called || !loading} />
        {showList()}
        <XErrorMessage title="Error Loading Tags" err={error} />
      </Container>
    </Container>
  );
};

export default XMultiTagView;
