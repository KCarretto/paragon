import * as React from "react";
import { Icon } from "semantic-ui-react";
import { Tag } from "../../graphql/models";

type TagParams = {
  tags: Tag[];
  defaultText?: string;
};

export default ({ tags = [], defaultText }: TagParams) => {
  if (!tags || tags.length < 1) {
    if (!defaultText) {
      return <span />;
    }
    return (
      <span>
        <Icon name="tags" />
        {defaultText}
      </span>
    );
  }
  return (
    <span>
      <Icon name="tags" /> {tags.map(tag => tag.name).join(", ")}
    </span>
  );
};
