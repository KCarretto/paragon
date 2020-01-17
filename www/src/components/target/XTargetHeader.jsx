import React from 'react';
import { Header, Icon } from 'semantic-ui-react';
import { XTags } from '../tag';

export default ({ name, tags, icon }) => (
    <Header size='huge'>
        {icon ? icon : <Icon name='desktop' />}
        <Header.Content>{name}</Header.Content>
        <Header.Subheader>
            <XTags tags={tags} />
        </Header.Subheader>
    </Header>
);