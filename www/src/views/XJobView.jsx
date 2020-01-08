import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import { default as React } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Container, Icon, Label, Loader } from 'semantic-ui-react';
import { XTaskStatus } from '../components/task';


export const JOB_QUERY = gql`
query Job($id: ID!) {
    job(id: $id) {
        id

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

    const { loading, error, data } = useQuery(JOB_QUERY, {
        variables: { id },
    });

    if (loading) return (<Loader active />);
    if (error) return (`${error}`);

    const taskCards = data.job.tasks.map(task => {
        let status = XTaskStatus.getStatus(task).icon;
        return <Card fluid centered>
            <Card.Content>
                <Card.Header>
                    <Label size='mini' icon={<Icon fitted size='big' {...status} />} attached='top right' />
                    {task.job.name ? task.job.name : 'Untitled Job'}
                </Card.Header>
                {
                    task.job.tags && task.job.tags.length != 0 ?
                        <Card.Meta>
                            <Icon name='tags' /> {task.job.tags.map(tag => tag.name).join(', ')}
                        </Card.Meta>
                        :
                        <Card.Meta>
                            <Icon name='tags' /> None
                    </Card.Meta>
                }
                <Card.Description>
                    <div className="ui visible message">
                        <div className="header">
                            Content
                    </div>
                        <pre>{task.content}</pre>
                    </div>
                    {task.output ?
                        <div className="ui positive message">
                            <div className="header">
                                Output
                        </div>
                            <pre>{task.output}</pre>
                        </div> :
                        <div></div>
                    }
                    {task.error ?
                        <div className="ui negative message">
                            <div className="header">
                                Error
                        </div>
                            <pre>{task.error}</pre>
                        </div>
                        :
                        <div></div>
                    }
                </Card.Description>
            </Card.Content>
            {task.sessionID ?
                <Card.Content extra>
                    <a>
                        <i aria-hidden="true" className="user icon"></i>
                        {task.sessionID}
                    </a>
                </Card.Content>
                :
                <div></div>
            }

        </Card>
    })

    return (
        <Container fluid style={{ padding: '20px' }}>
            {taskCards}
        </Container>
    );
}


export default XJobView