import { RESTDataSource } from "apollo-datasource-rest";

interface TargetsFilterInput {
  filter: String;
  offset: Number;
  limit: Number;
}

export class ParagonAPI extends RESTDataSource {
  constructor() {
    super();
    this.baseURL = "http://127.0.0.1:80";
  }

  async getAllTargets(args: any) {
    console.log("WAT");
    const result = await this.get("/api/v1/targets", args);
    console.log(result);
    return [];
  }

  async getATarget(id: Number) {
    const result = await this.get("/api/v1/targets", {
      id
    });
    return result;
  }

  async createTarget(name: String, primaryIP: String, tags: [String]) {
    const result = await this.post("/api/v1/targets/create", { name, primaryIP, tags });
    return result;
  }

  async setTargetFields(args: {
    id: Number;
    name: String;
    machineUUID: String;
    primaryIP: String;
    publicIP: String;
    primaryMAC: String;
    hostname: String;
  }) {
    const result = await this.post("/api/v1/targets/setTargetFields", args);
    return result;
  }
}
