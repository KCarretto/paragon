import * as monaco from "monaco-editor";
import * as spec from "./spec.json";

interface IParam {
  name: string;
  type: string;
}

interface IFunc {
  name: string;
  doc: string;
  params?: IParam[];
  retvals?: IParam[];
  // getSignature(): monaco.languages.SignatureInformation;
}

export const GetSignature = (
  libName: string,
  fn: IFunc
): monaco.languages.SignatureInformation => {
  let params =
    fn.params && fn.params.length > 0
      ? fn.params
        .map(param => {
          return `${param.name}: ${param.type}`;
        })
        .join(", ")
      : "";
  let retvals =
    fn.retvals && fn.retvals.length > 0
      ? " -> " +
      fn.retvals
        .map(retval => {
          return `${retval.name}: ${retval.type}`;
        })
        .join(", ")
      : "";

  return {
    label: `${libName}.${fn.name}(${params})${retvals}`,
    documentation: {
      value: fn.doc
    },
    parameters:
      fn.params && fn.params.length > 0
        ? fn.params.map(param => {
          return {
            label: param.name
          };
        })
        : []
  };
};

export const GetInsertText = (libName: string, fn: IFunc): string => {
  let name = `${libName}.${fn.name}`;

  let params =
    fn.params && fn.params.length > 0
      ? fn.params
        .filter(param => !param.type.startsWith("?"))
        .map((param, index) => {
          switch (param.type) {
            case "string":
              return `"\${${index + 1}:${param.name}}"`;
            default:
              return `\${${index + 1}:${param.name}}`;
          }
        })
        .join(", ")
      : "";

  let retvals =
    fn.retvals && fn.retvals.length > 0
      ? fn.retvals.map(retval => `${retval.name}`).join(", ") + " = "
      : "";

  return `${retvals}${name}(${params})`;
};

export const FunctionSignatures = new Map(
  spec.libraries.flatMap(lib => {
    return lib.functions && lib.functions.length > 0
      ? lib.functions.map(fn => [fn.name, GetSignature(lib.name, fn)])
      : [];
  })
);
