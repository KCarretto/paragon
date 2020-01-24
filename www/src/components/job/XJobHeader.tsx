import * as React from "react";
import { Header, Icon } from "semantic-ui-react";
import { Tag } from "../../graphql/models";
import { XTags } from "../tag";

type JobHeaderParams = {
  name: String;
  tags: Tag[];
  icon: String;
};

export default ({ name, tags, icon }: JobHeaderParams) => (
  <Header size="huge">
    {icon ? icon : <Icon name="cube" />}
    <Header.Content>{name}</Header.Content>
    <Header.Subheader>
      <XTags tags={tags} />
    </Header.Subheader>
  </Header>
);
