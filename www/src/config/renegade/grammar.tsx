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

// interface IParam {
//   name: string;
//   docs: string;
//   type: ParamType;

//   getSignature(): monaco.languages.ParameterInformation;
//   getInsertText(index: number): string;
//   getMarkdownDef(): string;
// }

// enum ParamType {
//   STR = "STR",
//   URL = "URL",
//   DICT = "DICT",
//   OWNERTUPLE = "OWNERTUPLE",
//   PATH = "PATH",
//   BOOL = "BOOL",
//   NONE = "NONE",
//   FILE = "FILE",
//   ERROR = "ERROR"
// }

// class Param implements IParam {
//   public constructor(
//     public name: string,
//     public docs: string,
//     public type: ParamType,
//     public optional: boolean = false
//   ) {}

//   public getSignature(): monaco.languages.ParameterInformation {
//     return {
//       label: this.name,
//       documentation: this.getDocs()
//     };
//   }

//   public getInsertText(index: number): string {
//     switch (this.type) {
//       case ParamType.OWNERTUPLE:
//         return `"\${${index}:owner:group}"`;
//       case ParamType.PATH:
//         return this.optional
//           ? `${this.name}="\${${index}:/path/to/${this.name}}"`
//           : `"\${${index}:/path/to/${this.name}}"`;
//       case ParamType.STR:
//         return this.optional
//           ? `${this.name}="\${${index}:/path/to/${this.name}}"`
//           : `"\${${index}:${this.name}}"`;
//       case ParamType.URL:
//         return `"\${${index}:http://${this.name}.com}"`;
//       case ParamType.BOOL:
//         return `${this.name}=\${${index}:False}`;
//       default:
//         return `\${${index}:${this.name}}`;
//     }
//   }

//   public getDocs(): monaco.IMarkdownString {
//     return {
//       value: `_${this.name}_: __${
//         this.optional ? "?" : ""
//       }${this.type.valueOf()}__  \n\t${this.docs}`
//     };
//   }

//   public getMarkdownDef(): string {
//     return `${this.name}: ${this.getMarkdownType()}`;
//   }

//   public getMarkdownType(): string {
//     switch (this.type) {
//       case ParamType.OWNERTUPLE:
//       case ParamType.PATH:
//       case ParamType.URL:
//       case ParamType.STR:
//         return this.optional ? "__?str__" : "__str__";
//       case ParamType.BOOL:
//         return this.optional ? "__?bool__" : "__bool__";
//       case ParamType.DICT:
//         return this.optional ? "__?dict__" : "__dict__";
//       default:
//         return "__void__";
//     }
//   }
// }

// class Func implements IFunc {
//   public constructor(
//     public name: string,
//     public docs: string,
//     public retVal: ParamType,
//     public params: IParam[]
//   ) {}

//   static From({ name, docs, retVal, params }): Func {
//     return new Func(name, docs, retVal, params);
//   }

//   public getInsertText(): string {
//     return `${this.name}(${this.params
//       .map((p, i) => p.getInsertText(i + 1))
//       .join(", ")})`;
//   }

//   public getDetail(): string {
//     return `${this.name}(${this.params
//       .map(p => p.name + ": " + p.type.valueOf().toLowerCase())
//       .join(", ")}): ${this.retVal.valueOf().toLowerCase()}`;
//   }

//   public getDocs(): monaco.IMarkdownString {
//     return {
//       value: this.docs
//     };
//   }

//   public getSignature(): monaco.languages.SignatureInformation {
//     return {
//       label: this.getDetail(),
//       documentation: {
//         value: this.docs
//       },
//       parameters: this.params.map(p => p.getSignature())
//     };
//   }
// }

// export const BuiltIns: Func[] = [
//   Func.From({
//     name: "file.move",
//     docs: "Move a file to the desired location",
//     retVal: ParamType.ERROR,
//     params: [
//       new Param("file", "TODO", ParamType.FILE),
//       new Param("dstPath", "TODO", ParamType.PATH)
//     ]
//   }),
//   Func.From({
//     name: "file.remove",
//     docs: "Delete a file",
//     retVal: ParamType.ERROR,
//     params: [new Param("file", "TODO", ParamType.FILE)]
//   }),
//   Func.From({
//     name: "file.copy",
//     docs: "Copy the contents of a file to the destination file.",
//     retVal: ParamType.ERROR,
//     params: [
//       new Param("src", "The file with contents to copy", ParamType.FILE),
//       new Param("dst", "The file to write the contents to", ParamType.FILE)
//     ]
//   }),
//   Func.From({
//     name: "file.chown",
//     docs: "Change the file's ownership.",
//     retVal: ParamType.ERROR,
//     params: [
//       new Param("file", "The file to modify", ParamType.FILE),
//       new Param("user", "The user ownership to set", ParamType.STR),
//       new Param("group", "The group ownership to set", ParamType.STR)
//     ]
//   }),
//   Func.From({
//     name: "sys.exec",
//     docs: "Execute a command on the system",
//     retVal: ParamType.ERROR,
//     params: [
//       new Param(
//         "cmd",
//         "The command to run. When running shell commands you'll need to prepend /bin/sh -c or powershell.exe depending on your platform.",
//         ParamType.FILE
//       ),
//       new Param(
//         "disown",
//         "Fork the command into a new process and disown it. Prevents output from being returned.",
//         ParamType.BOOL
//       )
//     ]
//   }),
//   Func.From({
//     name: "sys.detectOS",
//     docs: "Return an enum string indicating the current operating system",
//     retVal: ParamType.STR,
//     params: []
//   }),
//   Func.From({
//     name: "sys.openFile",
//     docs: "Open a file on the system. It will be created if it does not exist.",
//     retVal: ParamType.FILE,
//     params: [new Param("path", "The path to the file.", ParamType.PATH)]
//   }),
//   Func.From({
//     name: "cdn.openFile",
//     docs:
//       "Open a file on the CDN. It will be created (on write) if it does not exist.",
//     retVal: ParamType.FILE,
//     params: [
//       new Param("name", "The name of the file on the CDN.", ParamType.STR)
//     ]
//   }),
//   Func.From({
//     name: "ssh.exec",
//     docs:
//       "Execute a command on the remote system using SSH. Must have valid SSH credentials added to the Target. Only available when using paragon's worker service.",
//     retVal: ParamType.ERROR,
//     params: [
//       new Param(
//         "cmd",
//         "The command to run. When running shell commands you'll need to prepend /bin/sh -c or powershell.exe depending on your platform.",
//         ParamType.FILE
//       )
//     ]
//   }),
//   Func.From({
//     name: "ssh.openFile",
//     docs:
//       "Open a file on the remote system using SFTP over SSH. Must have valid SSH credentials added to the Target. Only available when using paragon's worker service.",
//     retVal: ParamType.FILE,
//     params: [new Param("path", "The remote path to the file.", ParamType.PATH)]
//   }),
//   Func.From({
//     name: "assets.openFile",
//     docs:
//       "Open a file that was packaged into the executable during compilation.",
//     retVal: ParamType.FILE,
//     params: [
//       new Param(
//         "path",
//         "The relative path to the file (assume / is your assets folder).",
//         ParamType.PATH
//       )
//     ]
//   })

// Func.From({
//   name: "chmod",
//   docs:
//     "Chmod uses [os.Chmod](https://godoc.org/os#Chmod) to change a file's permissions. All optional params are assumed to be false unless passed.",
//   retVal: ParamType.NONE,
//   params: [
//     new Param("file", "TODO", ParamType.PATH),
//     new Param("setUser", "TODO", ParamType.BOOL, true),
//     new Param("setGroup", "TODO", ParamType.BOOL, true),
//     new Param("setSticky", "TODO", ParamType.BOOL, true),
//     new Param("ownerRead", "TODO", ParamType.BOOL, true),
//     new Param("ownerWrite", "TODO", ParamType.BOOL, true),
//     new Param("ownerExec", "TODO", ParamType.BOOL, true),
//     new Param("groupRead", "TODO", ParamType.BOOL, true),
//     new Param("groupWrite", "TODO", ParamType.BOOL, true),
//     new Param("groupExec", "TODO", ParamType.BOOL, true),
//     new Param("worldRead", "TODO", ParamType.BOOL, true),
//     new Param("worldWrite", "TODO", ParamType.BOOL, true),
//     new Param("worldExec", "TODO", ParamType.BOOL, true)
//   ]
// }),
// Func.From({
//   name: "chown",
//   docs:
//     "Chown uses [os.Chown](https://godoc.org/os#Chown) to change the user/group ownership of a file/dir.",
//   retVal: ParamType.NONE,
//   params: [
//     new Param("file", "TODO", ParamType.PATH),
//     new Param("owner", "TODO", ParamType.OWNERTUPLE)
//   ]
// }),
// Func.From({
//   name: "connections",
//   docs:
//     "Connections uses [gopsutil/net.ConnectionsPid](https://godoc.org/github.com/shirou/gopsutil/net#ConnectionsPid) to get all the connections opened by a process, given a connection protocol. May work on windows.",
//   retVal: ParamType.DICT,
//   params: []
// }),
// Func.From({
//   name: "copy",
//   docs:
//     "Copy uses [io/ioutil.ReadFile](https://godoc.org/io/ioutil#ReadFile) and [io/ioutil.WriteFile](https://godoc.org/io/ioutil#WriteFile) to copy a file from source to destination.",
//   retVal: ParamType.NONE,
//   params: [
//     new Param("srcFile", "TODO", ParamType.PATH),
//     new Param("dstFile", "TODO", ParamType.PATH)
//   ]
// }),
// Func.From({
//   name: "detectOS",
//   docs:
//     "DetectOS uses [runtime.GOOS](https://godoc.org/runtime#pkg-constants) to detect what OS the agent is running on.",
//   retVal: ParamType.STR,
//   params: []
// }),

// Func.From({
//   name: "dir",
//   docs:
//     "Dir uses [io/ioutil.ReadDir](https://godoc.org/io/ioutil#ReadDir) to get the directory entries of a passed directory.",
//   retVal: ParamType.DICT,
//   params: [new Param("dir", "TODO", ParamType.PATH, true)]
// }),

// Func.From({
//   name: "exec",
//   docs:
//     "Exec uses [os/exec.Command](https://godoc.org/os/exec#Command) to execute the passed string",
//   retVal: ParamType.STR,
//   params: [
//     new Param("cmd", "TODO", ParamType.STR),
//     new Param("disown", "TODO", ParamType.BOOL, true)
//   ]
// }),
// Func.From({
//   name: "kill",
//   docs:
//     "Kill uses [gopsutil/process.Kill](https://godoc.org/github.com/shirou/gopsutil/process#Process.Kill) to kill and passed process pid.",
//   retVal: ParamType.NONE,
//   params: [new Param("pid", "TODO", ParamType.STR)]
// }),

// Func.From({
//   name: "move",
//   docs:
//     "Move uses [os.Rename](https://godoc.org/os#Rename) to move a file from source to destination.",
//   retVal: ParamType.NONE,
//   params: [
//     new Param(
//       "srcFile",
//       "A string for the path of the source file.",
//       ParamType.PATH
//     ),
//     new Param(
//       "dstFile",
//       "A string for the path of the destination file.",
//       ParamType.PATH
//     )
//   ]
// }),
// Func.From({
//   name: "processes",
//   docs:
//     "Processes uses [gopsutil/process.Pids](https://godoc.org/github.com/shirou/gopsutil/process#Pids) to get all pids for a box and then makes them into Process map structs.",
//   retVal: ParamType.DICT,
//   params: []
// }),
// Func.From({
//   name: "read",
//   docs:
//     "Read uses [io/ioutil.ReadFile](https://godoc.org/io/ioutil#ReadFile) to read an entire file's contents.",
//   retVal: ParamType.STR,
//   params: [new Param("file", "TODO", ParamType.PATH)]
// }),
// Func.From({
//   name: "remove",
//   docs:
//     "Remove uses [os.Rename](https://godoc.org/os#Rename) to remove a file/folder. WARNING: basically works like rm -rf.",
//   retVal: ParamType.NONE,
//   params: [new Param("file", "TODO", ParamType.PATH)]
// }),
// Func.From({
//   name: "replaceString",
//   docs:
//     "ReplaceString uses [regexp.MustCompile](https://godoc.org/regexp#MustCompile) to replace values in a string.",
//   retVal: ParamType.STR,
//   params: [
//     new Param("haystack", "TODO", ParamType.STR),
//     new Param("needle", "TODO", ParamType.STR),
//     new Param("newStr", "TODO", ParamType.STR)
//   ]
// }),
// Func.From({
//   name: "request",
//   docs:
//     "Request uses the [net/http](https://godoc.org/net/http) package to send a HTTP request and return a response. Currently we only support GET and POST.",
//   retVal: ParamType.STR,
//   params: [
//     new Param("request", "TODO", ParamType.URL),
//     new Param("method", "TODO", ParamType.STR, true),
//     new Param("writeToFile", "TODO", ParamType.PATH, true),
//     new Param("contentType", "TODO", ParamType.STR, true),
//     new Param("data", "TODO", ParamType.STR, true)
//   ]
// }),
// Func.From({
//   name: "write",
//   docs:
//     "WriteFile uses [io/ioutil.WriteFile](https://godoc.org/io/ioutil#WriteFile) to write an entire file's contents, perms are set to 0644.",
//   retVal: ParamType.NONE,
//   params: [
//     new Param("file", "TODO", ParamType.PATH),
//     new Param("content", "TODO", ParamType.STR)
//   ]
// })
// ];
