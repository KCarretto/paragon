import { RESTDataSource } from "apollo-datasource-rest";

export class ParagonAPI extends RESTDataSource {
  constructor() {
    super();
    this.baseURL = "https://enterurlhere";
  }

  async getAllTargers(filter: String, offset: Number, limit: Number) {
    return this.get("/api/v1/targets", { filter, offset, limit });
  }

  async getATarget(id: Number) {
    const result = await this.get("/api/v1/targets", {
      id
    });
    return result
  }

  async createTarget(name: String; primaryIP: String; tags: [String]) {
    const result = await this.get("/api/v1/targets/create", { name, primaryIP, tags })
    return result
  }
}
