// import { Monaco } from './monaco';
import * as monaco from 'monaco-editor';
import {CompletionProvider} from './completion_provider';
import {Language, LanguageConfig} from './language';
import {SignatureProvider} from './signature_provider';
import {Theme} from './theme';

export const LanguageID = 'renegade';

export const Register = () => {
  monaco.languages.register({id: LanguageID});

  monaco.editor.defineTheme(LanguageID, Theme);
  monaco.editor.setTheme(LanguageID);
  monaco.languages.setMonarchTokensProvider(LanguageID, Language);
  monaco.languages.setLanguageConfiguration(LanguageID, LanguageConfig);

  monaco.languages.registerCompletionItemProvider(
      LanguageID, CompletionProvider);
  monaco.languages.registerSignatureHelpProvider(LanguageID, SignatureProvider);
}
