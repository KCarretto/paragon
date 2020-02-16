import { useQuery } from "@apollo/react-hooks";
import { ApolloError } from "apollo-client/errors/ApolloError";
import gql from "graphql-tag";
import * as React from "react";
import { useState } from "react";
import { Dropdown, DropdownItemProps, Input } from "semantic-ui-react";
import { Tag, Service } from "../../graphql/models";

// Suggest services for the typeahead.
export const SUGGEST_SERVICES_QUERY = gql`
  query SuggestServices {
    services {
        name
        id
        tag {
          id
          name
        }
      }
  }
`;

type ServicesResult = {
    services: Service[];
};

// XServiceTypeahead adds a service tag field to a form, which has a value of a single tag id.
const XServiceTypeahead = ({ value, onChange, labeled }) => {
    const { loading, error, data: { services = [] } = {} } = useQuery<ServicesResult>(SUGGEST_SERVICES_QUERY);

    let options = services.map(svc => {
        return {
            text: svc.name || "unknown-service",
            value: svc.tag ? svc.tag.id : 0,
        }
    })

    const getDropdown = () => (
        <Dropdown
            placeholder="Select Service"
            icon=""
            fluid
            search
            selection
            error={error !== null}
            loading={loading}
            options={options}
            name="service_tag"
            value={value}
            onChange={onChange}
            style={{
                borderRadius: "0 4px 4px 0"
            }}
        />
    );

    if (labeled) {
        return <Input fluid label="Service" icon="cloud" input={getDropdown()} />;
    }
    return getDropdown();
};

export default XServiceTypeahead;
