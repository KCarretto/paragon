import PropTypes from 'prop-types';
import React from 'react';
import { Card, Icon, Menu } from 'semantic-ui-react';
import { XTaskSummary } from '../task';
import XJobQueueModal from './XJobQueueModal';

const XJobCard = ({ id, name, tasks, tags }) => (
    <Card fluid >
        <Card.Content>
            <Menu text floated='right'>
                <XJobQueueModal />
            </Menu>
            {/* <Menu text floated='right'>
                <Button compact basic circular>
                    <Icon.Group size='big'>
                        <Icon fitted name='cube' color='blue' />
                        <Icon corner='bottom left' name='chevron right' color='teal' />
                    </Icon.Group>
                </Button>
            </Menu> */}
            <Card.Header href={'/jobs/' + id}>{name}</Card.Header>

            <XTaskSummary tasks={tasks} />

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