export default {
  Query: {
    testMessage: (): string => "Hello World!",
    targets: (root, args, { dataSources }) => dataSources.paragonAPI.getAllTargets(),
    target: (root, { values }, { dataSources }) => dataSources.paragonAPI.getTarget(values)
  }
};
