import moment from "moment";
import { IconProps } from "semantic-ui-react";
import { Task } from "../../graphql/models";

type TaskState = {
  text: string;
  icon: IconProps;
};

interface TaskStatus {
  getStatus(t: Task): TaskState;
  getTimestamp(t: Task): string;
}

class XTaskStatus implements TaskStatus {
  // Task is waiting to be retrieved
  static QUEUED: TaskState = {
    text: "In the queue, waiting to be claimed for execution.",
    icon: { name: "wait", color: "violet" }
  };
  // Task has been sent to a target
  static CLAIMED: TaskState = {
    text: "Sent to target, waiting for execution to begin.",
    icon: { name: "send", color: "blue" }
  };
  // Task has begun execution on a target
  static IN_PROGRESS: TaskState = {
    text: "Execution is currently in-progress.",
    icon: {
      name: "circle notched",
      color: "green",
      className: "XCircleIcon",
      bordered: false,
      circular: true,
      loading: true
    }
  };
  // Task has successfully completed execution
  static COMPLETED: TaskState = {
    text: "Execution has successfully completed.",
    icon: { name: "check circle", color: "green" }
  };
  // Task has been queued for longer than expected
  static STALE: TaskState = {
    text: "Queue time is taking longer than expected.",
    icon: { name: "wait", color: "yellow" }
  };
  // Task has been claimed, but hasn't started execution within a reasonable amount of time
  static TIMED_OUT: TaskState = {
    text: "Execution is taking longer than expected to complete.",
    icon: {
      name: "circle notched",
      color: "red",
      className: "XCircleIcon",
      bordered: false,
      circular: true,
      loading: true
    }
  };
  // Task encountered fatal error during execution
  static ERRORED: TaskState = {
    text: "Encountered an error during execution.",
    icon: { name: "times circle", color: "red" }
  };

  public getStatus(t: Task): TaskState {
    if (t.error) {
      return XTaskStatus.ERRORED;
    } else if (t.execStopTime) {
      return XTaskStatus.COMPLETED;
    } else if (
      t.claimTime &&
      moment(t.queueTime).isBefore(moment().subtract(5, "minutes"))
    ) {
      return XTaskStatus.STALE;
    } else if (
      t.claimTime &&
      moment(t.claimTime).isBefore(moment().subtract(5, "minutes"))
    ) {
      return XTaskStatus.TIMED_OUT;
    } else if (t.execStartTime) {
      return XTaskStatus.IN_PROGRESS;
    } else if (t.claimTime) {
      return XTaskStatus.CLAIMED;
    }
    return XTaskStatus.QUEUED;
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
