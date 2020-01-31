import moment from "moment";
import * as React from "react";
import { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import { Divider, Feed, Icon } from "semantic-ui-react";
import { Event } from "../../graphql/models";

export enum EventKind {
  JobCreated = "CREATE_JOB",
  JobCompleted = "COMPLETE_JOB",
  ServiceActivated = "ACTIVATE_SERVICE",
  UserActivated = "CREATE_USER",
  FileUploaded = "UPLOAD_FILE",
  LinkCreated = "CREATE_LINK"
}

interface EventProps extends XEventProps {
  actor: EventActor;
}

interface XEventProps {
  kind: EventKind;
  event: Event;
}

type EventActor = {
  id: string;
  name: string;
  imgURL: string;
  isUser: boolean;
};

const XEventNounJob = () => <span />;
const XEventNounService = () => <span />;
const XEventNounUser = () => <span />;
const XEventNounFile = () => <span />;
const XEventNounLink = () => <span />;

const GetEventActor: (event: Event) => EventActor = (event: Event) => {
  if (event.owner !== null) {
    return {
      id: event!.owner!.id || "0",
      name: event!.owner!.name || "anonymous hippo",
      imgURL: event.owner.photoURL || "/app/default_profile.gif",
      isUser: true
    };
  }

  return {
    id: event.svcOwner!.id,
    name: event.svcOwner!.name,
    imgURL: "/app/default_profile.gif",
    isUser: false
  };
};

type XEventSummaryProps = {
  kind: EventKind;
  event: Event;
  actor: EventActor;
};

const XEventDescription: FunctionComponent<EventProps> = ({
  kind,
  event,
  actor
}) => {
  switch (kind) {
    case EventKind.JobCreated:
    case EventKind.JobCompleted:
      return (
        <span>
          {" "}
          completed job{" "}
          <Link to={"/jobs/" + event.job.id}>{event.job.name}</Link>
        </span>
      );
    case EventKind.ServiceActivated:
      return <XEventNounService />;
    case EventKind.FileUploaded:
      return <XEventNounJob />;
    case EventKind.LinkCreated:
      return <XEventNounLink />;
    default:
      return <span> Caused an invalid event to occur! ({kind})</span>;
  }
  return <span />;
};

const XEventDetails: FunctionComponent<EventProps> = ({
  kind,
  event,
  actor
}) => <span />;

const XEventSummary: FunctionComponent<EventProps> = ({
  kind,
  event,
  actor
}) => {
  return (
    <Feed.Summary>
      <Feed.User>{actor.name}</Feed.User>
      <XEventDescription kind={kind} event={event} actor={actor} />
      <Feed.Date>{moment(event.creationTime).fromNow()}</Feed.Date>
    </Feed.Summary>
  );
};

// return <Feed.Label image='/images/avatar/small/elliot.jpg' />

const XEvent: FunctionComponent<XEventProps> = ({ kind, event }) => {
  let actor = GetEventActor(event);

  return (
    <Feed.Event>
      <Feed.Label>
        <img src={actor.imgURL} />
      </Feed.Label>
      <Feed.Content>
        <XEventSummary kind={kind} event={event} actor={actor} />
        <XEventDetails kind={kind} event={event} actor={actor} />

        <Feed.Meta>
          <Feed.Like>
            <Icon name="like" />4 Likes
          </Feed.Like>
        </Feed.Meta>
        <Divider />
      </Feed.Content>
    </Feed.Event>
  );
};

export default XEvent;
