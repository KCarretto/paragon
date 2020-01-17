

export const renegade_conf = (monaco) => {
    return {
        comments: {
            lineComment: '#',
            blockComment: ['\'\'\'', '\'\'\''],
        },
        brackets: [
            ['{', '}'],
            ['[', ']'],
            ['(', ')']
        ],
        autoClosingPairs: [
            { open: '{', close: '}' },
            { open: '[', close: ']' },
            { open: '(', close: ')' },
            { open: '"', close: '"', notIn: ['string'] },
            { open: '\'', close: '\'', notIn: ['string', 'comment'] },
        ],
        surroundingPairs: [
            { open: '{', close: '}' },
            { open: '[', close: ']' },
            { open: '(', close: ')' },
            { open: '"', close: '"' },
            { open: '\'', close: '\'' },
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
};

export const renegade_language = {
    defaultToken: '',
    tokenPostfix: '.rg',

    keywords: [
        'and',
        'break',
        'continue',
        'def',
        'del',
        'elif',
        'else',
        'for',
        'from',
        'if',
        'in',
        'None',
        'not',
        'or',
        'pass',
        'print',
        'return',

        'dict',
        'len',
        'list',
        'map',
        'max',
        'min',
        'ord',
        'pow',
        'print',
        'range',
        'set',
        'slice',
        'str',
        'tuple',
        'xrange',
        'zip',

        'True',
        'False',


        'load',
        "move",
        "copy",
        "remove",
        "exec",
        "read",
        "write",
        "chmod",
        "chown",
        "processes",
        "kill",
        "connections",
        "dir",
        "replaceString",
        "request",
        "detectOS",
    ],

    brackets: [
        { open: '{', close: '}', token: 'delimiter.curly' },
        { open: '[', close: ']', token: 'delimiter.bracket' },
        { open: '(', close: ')', token: 'delimiter.parenthesis' }
    ],

    tokenizer: {
        root: [
            { include: '@whitespace' },
            { include: '@numbers' },
            { include: '@strings' },

            [/[,:;]/, 'delimiter'],
            [/[{}\[\]()]/, '@brackets'],

            [/@[a-zA-Z]\w*/, 'tag'],
            [/[a-zA-Z]\w*/, {
                cases: {
                    '@keywords': 'keyword',
                    '@default': 'identifier'
                }
            }]
        ],

        // Deal with white space, including single and multi-line comments
        whitespace: [
            [/\s+/, 'white'],
            [/(^#.*$)/, 'comment'],
            [/'''/, 'string', '@endDocString'],
            [/"""/, 'string', '@endDblDocString']
        ],
        endDocString: [
            [/[^']+/, 'string'],
            [/\\'/, 'string'],
            [/'''/, 'string', '@popall'],
            [/'/, 'string']
        ],
        endDblDocString: [
            [/[^"]+/, 'string'],
            [/\\"/, 'string'],
            [/"""/, 'string', '@popall'],
            [/"/, 'string']
        ],

        // Recognize hex, negatives, decimals, imaginaries, longs, and scientific notation
        numbers: [
            [/-?0x([abcdef]|[ABCDEF]|\d)+[lL]?/, 'number.hex'],
            [/-?(\d*\.)?\d+([eE][+\-]?\d+)?[jJ]?[lL]?/, 'number']
        ],

        // Recognize strings, including those broken across lines with \ (but not without)
        strings: [
            [/'$/, 'string.escape', '@popall'],
            [/'/, 'string.escape', '@stringBody'],
            [/"$/, 'string.escape', '@popall'],
            [/"/, 'string.escape', '@dblStringBody']
        ],
        stringBody: [
            [/[^\\']+$/, 'string', '@popall'],
            [/[^\\']+/, 'string'],
            [/\\./, 'string'],
            [/'/, 'string.escape', '@popall'],
            [/\\$/, 'string']
        ],
        dblStringBody: [
            [/[^\\"]+$/, 'string', '@popall'],
            [/[^\\"]+/, 'string'],
            [/\\./, 'string'],
            [/"/, 'string.escape', '@popall'],
            [/\\$/, 'string']
        ]
    }
};




export const renegade_autocomplete = (monaco) => {
    return {
        provideCompletionItems: () => {
            var suggestions = [{
                label: 'load',
                kind: monaco.languages.CompletionItemKind.Keyword,
                insertText: 'load(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'move',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'move(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'copy',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'copy(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'remove',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'remove(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'exec',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'exec(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'read',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'read(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'write',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'write(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'chmod',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'chmod(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'chown',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'chown(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'processes',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'processes(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'kill',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'kill(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'connections',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'connections(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'dir',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'dir(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'replaceString',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'replaceString(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'request',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'request(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
            {
                label: 'detectOS',
                kind: monaco.languages.CompletionItemKind.Function,
                insertText: 'detectOS(${1:param})',
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            },
                //     {
                //     label: 'testing',
                //     kind: monaco.languages.CompletionItemKind.Keyword,
                //     insertText: 'testing(${1:condition})',
                //     insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet
                // }, {
                //     label: 'ifelse',
                //     kind: monaco.languages.CompletionItemKind.Snippet,
                //     insertText: [
                //         'if (${1:condition}) {',
                //         '\t$0',
                //         '} else {',
                //         '\t',
                //         '}'
                //     ].join('\n'),
                //     insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                //     documentation: 'If-Else Statement'
                //     }
            ];
            return { suggestions: suggestions };
        }
    }
};