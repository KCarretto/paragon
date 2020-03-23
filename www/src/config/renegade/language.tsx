import * as monaco from "monaco-editor";
import * as spec from "./spec.json";

interface ILanguage extends monaco.languages.IMonarchLanguage {
  keywords: string[];
}

export const LanguageConfig: monaco.languages.LanguageConfiguration = {
  comments: { lineComment: "#", blockComment: ["'''", "'''"] },
  brackets: [
    ["{", "}"],
    ["[", "]"],
    ["(", ")"]
  ],
  autoClosingPairs: [
    { open: "{", close: "}" },
    { open: "[", close: "]" },
    { open: "(", close: ")" },
    { open: '"', close: '"', notIn: ["string"] },
    { open: "'", close: "'", notIn: ["string", "comment"] }
  ],
  surroundingPairs: [
    { open: "{", close: "}" },
    { open: "[", close: "]" },
    { open: "(", close: ")" },
    { open: '"', close: '"' },
    { open: "'", close: "'" }
  ],
  onEnterRules: [
    {
      beforeText: new RegExp("^\\s*(?:def|for|if|elif|else).*?:\\s*$"),
      action: { indentAction: monaco.languages.IndentAction.Indent }
    }
  ],
  folding: {
    offSide: true,
    markers: {
      start: new RegExp("^\\s*#region\\b"),
      end: new RegExp("^\\s*#endregion\\b")
    }
  }
};

export const Language: ILanguage = {
  defaultToken: "",
  tokenPostfix: ".rg",
  keywords: [
    "and",
    "break",
    "continue",
    "def",
    "del",
    "elif",
    "else",
    "for",
    "from",
    "if",
    "in",
    "None",
    "not",
    "or",
    "pass",
    "print",
    "return",

    "dict",
    "len",
    "list",
    "map",
    "max",
    "min",
    "ord",
    "pow",
    "print",
    "range",
    "set",
    "slice",
    "str",
    "tuple",
    "xrange",
    "zip",

    "True",
    "False",
    ...spec.libraries.map(({ name }) => name)
    // "load",
    // "move",
    // "copy",
    // "remove",
    // "exec",
    // "read",
    // "write",
    // "chmod",
    // "chown",
    // "processes",
    // "kill",
    // "connections",
    // "dir",
    // "replaceString",
    // "request",
    // "detectOS"
  ],

  brackets: [
    { open: "{", close: "}", token: "delimiter.curly" },
    { open: "[", close: "]", token: "delimiter.bracket" },
    { open: "(", close: ")", token: "delimiter.parenthesis" }
  ],

  tokenizer: {
    root: [
      { include: "@whitespace" },
      { include: "@numbers" },
      { include: "@strings" },
      { include: "@functions" },

      [/[,:;]/, "delimiter"],
      [/[{}\[\]()]/, "@brackets"],

      [/@[a-zA-Z]\w*/, "tag"],
      [
        /[a-zA-Z]\w*/,
        { cases: { "@keywords": "keyword", "@default": "identifier" } }
      ]
    ],

    functions: [
      // Function Declaration
      [/(def)(\s+)(\w+)/, ["keyword", "white", "funcDeclName"]]
    ],
    // Deal with white space, including single and multi-line comments
    whitespace: [
      [/\s+/, "white"],
      [/(^#.*$)/, "comment"],
      [/'''/, "string", "@endDocString"],
      [/"""/, "string", "@endDblDocString"]
    ],
    endDocString: [
      [/[^']+/, "string"],
      [/\\'/, "string"],
      [/'''/, "string", "@popall"],
      [/'/, "string"]
    ],
    endDblDocString: [
      [/[^"]+/, "string"],
      [/\\"/, "string"],
      [/"""/, "string", "@popall"],
      [/"/, "string"]
    ],

    // Recognize hex, negatives, decimals, imaginaries, longs, and scientific
    // notation
    numbers: [
      [/-?0x([abcdef]|[ABCDEF]|\d)+[lL]?/, "number.hex"],
      [/-?(\d*\.)?\d+([eE][+\-]?\d+)?[jJ]?[lL]?/, "number"]
    ],

    // Recognize strings, including those broken across lines with \ (but not
    // without)
    strings: [
      [/'$/, "string.escape", "@popall"],
      [/'/, "string.escape", "@stringBody"],
      [/"$/, "string.escape", "@popall"],
      [/"/, "string.escape", "@dblStringBody"]
    ],
    stringBody: [
      [/[^\\']+$/, "string", "@popall"],
      [/[^\\']+/, "string"],
      [/\\./, "string"],
      [/'/, "string.escape", "@popall"],
      [/\\$/, "string"]
    ],
    dblStringBody: [
      [/[^\\"]+$/, "string", "@popall"],
      [/[^\\"]+/, "string"],
      [/\\./, "string"],
      [/"/, "string.escape", "@popall"],
      [/\\$/, "string"]
    ]
  }
};
