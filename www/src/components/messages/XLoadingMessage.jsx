import React from 'react';
import { Icon, Message } from 'semantic-ui-react';

const XLoadingMessage = ({ title, msg, hidden }) => (
    <Message icon size='massive' hidden={hidden}>
        <Icon name='circle notched' loading />
        <Message.Content>
            <Message.Header>{title}</Message.Header>
            {msg}
        </Message.Content>
    </Message>
);

export default XLoadingMessage;