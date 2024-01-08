export {};

// Import complete from "sql-language-server/complete";
// import * as monaco from "monaco-editor";

// import { Schema } from "sql-language-server/database_libs/AbstractClient";

// type getColumns = (databaseName: string, tableName: string) => string[];
// type getTables = (databaseName: string) => string[];
// type getDatabases = () => string[];

// // Pos from sql-language-server
// type SQLLanguageServerPos = { line: number; column: number };

// export class SqlCompletionProvider implements monaco.languages.CompletionItemProvider {
//     private schema: Schema;

//     public triggerCharacters = [",", ".", " "];
//     constructor(schema: Schema) {
//         this.schema = schema;
//     }

//     fromMonacoPosition(position: monaco.Position): SQLLanguageServerPos {
//         return {
//             line: position.lineNumber - 1,
//             column: position.column - 1
//         };
//     }

//     public provideCompletionItems(
//         model: monaco.editor.ITextModel,
//         position: monaco.Position,
//     ): monaco.languages.ProviderResult<monaco.languages.CompletionList>
//     {
//         const completionsFromSQLLanguageServer = complete(
//             model.getValue(),
//             this.fromMonacoPosition(position),
//             this.schema
//         );

//         const completions: monaco.languages.CompletionList = {
//             suggestions: []
//         };
//         completionsFromSQLLanguageServer.candidates.forEach(comp => {
//             let monacoCompletion: monaco.languages.CompletionItem = {
//                 label: comp.label,
//                 kind: comp.kind,
//                 insertText: comp.insertText || comp.label,
//                 range: null,
//                 documentation: comp.detail ? comp.detail.replace(/^column/, "") : ""
//             };

//             completions.suggestions.push(monacoCompletion);
//         });
//         return completions;
//     }
// }
