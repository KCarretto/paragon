import * as React from "react";
import { Icon, Message } from "semantic-ui-react";

type LoadingMessageParams = {
  title: string;
  msg: string;
  hidden?: boolean;
};

const XLoadingMessage = ({ title, msg, hidden }: LoadingMessageParams) => (
  <Message icon size="massive" hidden={hidden}>
    <Icon name="circle notched" loading />
    <Message.Content>
      <Message.Header>{title}</Message.Header>
      {msg}
    </Message.Content>
  </Message>
);

export default XLoadingMessage;
