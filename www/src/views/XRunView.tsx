import * as React from "react";
import { useState } from "react";
import { XJobEditor, XJobResults } from "../components/job";

const XRunView = () => {
  const [name, setName] = useState<string>("Untitled Job...");
  return (
    <React.Fragment>
      <XJobEditor name={name} setName={setName} />
      <XJobResults name={name} />
    </React.Fragment>
  );
};
export default XRunView;
