import React from 'react';
import { Message } from 'semantic-ui-react';

export default ({ title, msg }) => (
    <Message negative>
        <Message.Header>{title}</Message.Header>
        {msg}
    </Message>
)

