import * as React from "react";
import { useState } from "react";
import { Button, Container, Icon, Menu, Segment, Sidebar } from "semantic-ui-react";
import { XJobEditor, XJobResults } from "../components/job";

const XRunView = () => {
  const [name, setName] = useState<string>("Untitled Job...");
  const [visible, setVisible] = useState<boolean>(false);

  return (
    <React.Fragment>
      <Sidebar.Pushable style={{ overflowY: "hidden" }}>
        <Sidebar
          visible={visible}
          animation='overlay'
          direction='bottom'
          as={Segment}
          id="XResultOverlay"
          style={{ margin: "0px", padding: "0px", border: "none", height: "50vh!important" }}
        >
          <Menu inverted borderless fluid style={{ margin: "0px", borderRadius: "0", backgroundColor: "rgb(63, 63, 63)" }}>
            <Menu.Item header>
              <Icon name="tasks" />
              Results
            </Menu.Item>
            <Menu.Item position="right">
              <Button icon="chevron down" onClick={() => setVisible(!visible)} />
            </Menu.Item>
          </Menu>
          <Container fluid style={{ overflowY: "auto", height: "100%" }}>
            <XJobResults name={name} />
          </Container>
        </Sidebar>
        <Sidebar.Pusher >
          <XJobEditor name={name} setName={setName} />
          <Menu fixed="bottom" style={{ margin: "0px", borderRadius: "0", backgroundColor: "rgb(63, 63, 63)" }}>
            <Menu.Item position="right">
              <Button icon="chevron up" onClick={() => setVisible(!visible)} />
            </Menu.Item>
          </Menu>
          {/* <Segment basic inverted style={{ marginTop: "0px" }}>
          </Segment> */}
        </Sidebar.Pusher>
      </Sidebar.Pushable>
    </React.Fragment >
  );
};
export default XRunView;
