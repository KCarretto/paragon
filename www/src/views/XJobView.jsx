import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Container, Header, Icon } from 'semantic-ui-react';
import { XJobHeader } from '../components/job';
import { XErrorMessage, XLoadingMessage } from '../components/messages';
import { XTaskCard, XTaskContent } from '../components/task';

export const JOB_QUERY = gql`
query Job($id: ID!) {
    job(id: $id) {
        id
        name
        content
        tags {
            id
            name
        }
        tasks {
            id
            queueTime
            claimTime
            execStartTime
            execStopTime
            content
            output
            error
            sessionID

            target {
                id
              	name
            }
            job {
                id
                name
                tags {
                    id
                    name
                }
            }
        }
    }
}`

const XJobView = () => {
    let { id } = useParams();
    const [error, setError] = useState(null);

    const [name, setName] = useState('');
    const [content, setContent] = useState('');
    const [tags, setTags] = useState([]);
    const [tasks, setTasks] = useState([]);

    const { called, loading } = useQuery(JOB_QUERY, {
        variables: { id },
        pollInterval: 5000,
        onCompleted: data => {
            setError(null);
            if (!data || !data.job) {
                data = { job: { name: '', content: '', tags: [], tasks: [] } }
            }

            setName(data.job.name || '');
            setContent(data.job.content || '')
            setTags(data.job.tags || []);
            setTasks(data.job.tasks || []);
        },
        onError: err => setError(err),
    });

    const showCards = () => {
        if (!tasks || tasks.length < 1) {
            return (
                // TODO: Better styling
                <h1>No tasks found!</h1>
            );
        }
        return (
            <Card.Group centered itemsPerRow={4}>
                {tasks.map(task => (<XTaskCard key={task.id} task={task} />))}
            </Card.Group>
        );
    };

    return (
        <Container fluid style={{ padding: '20px' }}>
            <XJobHeader name={name} tags={tags} />

            <XErrorMessage title='Error Loading Job' err={error} />
            <XLoadingMessage
                title='Loading Job'
                msg='Fetching job information...'
                hidden={called && !loading}
            />

            <XTaskContent content={content} />

            <Header size='large' block inverted>
                <Icon name='tasks' />
                <Header.Content>Tasks</Header.Content>
            </Header>

            {showCards()}
        </Container>
    );
}


export default XJobView