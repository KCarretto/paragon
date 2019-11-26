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
      { filter, offset, limit }: { filter: String; offset: Number; limit: Number },
      { dataSources }: any
    ) => dataSources.paragonAPI.getAllTargets(filter, offset, limit),
    target: (root: any, { id }: { id: Number }, { dataSources }: any) =>
      dataSources.paragonAPI.getTarget(id)
  },
  Mutation: {
    createTarget: (
      root: any,
      { name, primaryIP, tags }: { name: String; primaryIP: String; tags: [String] },
      { dataSources }: any
    ) => dataSources.paragonAPI.createTarget(name, primaryIP, tags),
    setTargetFields: (root: any, args: setTargetFieldsInput, { dataSources }: any) =>
      dataSources.paragonAPI.setTargetFields(args)
  }
};
