import { IconProps } from "semantic-ui-react";
import { Credential } from "../../graphql/models";

type CredentialState = {
  text: string;
  icon: IconProps;
};

interface CredentialStatus {
  WORKING: CredentialState;
  UNSURE: CredentialState;
  FAILING: CredentialState;
  getStatus(c: Credential): CredentialState;
}

class XCredentialStatus implements CredentialStatus {
  WORKING: {
    text: "The Credential has never failed.";
    icon: {
      name: "check circle";
      color: "green";
      className: "XCircleIcon";
      bordered: false;
      circular: true;
    };
  };

  UNSURE: {
    text: "The credential has failed atleast once.";
    icon: {
      name: "times circle";
      color: "yellow";
      className: "XCircleIcon";
      bordered: false;
      circular: true;
    };
  };

  FAILING: {
    text: "The credential has failed atleast eight times.";
    icon: {
      name: "times circle";
      color: "red";
      className: "XCircleIcon";
      bordered: false;
      circular: true;
    };
  };

  public getStatus(c: Credential): CredentialState {
    if (c.fails > 8) {
      return this.FAILING;
    } else if (c.fails > 0) {
      return this.UNSURE;
    }
    return this.WORKING;
  }
}

export default XCredentialStatus;
