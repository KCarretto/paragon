import * as React from "react";
import { Dropdown, Input } from "semantic-ui-react";

const XCredentialKindDropdown = ({ value, onChange }) => (
  <Input
    fluid
    label="Kind"
    input={
      <Dropdown
        fluid
        selection
        labeled
        options={[
          {
            text: "Password",
            value: "password"
          },
          {
            text: "SSH Private Key",
            value: "key"
          }
        ]}
        name="kind"
        value={value}
        onChange={onChange}
        style={{
          borderRadius: "0 4px 4px 0"
        }}
      />
    }
  />
);

export default XCredentialKindDropdown;
