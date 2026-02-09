import * as vscode from 'vscode';
import * as path from 'path';

type KeywordMap = {
  keywords: Record<string, string>;
  predeclared: Record<string, string>;
};

type WorkspaceState = {
  locale: string;
  keywords: string[];
  predeclared: string[];
};

const defaultLocale = 'bn';

export function activate(context: vscode.ExtensionContext) {
  const stateByFolder = new Map<string, WorkspaceState>();

  const provider: vscode.CompletionItemProvider = {
    provideCompletionItems(doc, _pos) {
      if (!doc.fileName.endsWith('.p.go')) {
        return undefined;
      }
      const folder = workspaceFolderForDoc(doc);
      if (!folder) {
        return undefined;
      }
      const state = stateByFolder.get(folder.uri.fsPath);
      if (!state) {
        return undefined;
      }
      const items: vscode.CompletionItem[] = [];
      for (const kw of state.keywords) {
        items.push(new vscode.CompletionItem(kw, vscode.CompletionItemKind.Keyword));
      }
      for (const kw of state.predeclared) {
        items.push(new vscode.CompletionItem(kw, vscode.CompletionItemKind.Keyword));
      }
      return items;
    },
  };

  context.subscriptions.push(vscode.languages.registerCompletionItemProvider({ language: 'go' }, provider));

  // Initialize
  refreshAll(stateByFolder).catch(console.error);

  // Watch for locale/config changes
  const langWatcher = vscode.workspace.createFileSystemWatcher('**/lang/*.json');
  const localeWatcher = vscode.workspace.createFileSystemWatcher('**/.pgo_lang');
  const configWatcher = vscode.workspace.onDidChangeConfiguration((e) => {
    if (e.affectsConfiguration('polygo.locale')) {
      refreshAll(stateByFolder).catch(console.error);
    }
  });

  const refresh = () => refreshAll(stateByFolder).catch(console.error);
  langWatcher.onDidChange(refresh);
  langWatcher.onDidCreate(refresh);
  langWatcher.onDidDelete(refresh);
  localeWatcher.onDidChange(refresh);
  localeWatcher.onDidCreate(refresh);
  localeWatcher.onDidDelete(refresh);

  context.subscriptions.push(langWatcher, localeWatcher, configWatcher);
}

export function deactivate() {}

async function refreshAll(stateByFolder: Map<string, WorkspaceState>) {
  const folders = vscode.workspace.workspaceFolders ?? [];
  for (const folder of folders) {
    const state = await loadWorkspaceState(folder);
    if (state) {
      stateByFolder.set(folder.uri.fsPath, state);
    }
  }
}

async function loadWorkspaceState(folder: vscode.WorkspaceFolder): Promise<WorkspaceState | null> {
  const configLocale = vscode.workspace.getConfiguration('polygo', folder).get<string>('locale')?.trim();
  const locale = configLocale || (await readLocaleFile(folder)) || defaultLocale;
  const map = await readLangMap(folder, locale);
  if (!map) {
    return null;
  }
  return {
    locale,
    keywords: Object.keys(map.keywords || {}),
    predeclared: Object.keys(map.predeclared || {}),
  };
}

async function readLocaleFile(folder: vscode.WorkspaceFolder): Promise<string | null> {
  const file = vscode.Uri.joinPath(folder.uri, '.pgo_lang');
  try {
    const data = await vscode.workspace.fs.readFile(file);
    const value = data.toString().trim();
    return value || null;
  } catch {
    return null;
  }
}

async function readLangMap(folder: vscode.WorkspaceFolder, locale: string): Promise<KeywordMap | null> {
  const file = vscode.Uri.joinPath(folder.uri, 'lang', `${locale}.json`);
  try {
    const data = await vscode.workspace.fs.readFile(file);
    return JSON.parse(data.toString()) as KeywordMap;
  } catch {
    return null;
  }
}

function workspaceFolderForDoc(doc: vscode.TextDocument): vscode.WorkspaceFolder | undefined {
  return vscode.workspace.getWorkspaceFolder(doc.uri);
}
