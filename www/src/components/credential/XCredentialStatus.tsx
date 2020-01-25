import { IconProps } from "semantic-ui-react";
import { Credential } from "../../graphql/models";

type CredentialState = {
  text: string;
  icon: IconProps;
};

interface CredentialStatus {
  getStatus(c: Credential): CredentialState;
}

class XCredentialStatus implements CredentialStatus {
  static WORKING: CredentialState = {
    text: "The Credential has never failed.",
    icon: {
      name: "check circle",
      color: "green",
      className: "XCircleIcon",
      bordered: false,
      circular: true
    }
  };

  static UNSURE: CredentialState = {
    text: "The credential has failed atleast once.",
    icon: {
      name: "times circle",
      color: "yellow",
      className: "XCircleIcon",
      bordered: false,
      circular: true
    }
  };

  static FAILING: CredentialState = {
    text: "The credential has failed atleast eight times.",
    icon: {
      name: "times circle",
      color: "red",
      className: "XCircleIcon",
      bordered: false,
      circular: true
    }
  };

  public getStatus(c: Credential): CredentialState {
    if (c.fails > 8) {
      return XCredentialStatus.FAILING;
    } else if (c.fails > 0) {
      return XCredentialStatus.UNSURE;
    }
    return XCredentialStatus.WORKING;
  }
}

export default XCredentialStatus;
