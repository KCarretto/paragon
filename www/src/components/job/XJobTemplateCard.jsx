import PropTypes from 'prop-types';
import React from 'react';
import { Button, Card, Icon, Menu } from 'semantic-ui-react';

const XJobTemplateCard = ({ id, name, tags }) => (
    <Card fluid >
        <Card.Content>
            <Menu text floated='right'>
                <Button compact basic circular>
                    <Icon.Group size='big'>
                        <Icon fitted name='cube' color='blue' />
                        <Icon corner='bottom left' name='chevron right' color='teal' />
                    </Icon.Group>
                </Button>
            </Menu>
            <Card.Header href={'/jobs/' + id}>{name}</Card.Header>


        </Card.Content>
        <Card.Content extra>
            <Icon name='tags' /> {tags ? tags.map(tag => tag.name).join(', ') : 'None'}
        </Card.Content>
    </Card>
)

XJobTemplateCard.propTypes = {
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    tags: PropTypes.arrayOf(PropTypes.shape({
        id: PropTypes.number.isRequired,
        name: PropTypes.string.isRequired
    })),
}

export default XJobTemplateCard