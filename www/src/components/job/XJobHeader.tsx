import * as React from "react";
import { Header, Icon, IconProps } from "semantic-ui-react";
import { Tag } from "../../graphql/models";
import { XTags } from "../tag";

type JobHeaderParams = {
  name: string;
  tags: Tag[];
  icon?: React.CElement<IconProps, Icon>;
};

export default ({ name, tags, icon }: JobHeaderParams) => (
  <Header size="huge">
    {icon !== null ? icon : React.createElement(Icon, { name: "cube" })}
    <Header.Content>{name}</Header.Content>
    <Header.Subheader>
      <XTags tags={tags} />
    </Header.Subheader>
  </Header>
);
