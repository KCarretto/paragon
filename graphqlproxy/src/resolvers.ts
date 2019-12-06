interface setTargetFieldsInput {
  id: Number;
  name: String;
  machineUUID: String;
  primaryIP: String;
  publicIP: String;
  primaryMAC: String;
  hostname: String;
}

export default {
  Query: {
    testMessage: (): string => "Hello World!",
    targets: (
      root: any,
      args: { filter: String; offset: Number; limit: Number },
      { dataSources }: any
    ) => dataSources.paragonAPI.getAllTargets(args),
    target: (root: any, { id }: { id: String }, { dataSources }: any) =>
      dataSources.paragonAPI.getTarget(id)
  },
  Mutation: {
    createTarget: (
      root: any,
      args: { name: String; primaryIP: String; tags: [String] },
      { dataSources }: any
    ) => dataSources.paragonAPI.createTarget(args),
    setTargetFields: (root: any, args: setTargetFieldsInput, { dataSources }: any) =>
      dataSources.paragonAPI.setTargetFields(args)
  }
};
