import moment from 'moment';

const XTaskStatus = {
    // Task is waiting to be retrieved
    QUEUED: {
        text: 'In the queue, waiting to be claimed for execution.',
        icon: { name: 'wait', color: 'violet' }
    },
    // Task has been sent to a target
    CLAIMED: {
        text: 'Sent to target, waiting for execution to begin.',
        icon: { name: 'send', color: 'blue' }
    },
    // Task has begun execution on a target
    IN_PROGRESS: { text: 'Execution is currently in-progress.', icon: { name: 'circle notched', color: 'green', className: 'XCircleIcon', bordered: false, circular: true, loading: true } },
    // Task has successfully completed execution
    COMPLETED: { text: 'Execution has successfully completed.', icon: { name: 'check circle', color: 'green', className: 'XCircleIcon', bordered: false, circular: true } },
    // Task has been queued for longer than expected
    STALE: { text: 'Queue time is taking longer than expected.', icon: { name: 'wait', color: 'yellow', className: 'XCircleIcon', bordered: false, circular: true } },
    // Task has been claimed, but hasn't started execution within a reasonable amount of time
    TIMED_OUT: { text: 'Execution is taking longer than expected to complete.', icon: { name: 'circle notched', color: 'red', className: 'XCircleIcon', bordered: false, circular: true, loading: true } },
    // Task encountered fatal error during execution
    ERRORED: {
        text: 'Encountered an error during execution.', icon: { name: 'times circle', color: 'red', className: 'XCircleIcon', bordered: false, circular: true }
    },
}

XTaskStatus.getStatus = ({ queueTime, claimTime, execStartTime, execStopTime, error }) => {
    if (error) {
        return XTaskStatus.ERRORED;
    } else if (execStopTime) {
        return XTaskStatus.COMPLETED;
    }
    else if (!claimTime && moment(queueTime).isBefore(moment().subtract(5, 'minutes'))) {
        return XTaskStatus.STALE;
    } else if (claimTime && moment(claimTime).isBefore(moment().subtract(5, 'minutes'))) {
        return XTaskStatus.TIMED_OUT;
    } else if (execStartTime) {
        return XTaskStatus.IN_PROGRESS;
    } else if (claimTime) {
        return XTaskStatus.CLAIMED;
    }
    return XTaskStatus.QUEUED;
}

XTaskStatus.getTimestamp = ({ queueTime, claimTime, execStartTime, execStopTime, error }) => {
    if (execStopTime) {
        return execStopTime;
    } else if (execStartTime) {
        return execStartTime;
    } else if (claimTime) {
        return claimTime;
    }
    return queueTime;
}

export default XTaskStatus;
