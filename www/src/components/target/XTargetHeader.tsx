import moment from "moment";
import * as React from "react";
import { Header, Icon, IconProps } from "semantic-ui-react";
import { Tag } from "../../graphql/models";
import { XTags } from "../tag";

type TargetHeaderParams = {
  name: string;
  tags: Tag[];
  lastSeen?: string;
  icon?: React.CElement<IconProps, Icon>;
};

export default ({ name, lastSeen, tags, icon }: TargetHeaderParams) => (
  <Header size="huge">
    {icon ? icon : <Icon name="desktop" />}
    {!lastSeen || moment(lastSeen).isBefore(moment().subtract(5, "minutes")) ? (
      <Icon name="times circle" color="red" size="large" />
    ) : (
      <Icon name="check circle" color="green" size="large" />
    )}
    <Header.Content>{name}</Header.Content>
    <Header.Subheader>
      <XTags tags={tags} />
    </Header.Subheader>
  </Header>
);
