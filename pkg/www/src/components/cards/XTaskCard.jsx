import PropTypes from 'prop-types';
import React from 'react';
import { Button, Card, Icon, Progress } from 'semantic-ui-react';

const XTaskCard = ({ name, tags }) => (
    <Card>
        <Card.Content>
            <Card.Header>
                <Icon name='linkify' size='small' />
                {name ? name : 'Untitled Task'}
            </Card.Header>
            <Card.Meta>
                <Icon name='tags' /> {tags ? tags.join(', ') : 'None'}
            </Card.Meta>
            <Card.Description>
                <Button icon labelPosition='right'>Subscribe<Icon name='bell' /></Button>
                <Button icon labelPosition='right'>View<Icon name='external' /></Button>
            </Card.Description>
        </Card.Content>
        <Card.Content extra>
            <Progress color='red' size='small' percent={50} active>In Progress</Progress>
        </Card.Content>
    </Card>
)

XTaskCard.propTypes = {
    name: PropTypes.string,
    tags: PropTypes.arrayOf(PropTypes.string),
}

export default XTaskCard