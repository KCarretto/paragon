import moment from "moment";
import * as React from "react";
import { FunctionComponent } from "react";
import { Link } from "react-router-dom";
import { Divider, Feed, Icon } from "semantic-ui-react";
import { Event } from "../../graphql/models";

export enum EventKind {
  CREATE_JOB = "CREATE_JOB",
  COMPLETE_JOB = "COMPLETE_JOB",
  ADD_CREDENTIAL_FOR_TARGET = "ADD_CREDENTIAL_FOR_TARGET",
  ACTIVATE_SERVICE = "ACTIVATE_SERVICE",
  ACTIVATE_USER = "ACTIVATE_USER",
  UPLOAD_FILE = "UPLOAD_FILE",
  CREATE_LINK = "CREATE_LINK",
  CREATE_USER = "CREATE_USER",
  CREATE_SERVICE = "CREATE_SERVICE",
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
  if (event.kind === EventKind.ACTIVATE_USER) {
    return {
      id: event!.user!.id || "0",
      name: event!.user!.name || "anonymous walrus",
      imgURL: event!.user.photoURL || "/app/default_profile.gif",
      isUser: true,
    };
  }

  if (event.owner !== null) {
    return {
      id: event!.owner!.id || "0",
      name: event!.owner!.name || "anonymous hippo",
      imgURL: event.owner.photoURL || "/app/default_profile.gif",
      isUser: true,
    };
  }

  return {
    id: event.svcOwner!.id,
    name: event.svcOwner!.name,
    imgURL: "/app/default_profile.gif",
    isUser: false,
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
  actor,
}) => {
  switch (kind) {
    case EventKind.CREATE_JOB:
      return (
        <span>
          {" "}
          created job{" "}
          <Link to={"/jobs/" + event.job ? event.job.id : ""}>
            {event.job ? event.job.name : "[DELETED]"}
          </Link>
        </span>
      );
    case EventKind.COMPLETE_JOB:
      return (
        <span>
          {" "}
          completed job{" "}
          <Link to={"/jobs/" + event.job ? event.job.id : ""}>
            {event.job ? event.job.name : "[DELETED]"}
          </Link>
        </span>
      );
    case EventKind.ADD_CREDENTIAL_FOR_TARGET:
      return (
        <span>
          {" "}
          added credentials for{" "}
          <Link to={"/targets/" + event.target ? event.target.id : ""}>
            {event.credential ? event.credential.principal : "[DELETED]"}:
            {event.credential ? event.credential.secret : "[DELETED]"}@
            {event.target ? event.target.name : "[DELETED]"}{" "}
          </Link>
        </span>
      );
    case EventKind.ACTIVATE_SERVICE:
      return (
        <span>
          {" "}
          activated service{" "}
          <Link to={"/admin"}>
            {event.service ? event.service.name : "[DELETED]"}
          </Link>
        </span>
      );
    case EventKind.UPLOAD_FILE:
      return (
        <span>
          {" "}
          uploaded file{" "}
          <Link to={"/files"}>
            {event.file ? event.file.name : "[DELETED]"}
          </Link>
        </span>
      );
    case EventKind.CREATE_LINK:
      return (
        <span>
          {" "}
          created link{" "}
          <Link to={"/files"}>
            {event.link ? event.link.alias : "[DELETED]"}
          </Link>
        </span>
      );
    case EventKind.CREATE_USER:
      return <span> requested to join</span>;
    case EventKind.ACTIVATE_USER:
      return <span> joined!</span>;
    case EventKind.CREATE_SERVICE:
      return <span> service requested registration</span>;
    default:
      return <span> Caused an unhandled event to occur! ({kind})</span>;
  }
  return <span />;
};

const XEventDetails: FunctionComponent<EventProps> = ({
  kind,
  event,
  actor,
}) => <span />;

const XEventSummary: FunctionComponent<EventProps> = ({
  kind,
  event,
  actor,
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

const XEvent: FunctionComponent<XEventProps> = ({ event }) => {
  let actor = GetEventActor(event);
  let kind: EventKind = EventKind[event.kind];
  if (!kind) {
    return <span />;
  }

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
            <Icon name="like" />
            420 Likes
          </Feed.Like>
        </Feed.Meta>
        <Divider />
      </Feed.Content>
    </Feed.Event>
  );
};

export default XEvent;
