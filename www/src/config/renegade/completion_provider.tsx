import * as monaco from "monaco-editor";
import { BuiltIns } from "./grammar";

export const CompletionProvider: monaco.languages.CompletionItemProvider = {
  provideCompletionItems: (model, position) => {
    let word = model.getWordUntilPosition(position);
    let range = {
      startLineNumber: position.lineNumber,
      startColumn: word.startColumn,
      endLineNumber: position.lineNumber,
      endColumn: word.endColumn
    };
    return {
      suggestions: BuiltIns.map(func => {
        return {
          label: func.name,
          detail: func.getDetail(),
          documentation: func.getDocs(),
          insertText: func.getInsertText(),
          insertTextRules:
            monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          kind: monaco.languages.CompletionItemKind.Function,
          range: range
          // command: {
          //   id: "editor.action.triggerParameterHint",
          //   title: "editor.action.triggerParameterHint"
          // }
        };
      })
    };
  }
};
