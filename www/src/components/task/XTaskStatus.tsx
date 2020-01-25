import moment from "moment";
import { IconProps } from "semantic-ui-react";
import { Task } from "../../graphql/models";

type TaskState = {
  text: string;
  icon: IconProps;
};

interface TaskStatus {
  QUEUED: TaskState;
  CLAIMED: TaskState;
  IN_PROGRESS: TaskState;
  COMPLETED: TaskState;
  STALE: TaskState;
  TIMED_OUT: TaskState;
  ERRORED: TaskState;
  getStatus(t: Task): TaskState;
  getTimestamp(t: Task): string;
}

class XTaskStatus implements TaskStatus {
  // Task is waiting to be retrieved
  QUEUED: {
    text: "In the queue, waiting to be claimed for execution.";
    icon: { name: "wait"; color: "violet" };
  };
  // Task has been sent to a target
  CLAIMED: {
    text: "Sent to target, waiting for execution to begin.";
    icon: { name: "send"; color: "blue" };
  };
  // Task has begun execution on a target
  IN_PROGRESS: {
    text: "Execution is currently in-progress.";
    icon: {
      name: "circle notched";
      color: "green";
      className: "XCircleIcon";
      bordered: false;
      circular: true;
      loading: true;
    };
  };
  // Task has successfully completed execution
  COMPLETED: {
    text: "Execution has successfully completed.";
    icon: { name: "check circle"; color: "green" };
  };
  // Task has been queued for longer than expected
  STALE: {
    text: "Queue time is taking longer than expected.";
    icon: { name: "wait"; color: "yellow" };
  };
  // Task has been claimed, but hasn't started execution within a reasonable amount of time
  TIMED_OUT: {
    text: "Execution is taking longer than expected to complete.";
    icon: {
      name: "circle notched";
      color: "red";
      className: "XCircleIcon";
      bordered: false;
      circular: true;
      loading: true;
    };
  };
  // Task encountered fatal error during execution
  ERRORED: {
    text: "Encountered an error during execution.";
    icon: { name: "times circle"; color: "red" };
  };

  public getStatus(t: Task): TaskState {
    if (t.error) {
      return this.ERRORED;
    } else if (t.execStopTime) {
      return this.COMPLETED;
    } else if (
      !t.claimTime &&
      moment(t.queueTime).isBefore(moment().subtract(5, "minutes"))
    ) {
      return this.STALE;
    } else if (
      t.claimTime &&
      moment(t.claimTime).isBefore(moment().subtract(5, "minutes"))
    ) {
      return this.TIMED_OUT;
    } else if (t.execStartTime) {
      return this.IN_PROGRESS;
    } else if (t.claimTime) {
      return this.CLAIMED;
    }
    return this.QUEUED;
  }
  public getTimestamp(t: Task): string {
    if (t.execStopTime) {
      return t.execStopTime;
    } else if (t.execStartTime) {
      return t.execStartTime;
    } else if (t.claimTime) {
      return t.claimTime;
    }
    return t.queueTime;
  }
}

export default XTaskStatus;
