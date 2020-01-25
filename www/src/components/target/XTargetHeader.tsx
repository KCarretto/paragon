import * as React from "react";
import { Header, Icon } from "semantic-ui-react";
import { Tag } from "../../graphql/models";
import { XTags } from "../tag";

type TargetHeaderParams = {
  name: String;
  tags: Tag[];
  icon?: Icon;
};

export default ({ name, tags, icon }: TargetHeaderParams) => (
  <Header size="huge">
    {icon ? icon : <Icon name="desktop" />}
    <Header.Content>{name}</Header.Content>
    <Header.Subheader>
      <XTags tags={tags} />
    </Header.Subheader>
  </Header>
);
