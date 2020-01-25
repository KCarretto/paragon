import moment from "moment";
import * as React from "react";
import {
  Button,
  Card,
  Feed,
  Header,
  Icon,
  List,
  SemanticCOLORS
} from "semantic-ui-react";
import { File } from "../../graphql/models";

const XFileCard = (f: File) => {
  let colors: SemanticCOLORS[] = [
    "olive",
    "green",
    "teal",
    "blue",
    "violet",
    "purple",
    "pink"
  ];

  return (
    <Card fluid>
      <Card.Content>
        <Button.Group floated="right">
          <Button basic color="blue" icon="linkify" /> {/* TODO: Implement */}
          <Button basic color="blue" icon="cloud upload" />{" "}
          {/* TODO: Implement */}
          <Button
            basic
            color="blue"
            icon="cloud download"
            href={"/cdn/download/" + f.name}
          />
        </Button.Group>
        <Card.Header>{f.name}</Card.Header>
        <Card.Meta>{f.size} bytes</Card.Meta>
        <Card.Description>
          <Header
            sub
            disabled={!f.links || f.links.length < 1}
            style={{ marginTop: "5px" }}
          >
            <Header.Content>
              {f.links && f.links.length > 0
                ? "Links (" + f.links.length + " total)"
                : "No Active Links"}
            </Header.Content>
          </Header>
          <Feed style={{ maxHeight: "25vh", overflowY: "auto" }}>
            {f.links && f.links.length > 0 ? (
              f.links.map((link, index) => (
                <Feed.Event key={index}>
                  <Feed.Label>
                    <Icon
                      fitted
                      name="linkify"
                      color={colors[Math.floor(Math.random() * colors.length)]}
                    />
                  </Feed.Label>
                  <Feed.Content>
                    <Feed.Summary>
                      <List.Header>{link.alias}</List.Header>
                      <Feed.Date>
                        {link.expirationTime
                          ? "Expires in " + moment().to(link.expirationTime)
                          : "Never expires"}
                      </Feed.Date>
                    </Feed.Summary>
                    <Feed.Meta>
                      {link.clicks && link.clicks >= 0
                        ? link.clicks + " Clicks left"
                        : "Unlimited clicks"}
                    </Feed.Meta>
                  </Feed.Content>
                </Feed.Event>
              ))
            ) : (
              <span />
            )}
          </Feed>
        </Card.Description>
      </Card.Content>
      <Card.Content extra>
        Created: {moment(f.creationTime).fromNow()}
        <br />
        Last Modified: {moment(f.lastModifiedTime).fromNow()}
      </Card.Content>
    </Card>
  );
};

export default XFileCard;
