import { RESTDataSource } from "apollo-datasource-rest";

export class ParagonAPI extends RESTDataSource {
  constructor() {
    super();
    this.baseURL = "https://enterurlhere";
  }

  async getAllTargers() {
    return this.get("targets");
  }

  async getATarget(values) {
    const result = await this.get("target", {
      values
    });
  }
}
