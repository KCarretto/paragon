import PropTypes from 'prop-types';
import React from 'react';
import { Card } from 'semantic-ui-react';
import { XTags } from '../tag';
import { XTaskSummary } from '../task';

const XJobCard = ({ id, name, tasks, tags }) => (
    <Card fluid >
        <Card.Content>
            <Card.Header href={'/jobs/' + id}>{name}</Card.Header>

            <XTaskSummary tasks={tasks} />

        </Card.Content>
        <Card.Content extra>
            <XTags defaultText='None' />
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