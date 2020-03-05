import * as monaco from "monaco-editor";
import { GetInsertText } from "./grammar";
import * as spec from "./spec.json";

export const CompletionProvider: monaco.languages.CompletionItemProvider = {
  provideCompletionItems: (model, position) => {
    let word = model.getWordUntilPosition(position);
    let range = {
      startLineNumber: position.lineNumber,
      startColumn: word.startColumn,
      endLineNumber: position.lineNumber,
      endColumn: word.endColumn
    };

    let suggestions = spec.libraries.flatMap(lib =>
      lib.functions.map(fn => {
        return {
          label: `${lib.name}.${fn.name}`,
          // detail: func.getDetail(),
          documentation: fn.doc,
          insertText: GetInsertText(lib.name, fn),
          insertTextRules:
            monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
          kind: monaco.languages.CompletionItemKind.Function,
          range: range,
          command: {
            id: "editor.action.triggerParameterHints",
            title: "editor.action.triggerParameterHints"
          }
        };
      })
    );

    return {
      suggestions: suggestions
      //   BuiltIns.map(func => {
      //   return {
      //     label: func.name,
      //     detail: func.getDetail(),
      //     documentation: func.getDocs(),
      //     insertText: func.getInsertText(),
      //     insertTextRules:
      //       monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
      //     kind: monaco.languages.CompletionItemKind.Function,
      //     range: range
      //     // command: {
      //     //   id: "editor.action.triggerParameterHint",
      //     //   title: "editor.action.triggerParameterHint"
      //     // }
      //   };
      // })
    };
  }
};
