import * as monaco from "monaco-editor/esm/vs/editor/editor.api";

export const Theme: monaco.editor.IStandaloneThemeData = {
  base: "vs-dark",
  inherit: true,
  rules: [{ token: "funcDeclName", foreground: "ff6ec7" }],
  colors: {}
};
