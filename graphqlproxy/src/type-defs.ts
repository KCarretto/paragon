import { gql } from "apollo-server";

export default gql`
  type setTargetFieldsResponse {
    id: ID
  }
  type Target {
    name: String
    machineUUID: String
    primaryIP: String
    publicIP: String
    primaryMAC: String
    hostname: String
    lastSeen: Int
    tags: [Int]
    tasks: [Int]
    credentials: [Int]
  }
  type Query {
    testMessage: String!
    targets: [ID]
    target(id: ID): Target
  }
  type Mutation {
    createTarget(name: String, primaryIP: String, tags: [String]): String
    setTargetFields(
      id: ID!
      name: String
      machineUUID: String
      primaryIP: String
      publicIP: String
      primaryMAC: String
      hostname: String
    ): String
  }
`;
