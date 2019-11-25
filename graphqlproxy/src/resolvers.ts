export default {
  Query: {
    testMessage: (): string => "Hello World!",
    targets: (
      root,
      { filter, offset, limit }: { filter: String; offset: Number; limit: Number },
      { dataSources }
    ) => dataSources.paragonAPI.getAllTargets(filter, offset, limit),
    target: (root, { id }: { id: Number }, { dataSources }) => dataSources.paragonAPI.getTarget(id)
  },
  Mutation: {
    createTarget: (
      root,
      { name, primaryIP, tags }: { name: String; primaryIP: String; tags: [String] },
      { dataSources }
    ) => dataSources.paragonAPI.createTarget(name, primaryIP, tags)
  }
};
