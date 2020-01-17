import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Button, Container, Icon, Label } from 'semantic-ui-react';
import { XJobHeader } from '../components/job';
import { XErrorMessage, XLoadingMessage } from '../components/messages';
import { XTaskContent, XTaskError, XTaskOutput, XTaskStatus } from '../components/task';

const TASK_QUERY = gql`
    query Task($id: ID!) {
    task(id: $id) {
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
`;

const XTaskView = () => {
    let { id } = useParams();
    const [loadingError, setLoadingError] = useState(null);

    const [queueTime, setQueueTime] = useState(null);
    const [claimTime, setClaimTime] = useState(null);
    const [execStartTime, setExecStartTime] = useState(null);
    const [execStopTime, setExecStopTime] = useState(null);
    const [content, setContent] = useState(null);
    const [output, setOutput] = useState(null);
    const [error, setError] = useState(null);
    const [sessionID, setSessionID] = useState(null);
    const [target, setTarget] = useState({});
    const [jobID, setJobID] = useState(null);
    const [name, setName] = useState('');
    const [tags, setTags] = useState([]);


    const { called, loading } = useQuery(TASK_QUERY, {
        variables: { id },
        onCompleted: data => {
            setLoadingError(null);

            if (!data || !data.task) {
                data = {
                    task: {
                        queueTime: null,
                        claimTime: null,
                        execStartTime: null,
                        execStopTime: null,
                        content: null,
                        output: null,
                        error: null,
                        sessionID: null,
                    }
                }
            }
            if (!data.task.target) {
                data.task.target = {};
            }
            if (!data.task.job) {
                data.task.job = { id: null, name: '', tags: [] }
            }

            setQueueTime(data.task.queueTime);
            setClaimTime(data.task.claimTime);
            setExecStartTime(data.task.execStartTime);
            setExecStopTime(data.task.execStopTime);
            setContent(data.task.content);
            setOutput(data.task.output);
            setError(data.task.error);
            setSessionID(data.task.sessionID);
            setTarget(data.task.target);
            setJobID(data.task.job.id);
            setTags(data.task.job.tags);
            setName(data.task.job.name);
        },
        onError: err => setLoadingError(err),
    });

    let status = XTaskStatus.getStatus({
        queueTime: queueTime,
        claimTime: claimTime,
        execStartTime: execStartTime,
        execStopTime: execStopTime,
        error: error,
    }).icon;

    return (
        <Container fluid style={{ padding: '20px' }}>
            <a href={'/jobs/' + jobID}>
                <XJobHeader name={name} tags={tags} icon={<Icon size='large' {...status} />} />
            </a>
            {!target || !target.id ? <span /> :
                <Button
                    basic
                    animated
                    href={'/targets/' + target.id}
                    color='blue'
                    size='small'
                    style={{ margin: '15px' }}
                >
                    <Button.Content visible>{target.name || 'View Target'}</Button.Content>
                    <Button.Content hidden>
                        <Icon name='arrow right' />
                    </Button.Content>
                </Button>}
            {!sessionID ? <span /> : <Label>SessionID: {sessionID}</Label>}

            <XErrorMessage title='Error Loading Task' err={loadingError} />
            <XLoadingMessage
                title='Loading Task'
                msg='Fetching task information...'
                hidden={called && !loading}
            />

            <XTaskContent content={content} />
            <XTaskOutput output={output} />
            <XTaskError error={error} />

        </Container>
    );
}

export default XTaskView;
