import * as monaco from "monaco-editor/esm/vs/editor/editor.api";

interface IFunc {
  name: string;
  docs: string;
  retVal: ParamType;
  params: IParam[];

  getSignature(): monaco.languages.SignatureInformation;
}

interface IParam {
  name: string;
  docs: string;
  type: ParamType;

  getSignature(): monaco.languages.ParameterInformation;
  getInsertText(index: number): string;
  getMarkdownDef(): string;
}

enum ParamType {
  STR = "STR",
  URL = "URL",
  DICT = "DICT",
  OWNERTUPLE = "OWNERTUPLE",
  PATH = "PATH",
  BOOL = "BOOL",
  NONE = "NONE"
}

class Param implements IParam {
  public constructor(
    public name: string,
    public docs: string,
    public type: ParamType,
    public optional: boolean = false
  ) {}

  public getSignature(): monaco.languages.ParameterInformation {
    return {
      label: this.name,
      documentation: this.getDocs()
    };
  }

  public getInsertText(index: number): string {
    switch (this.type) {
      case ParamType.OWNERTUPLE:
        return `"\${${index}:owner:group}"`;
      case ParamType.PATH:
        return `"\${${index}:/path/to/${this.name}}"`;
      case ParamType.STR:
        return `"\${${index}:${this.name}}"`;
      case ParamType.URL:
        return `"\${${index}:http://${this.name}.com}"`;
      case ParamType.BOOL:
        return `${this.name}=\${${index}:False}`;
      default:
        return `\${${index}:${this.name}}`;
    }
  }

  public getDocs(): monaco.IMarkdownString {
    return {
      value: `_${this.name}_: __${
        this.optional ? "?" : ""
      }${this.type.valueOf()}__  \n\t${this.docs}`
    };
  }

  public getMarkdownDef(): string {
    return `${this.name}: ${this.getMarkdownType()}`;
  }

  public getMarkdownType(): string {
    switch (this.type) {
      case ParamType.OWNERTUPLE:
      case ParamType.PATH:
      case ParamType.URL:
      case ParamType.STR:
        return this.optional ? "__?str__" : "__str__";
      case ParamType.BOOL:
        return this.optional ? "__?bool__" : "__bool__";
      case ParamType.DICT:
        return this.optional ? "__?dict__" : "__dict__";
      default:
        return "__void__";
    }
  }
}

class Func implements IFunc {
  public constructor(
    public name: string,
    public docs: string,
    public retVal: ParamType,
    public params: IParam[]
  ) {}

  static From({ name, docs, retVal, params }): Func {
    return new Func(name, docs, retVal, params);
  }

  public getInsertText(): string {
    return `${this.name}(${this.params
      .map((p, i) => p.getInsertText(i + 1))
      .join(", ")})`;
  }

  public getDetail(): string {
    return `${this.name}(${this.params
      .map(p => p.name + ": " + p.type.valueOf().toLowerCase())
      .join(", ")}): ${this.retVal.valueOf().toLowerCase()}`;
  }

  public getDocs(): monaco.IMarkdownString {
    return {
      value: this.docs
    };
  }

  public getSignature(): monaco.languages.SignatureInformation {
    return {
      label: this.name,
      documentation: {
        value: this.docs
      },
      parameters: this.params.map(p => p.getSignature())
    };
  }
}

export const BuiltIns: Func[] = [
  Func.From({
    name: "chmod",
    docs:
      "Chmod uses [os.Chmod](https://godoc.org/os#Chmod) to change a file's permissions. All optional params are assumed to be false unless passed.",
    retVal: ParamType.NONE,
    params: [
      new Param("file", "TODO", ParamType.PATH),
      new Param("setUser", "TODO", ParamType.BOOL, true),
      new Param("setGroup", "TODO", ParamType.BOOL, true),
      new Param("setSticky", "TODO", ParamType.BOOL, true),
      new Param("ownerRead", "TODO", ParamType.BOOL, true),
      new Param("ownerWrite", "TODO", ParamType.BOOL, true),
      new Param("ownerExec", "TODO", ParamType.BOOL, true),
      new Param("groupRead", "TODO", ParamType.BOOL, true),
      new Param("groupWrite", "TODO", ParamType.BOOL, true),
      new Param("groupExec", "TODO", ParamType.BOOL, true),
      new Param("worldRead", "TODO", ParamType.BOOL, true),
      new Param("worldWrite", "TODO", ParamType.BOOL, true),
      new Param("worldExec", "TODO", ParamType.BOOL, true)
    ]
  }),
  Func.From({
    name: "chown",
    docs:
      "Chown uses [os.Chown](https://godoc.org/os#Chown) to change the user/group ownership of a file/dir.",
    retVal: ParamType.NONE,
    params: [
      new Param("file", "TODO", ParamType.PATH),
      new Param("owner", "TODO", ParamType.OWNERTUPLE)
    ]
  }),
  Func.From({
    name: "connections",
    docs:
      "Connections uses [gopsutil/net.ConnectionsPid](https://godoc.org/github.com/shirou/gopsutil/net#ConnectionsPid) to get all the connections opened by a process, given a connection protocol. May work on windows.",
    retVal: ParamType.DICT,
    params: []
  }),
  Func.From({
    name: "copy",
    docs:
      "Copy uses [io/ioutil.ReadFile](https://godoc.org/io/ioutil#ReadFile) and [io/ioutil.WriteFile](https://godoc.org/io/ioutil#WriteFile) to copy a file from source to destination.",
    retVal: ParamType.NONE,
    params: [
      new Param("srcFile", "TODO", ParamType.PATH),
      new Param("dstFile", "TODO", ParamType.PATH)
    ]
  }),
  Func.From({
    name: "detectOS",
    docs:
      "DetectOS uses [runtime.GOOS](https://godoc.org/runtime#pkg-constants) to detect what OS the agent is running on.",
    retVal: ParamType.STR,
    params: []
  }),

  Func.From({
    name: "dir",
    docs:
      "Dir uses [io/ioutil.ReadDir](https://godoc.org/io/ioutil#ReadDir) to get the directory entries of a passed directory.",
    retVal: ParamType.DICT,
    params: [new Param("dir", "TODO", ParamType.PATH, true)]
  }),

  Func.From({
    name: "exec",
    docs:
      "Exec uses [os/exec.Command](https://godoc.org/os/exec#Command) to execute the passed string",
    retVal: ParamType.STR,
    params: [
      new Param("cmd", "TODO", ParamType.STR),
      new Param("disown", "TODO", ParamType.BOOL, true)
    ]
  }),
  Func.From({
    name: "kill",
    docs:
      "Kill uses [gopsutil/process.Kill](https://godoc.org/github.com/shirou/gopsutil/process#Process.Kill) to kill and passed process pid.",
    retVal: ParamType.NONE,
    params: [new Param("pid", "TODO", ParamType.STR)]
  }),

  Func.From({
    name: "move",
    docs:
      "Move uses [os.Rename](https://godoc.org/os#Rename) to move a file from source to destination.",
    retVal: ParamType.NONE,
    params: [
      new Param(
        "srcFile",
        "A string for the path of the source file.",
        ParamType.PATH
      ),
      new Param(
        "dstFile",
        "A string for the path of the destination file.",
        ParamType.PATH
      )
    ]
  }),
  Func.From({
    name: "processes",
    docs:
      "Processes uses [gopsutil/process.Pids](https://godoc.org/github.com/shirou/gopsutil/process#Pids) to get all pids for a box and then makes them into Process map structs.",
    retVal: ParamType.DICT,
    params: []
  }),
  Func.From({
    name: "read",
    docs:
      "Read uses [io/ioutil.ReadFile](https://godoc.org/io/ioutil#ReadFile) to read an entire file's contents.",
    retVal: ParamType.STR,
    params: [new Param("file", "TODO", ParamType.PATH)]
  }),
  Func.From({
    name: "remove",
    docs:
      "Remove uses [os.Rename](https://godoc.org/os#Rename) to remove a file/folder. WARNING: basically works like rm -rf.",
    retVal: ParamType.NONE,
    params: [new Param("file", "TODO", ParamType.PATH)]
  }),
  Func.From({
    name: "replaceString",
    docs:
      "ReplaceString uses [regexp.MustCompile](https://godoc.org/regexp#MustCompile) to replace values in a string.",
    retVal: ParamType.STR,
    params: [
      new Param("haystack", "TODO", ParamType.STR),
      new Param("needle", "TODO", ParamType.STR),
      new Param("newStr", "TODO", ParamType.STR)
    ]
  }),
  Func.From({
    name: "request",
    docs:
      "Request uses the [net/http](https://godoc.org/net/http) package to send a HTTP request and return a response. Currently we only support GET and POST.",
    retVal: ParamType.STR,
    params: [
      new Param("request", "TODO", ParamType.URL),
      new Param("method", "TODO", ParamType.STR),
      new Param("writeToFile", "TODO", ParamType.PATH),
      new Param("contentType", "TODO", ParamType.STR),
      new Param("data", "TODO", ParamType.STR)
    ]
  }),
  Func.From({
    name: "write",
    docs:
      "WriteFile uses [io/ioutil.WriteFile](https://godoc.org/io/ioutil#WriteFile) to write an entire file's contents, perms are set to 0644.",
    retVal: ParamType.NONE,
    params: [
      new Param("file", "TODO", ParamType.PATH),
      new Param("content", "TODO", ParamType.STR)
    ]
  })
];
