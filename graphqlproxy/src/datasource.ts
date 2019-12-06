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
    const result = await this.get("/api/v1/targets", args);
    console.log(result);
    return result ? result : [];
  }

  async getTarget(id: String) {
    const result = await this.get(`/api/v1/targets/${id}`);
    return result ? result : {};
  }

  async createTarget(args: any) {
    const result = await this.post("/api/v1/targets/create", args);
    return result ? result.id : null;
  }

  async setTargetFields(args: {
    id: String;
    name: String;
    machineUUID: String;
    primaryIP: String;
    publicIP: String;
    primaryMAC: String;
    hostname: String;
  }) {
    const result = await this.post("/api/v1/targets/setTargetFields", args);
    return true;
  }
}
