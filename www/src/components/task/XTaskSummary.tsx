import moment from "moment";
import React from "react";
import { Link } from "react-router-dom";
import { Divider, Feed, Header, Icon, List } from "semantic-ui-react";
import { XTaskStatus } from ".";
import { Task } from "../../graphql/models";

type TaskSummaryParams = {
  tasks?: Task[];
  limit?: number;
};

const XTaskSummary = ({ tasks = [], limit = 3 }: TaskSummaryParams) => {
  const unshown = tasks.length - limit;

  return (
    <Feed>
      <Header sub>Recent Tasks</Header>
      {tasks.length > 0 ? (
        tasks
          .sort(
            (a, b) =>
              (XTaskStatus.getTimestamp(a) || 0) -
              (XTaskStatus.getTimestamp(b) || 0)
          )
          .slice(0, limit)
          .map((task, index) => (
            <Feed.Event key={index}>
              <Feed.Label>
                <Icon fitted size="big" {...XTaskStatus.getStatus(task).icon} />
              </Feed.Label>
              <Feed.Content>
                <Feed.Date>
                  {moment(XTaskStatus.getTimestamp(task)).fromNow()}
                </Feed.Date>
                <Feed.Summary>
                  <Link to={"/tasks/" + task.id}>
                    <List.Header>{task.job.name}</List.Header>
                  </Link>
                </Feed.Summary>
                <Divider />
              </Feed.Content>
            </Feed.Event>
          ))
      ) : (
        <Header sub disabled>
          No recent tasks
        </Header>
      )}
      {unshown > 0 ? <span>and {unshown} more...</span> : ""}
    </Feed>
  );
};

export default XTaskSummary;
