import * as React from "react";
import { Header, Icon, IconProps } from "semantic-ui-react";
import { Tag } from "../../graphql/models";
import { XTags } from "../tag";

type TargetHeaderParams = {
  name: string;
  tags: Tag[];
  icon?: React.CElement<IconProps, Icon>;
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
