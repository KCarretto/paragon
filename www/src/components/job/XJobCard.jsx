import PropTypes from 'prop-types';
import React from 'react';
import { Card, Icon } from 'semantic-ui-react';

const XJobCard = ({ id, name, tags }) => (
    <Card fluid >
        <Card.Content>
            <Card.Header href={'/jobs/' + id}>{name} </Card.Header>
        </Card.Content>
        <Card.Content extra>
            <Icon name='tags' /> {tags ? tags.map(tag => tag.name).join(', ') : 'None'}
        </Card.Content>
    </Card>
)

XJobCard.propTypes = {
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    tags: PropTypes.arrayOf(PropTypes.shape({
        id: PropTypes.number.isRequired,
        name: PropTypes.string.isRequired
    })),
}

export default XJobCard