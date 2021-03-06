import * as monaco from "monaco-editor";
import { FunctionSignatures } from "./grammar";

export const SignatureProvider: monaco.languages.SignatureHelpProvider = {
  signatureHelpTriggerCharacters: ["(", ","],
  // signatureHelpRetriggerCharacters: [','],
  provideSignatureHelp: (model, position, token, context) => {
    // Default value if no signatures are found
    let noSignatures = {
      value: { signatures: [], activeParameter: -1, activeSignature: -1 },
      dispose: () => {}
    };
    let startPos = position;
    let endPos = position;

    // Start of function call
    let startMatch = model.findPreviousMatch(
      /\w+(?=\()/.source,
      position,
      true,
      true,
      null,
      true
    );
    if (!startMatch) {
      console.log("NO PREV FUNCTION", startPos, endPos, startMatch);
      return noSignatures;
    }
    startPos = startMatch.range.getEndPosition();

    // End of function call
    let endMatch = model.findNextMatch(
      /\)/.source,
      position,
      true,
      true,
      null,
      true
    );
    if (endMatch) {
      endPos = endMatch.range.getEndPosition();
    }

    // Out of function call bounds
    if (endPos.isBefore(position)) {
      console.log("POS OUT OF BOUNDS", startPos, endPos, startMatch, endMatch);
      return noSignatures;
    }

    // Find signature index
    let funcName = model.getWordUntilPosition(startPos).word;
    let sig = FunctionSignatures.get(funcName);

    // No signature found
    if (!sig) {
      console.log(
        "NO SIG INDEX",
        startPos,
        endPos,
        startMatch,
        endMatch,
        funcName
      );
      return noSignatures;
    }

    // Get all param characters within the function call
    let funcBody = model.getValueInRange({
      startColumn: startPos.column,
      startLineNumber: startPos.lineNumber,
      endLineNumber: position.lineNumber,
      endColumn: position.column
    });

    // Determine param index
    let funcBodyTokens = funcBody.match(
      /(?:(?<!['"])\b\w+|['"][^'"]*['"])\s*,/
    );
    let paramIndex =
      !funcBodyTokens || funcBodyTokens.length <= 0
        ? 0
        : funcBodyTokens.length >= sig.parameters.length
        ? sig.parameters.length - 1
        : funcBodyTokens.length;

    console.log(
      "SIG FOUND",
      startPos,
      endPos,
      startMatch,
      endMatch,
      funcName,
      sig,
      funcBody,
      funcBodyTokens,
      paramIndex
    );
    return {
      value: {
        signatures: [sig],
        activeSignature: 0,
        activeParameter: paramIndex
      },
      dispose: () => {}
    };
  }
};
