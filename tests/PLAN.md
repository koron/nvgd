# OPFS Web UI Playwright E2E テスト計画

対象: NVGD の OPFS プロトコル UI (`/opfs/`)
バージョン: 計画 v1 (2026-05-27)

---

## 1. ゴールと方針

### ゴール
NVGD の OPFS Web UI (`protocol/opfs/assets/index.html` + `main.js`) を Chromium / Firefox で E2E 検証し、DuckDB 連携を含むユーザーフロー全体が壊れないことを保証する。

### スコープ
- ディレクトリの作成・移動・履歴ナビゲーション
- ローカルファイルのアップロード（同名上書き確認を含む）
- 簡易エディタによるファイル作成・更新、Edit アクションでのロード
- URL からのダウンロード（同名上書き確認を含む）
- ファイル/ディレクトリの選択・一括削除
- DuckDB 連携（新規タブで対応 URL が開くこと、生成クエリの妥当性）
- Reload による再描画、Mithril 再レンダリング後の DOM 一貫性
- パンくず・ブラウザ戻る/進む（popstate）
- File System Access API のサポート差異（Chromium のみ Save As 検証）

### スコープ外
- DuckDB WASM 内部の実クエリ実行結果の検証（クエリ生成 URL までで止める）
- バックエンド (Go) のロジックテスト（既存の `opfs_test.go` でカバー済み）
- WebKit（Safari）— OPFS と File System Access API の互換性が不安定なため除外
- モバイルブラウザ
- アクセシビリティ・i18n の網羅検証
- 視覚回帰（必要であれば後段で `toHaveScreenshot` を追加）

### 非機能要件
- 1 テスト = 1 シナリオ、独立に実行可能
- 並列実行可能（context ごとに OPFS は隔離される）
- 各テストは自前のテンポラリ OPFS から開始（前テストの残骸に依存しない）

---

## 2. 前提環境

### サーバー起動（手動）
テスト実行前にユーザーが nvgd を起動する。`playwright.config.ts` には webServer は設定しない方針。

```bash
# ターミナル1
make build
./nvgd -c tests/fixtures/nvgd.conf.yml   # ポート 9280 で待受
```

ベース URL は `http://127.0.0.1:9280`。`playwright.config.ts` の `use.baseURL` を `http://127.0.0.1:9280` に設定する。

### CI の扱い
`.github/workflows/playwright.yml` に Go ビルド & バックグラウンド起動ステップを追加する（手動運用を CI 上で再現）：

```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with: { go-version: '1.24' }
- name: Build nvgd
  run: go build -o nvgd .
- name: Start nvgd
  run: ./nvgd -c tests/fixtures/nvgd.conf.yml &
- name: Wait for server
  run: npx wait-on http://127.0.0.1:9280/version/
```

### Playwright 設定変更（最小限）
| 項目 | 現在 | 変更後 |
|------|------|--------|
| `use.baseURL` | コメントアウト | `http://127.0.0.1:9280` |
| `projects` | chromium / firefox / webkit | chromium / firefox（webkit 削除またはコメントアウト） |
| `use.trace` | `on-first-retry` | 据え置き |
| `expect.timeout` | デフォルト | 5000ms に明示（Mithril 再描画待ち余裕を持たせる） |

---

## 3. ディレクトリ構成

```
tests/
├── PLAN.md                          # 本ドキュメント
├── fixtures/
│   ├── nvgd.conf.yml                # E2E 用最小設定
│   ├── small.txt                    # < 64KiB の編集対象
│   ├── large.bin                    # >= 64KiB（Edit 非表示確認用）
│   ├── sample.csv                   # DuckDB サポート対象
│   ├── sample.json                  # DuckDB サポート対象
│   └── sample.parquet               # DuckDB サポート対象
├── helpers/
│   ├── opfs.ts                      # OPFS クリア・seed・列挙ユーティリティ
│   ├── dialogs.ts                   # confirm/alert ハンドラ
│   └── selectors.ts                 # 共通 locator ヘルパ
├── pages/
│   └── OpfsPage.ts                  # Page Object（UI 操作のラッパ）
└── specs/
    ├── 01-initial-load.spec.ts
    ├── 02-mkdir-navigation.spec.ts
    ├── 03-upload.spec.ts
    ├── 04-editor.spec.ts
    ├── 05-download-url.spec.ts
    ├── 06-selection-delete.spec.ts
    ├── 07-duckdb-integration.spec.ts
    ├── 08-save-as.spec.ts           # Chromium only
    └── 09-history-breadcrumb.spec.ts
```

### Page Object 概要
`OpfsPage.ts` は次の API を提供：
- `goto(path?)` / `gotoRoot()`
- `breadcrumb()` / `currentPath()`
- `rows()` / `rowByName(name)` / `selectRow(name)` / `selectAll()`
- `mkdir(name)` / `uploadFile(localPath, asName?)` / `saveEditor(name, body)` / `clickEdit(name)`
- `downloadFromUrl(url, asName)` / `clickReload()` / `clickDelete()` / `clickDuckDB()`
- `acceptNextDialog()` / `dismissNextDialog()`

---

## 4. テストシナリオ一覧

ID 規則: `<カテゴリ><連番>`。優先度 P0 = 必須、P1 = 推奨、P2 = nice-to-have。

### A. 初期表示・タイトル

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| A1 | P0 | `/opfs/` を開く（hash なし） | `<title>` が `OPFS: (Root)/`、パンくずに `(Root)` 単独表示、grid-header が描画、ファイル行なし |
| A2 | P0 | URL hash `#sub1/sub2/` で開く（事前に seed） | パンくずに `(Root) / sub1 / sub2`、`<title>` がパス反映 |
| A3 | P1 | 存在しない hash パス | `getDirectoryHandle` 失敗 → アラート表示、Root にフォールバック |

### B. ディレクトリ作成と移動

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| B1 | P0 | `mkdir-name` に `foo` 入力 → `Create new directory` | 行に `foo/` が追加、`Type=dir`、入力欄クリア |
| B2 | P0 | mkdir で空文字 | `alert("Need directory name")`、行追加なし |
| B3 | P1 | 同名 mkdir 2 回 | 2 回目で alert、行は重複しない |
| B4 | P1 | 複数 mkdir → 自然順ソート確認 | `file2`, `file10` の順序が正（numeric ソート） |
| B5 | P0 | ディレクトリ名クリック → cd | パンくず増加、`<title>` 更新、`history.pushState` 発火 |
| B6 | P0 | パンくずの親リンククリック | 親に戻る、選択状態は解除 |
| B7 | P0 | 戻る/進む（`page.goBack/goForward`） | popstate でディレクトリ復元、リスト再描画 |
| B8 | P1 | 階層深いディレクトリ作成 → リロード（F5） | hash 経由でカレント復元 |

### C. ファイルアップロード

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| C1 | P0 | `setInputFiles(small.txt)` | `upload-name` に元ファイル名が自動入力、`upload-upload` 有効化 |
| C2 | P0 | アップロード成功 | 行追加、size 一致、`upload-name`/`upload-file` リセット、Upload ボタン無効化 |
| C3 | P0 | 同名アップロード → confirm OK | confirm ダイアログ表示、上書き成功、再アラート `Uploaded ... successfully.` |
| C4 | P1 | 同名アップロード → confirm キャンセル | ファイル内容変化なし、入力欄維持 |
| C5 | P1 | ファイル選択解除（`setInputFiles([])`） | name クリア、Upload ボタン無効化 |
| C6 | P1 | 異なる `upload-name` で保存 | 元名でなく指定名でリスト表示 |

### D. 簡易エディタ

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| D1 | P0 | 新規 `name + body` → Save | ファイル作成、行表示、エディタクリア |
| D2 | P0 | name 空のまま Save | alert、ファイル未作成 |
| D3 | P0 | 既存ファイルを別 body で Save | 内容更新、lastModified 更新 |
| D4 | P0 | `Edit` アクションクリック（< 64KiB） | `editor-name` `editor-edit` に値ロード |
| D5 | P0 | 64KiB 以上ファイルの行 | `Edit` リンクが存在しない、`Save as` のみ表示 |
| D6 | P1 | `Clear` ボタン | name/body 共に空 |
| D7 | P2 | textarea で Tab キー押下 | カーソル位置に `\t` 挿入（`execCommand` 依存） |

### E. URL ダウンロード

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| E1 | P0 | 初期状態 | `download-download` 無効、`download-clear` 無効 |
| E2 | P0 | `http://127.0.0.1:9280/examples/...` + `download-as` 入力 | Download ボタン有効化 |
| E3 | P1 | `ftp://...` を入力 | ボタン無効のまま |
| E4 | P0 | ダウンロード成功 | OPFS にファイル作成、行表示、size 一致 |
| E5 | P0 | 同名ファイル存在時 | confirm ダイアログ、OK で上書き |
| E6 | P1 | 404 を返す URL | alert にエラーメッセージ、ファイル未作成 |
| E7 | P2 | `Clear` ボタン | URL/name 入力リセット、ボタン無効化に戻る |

### F. 選択と一括削除

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| F1 | P0 | チェックボックスを 1 つ ON | Delete/DuckDB ボタン有効化、`selection:change` イベント 1 件 |
| F2 | P0 | 全選択 (`#toggle-selection-all`) | 全行のチェック ON、`toggle.checked=true` |
| F3 | P1 | 一部選択 | toggle が `indeterminate=true` |
| F4 | P0 | Delete → confirm OK | 選択行のみ削除、`Delete` ボタン再度無効化 |
| F5 | P0 | ディレクトリ選択 → Delete | 再帰削除 (`recursive: true`) |
| F6 | P1 | Delete → confirm キャンセル | 削除されない、選択状態維持 |
| F7 | P1 | 削除後の Reload | リスト整合性 |

### G. DuckDB 連携

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| G1 | P0 | `sample.csv` 1 つ選択 → DuckDB | 新タブが開く、URL が `/duckdb/?opfs=...#,...CREATE VIEW opfs0 AS SELECT * FROM 'opfs://...sample.csv'` を含む |
| G2 | P0 | `.csv` + `.json` + `.parquet` 選択 | クエリに `opfs0`/`opfs1`/`opfs2` が連番で含まれる |
| G3 | P1 | 非サポート（`.txt`）混在 | qparams には含まれるが CREATE VIEW は対象外（パスのみ転送） |
| G4 | P1 | ディレクトリ選択 | 再帰的にファイル列挙され、各々が opfsN に展開 |
| G5 | P2 | パスのエスケープ | スペース・記号を含むファイル名で URL が正しくエンコード |

実装メモ: `page.context().waitForEvent('page')` で新タブを捕捉し、URL のみ検証してクローズ。DuckDB 側を実行しない。

### H. Save as（Chromium のみ）

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| H1 | P1 | `Save as` クリック → `showSaveFilePicker` モック | mock 経由でハンドル取得、`createWritable().write()` 呼び出し検証 |
| H2 | P2 | ピッカーキャンセル（AbortError） | アラート表示なし、UI 状態維持 |

実装メモ: Firefox には `showSaveFilePicker` が無いため、テストファイル冒頭で `test.skip(browserName === 'firefox', '...')` を入れる。Chromium でも実ファイルダイアログは出せないので `page.addInitScript` で `window.showSaveFilePicker` を差し替える：

```ts
await page.addInitScript(() => {
  const writes: Uint8Array[] = [];
  (window as any).__writes = writes;
  (window as any).showSaveFilePicker = async () => ({
    name: 'mocked.txt',
    createWritable: async () => ({
      write: async (b: Blob) => writes.push(new Uint8Array(await b.arrayBuffer())),
      close: async () => {},
    }),
  });
});
```

### I. リロードと再描画

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| I1 | P1 | `Reload` ボタン | `renderEntries` 呼び出し、DOM 更新（行数一致） |
| I2 | P1 | mkdir 直後 → リスト即時反映 | `await render()` 後にアサーション通過 |
| I3 | P2 | Mithril 再描画後もイベントリスナが効く | アップロード → reload → 再度選択操作可能 |

### J. 履歴・パンくず詳細

| ID | 優先度 | シナリオ | 主アサーション |
|----|-------|--------|---------------|
| J1 | P1 | 3 階層ディレクトリ作成 → 各階層に cd → goBack 3 回 | popstate 経由で Root まで戻る |
| J2 | P1 | パンくずの中間階層クリック | 該当階層まで pop |
| J3 | P2 | ブラウザリロード（F5） | URL hash 維持、`setCurrPath` でリストア |

合計テストケース概算: **P0 = 26、P1 = 21、P2 = 7、合計 54 ケース**。Chromium / Firefox の 2 プロジェクトで実行するため実行数は約 100（H グループは Chromium のみ）。

---

## 5. テスト基盤の設計詳細

### 5.1 OPFS の隔離戦略

Playwright の各 `test` は新しい `BrowserContext` で起動するため OPFS は空。ただし保険として `beforeEach` で念のためクリア：

```ts
test.beforeEach(async ({ page }) => {
  await page.goto('/opfs/');
  await page.evaluate(async () => {
    const root = await navigator.storage.getDirectory();
    for await (const [name] of (root as any).entries()) {
      await root.removeEntry(name, { recursive: true });
    }
  });
  await page.reload();
});
```

並列実行（`fullyParallel: true`）でも origin 単位で context が分離されるので競合しない。

### 5.2 ダイアログ処理

`confirm` / `alert` は `page.on('dialog', ...)` で受ける。テストごとに振る舞いを切り替えられるよう、ヘルパーで「次の N 回」を制御：

```ts
export async function expectDialog(page: Page, action: 'accept' | 'dismiss', message?: RegExp) {
  return new Promise<Dialog>((resolve) => {
    page.once('dialog', async (d) => {
      if (message) expect(d.message()).toMatch(message);
      await (action === 'accept' ? d.accept() : d.dismiss());
      resolve(d);
    });
  });
}
```

### 5.3 OPFS への seed（前準備）

UI 操作で作る方が忠実だが、深い階層や大量データは `page.evaluate` で直接書き込む：

```ts
await page.evaluate(async ({ name, content }) => {
  const root = await navigator.storage.getDirectory();
  const fh = await root.getFileHandle(name, { create: true });
  const w = await fh.createWritable();
  await w.write(content);
  await w.close();
}, { name: 'seed.txt', content: 'hello' });
```

### 5.4 新規タブ（DuckDB）の捕捉

```ts
const [duckPage] = await Promise.all([
  page.context().waitForEvent('page'),
  page.click('#command-duckdb'),
]);
await expect(duckPage).toHaveURL(/\/duckdb\/\?opfs=.*#,.*CREATE\s+VIEW\s+opfs0/);
await duckPage.close();
```

### 5.5 Mithril 再描画の安定化

Mithril は同期的に DOM を差し替えるが、アクションは async。基本は `await expect(locator).toBeVisible()` で十分。ただし行数アサーションは `await expect(rows).toHaveCount(n)` を使い、ポーリング前提とする。

### 5.6 ファイル fixture の用意

```
tests/fixtures/
├── small.txt        # "hello world" など（< 64KiB）
├── large.bin        # 65536 バイト以上、ランダム or 0 詰め
├── sample.csv       # 3〜5 行、ヘッダ付き
├── sample.json      # JSON 配列、数件
└── sample.parquet   # 数行（生成スクリプトを `tests/fixtures/gen.ts` に置く）
```

`large.bin` は CI でリポジトリに含めず、`global-setup.ts` で生成（`fs.writeFileSync(..., Buffer.alloc(64 * 1024 + 1))`）。

### 5.7 nvgd 用テスト設定 (`tests/fixtures/nvgd.conf.yml`)

最小限の Examples プロトコルが効くように：

```yaml
addr: "127.0.0.1:9280"
access_control_allow_origin: "*"
# OPFS UI に必要なアセットのみで足りるため protocols セクションはほぼデフォルト
```

`examples` プロトコルは E5 で利用（ローカル完結のダウンロード元として）。

---

## 6. ブラウザ別の留意点

| 項目 | Chromium | Firefox |
|------|---------|---------|
| OPFS (`navigator.storage.getDirectory`) | ✓ | ✓ (FF111+) |
| `FileSystemWritableFileStream` | ✓ | ✓ |
| `showSaveFilePicker` | ✓ | ✗ → H グループスキップ |
| `entries()` イテレータ | ✓ | ✓ |
| `execCommand('insertText')` | ✓ | ✓（非標準だが両者対応） |
| Mithril 再描画 | 差なし | 差なし |
| `window.open` 新タブ | ✓ | ✓ |

WebKit は除外（playwright.config.ts の `projects` から削る）。

---

## 7. リスクと注意点

1. **Mithril の onclick が `e.preventDefault()` を呼ぶ** → ディレクトリリンクは標準ナビゲーションしない。`page.click()` で問題なし。
2. **`confirm` 連発** — 削除と上書きで重複しやすい。`page.on('dialog', ...)` の登録漏れに注意し、テスト終了時に `page.removeAllListeners('dialog')`。
3. **同一 origin の OPFS は context 間で共有されうる** — `fullyParallel: true` でも workers が別プロセスなら隔離。同 worker 内の連続テストでは `beforeEach` クリアを必ず実行。
4. **`actCd` 内の `await this.unselectAll()` がある** — cd 直後の選択状態テストは要注意。
5. **`makehash` の `swapchars`** — ` ` ↔ `-`、`;` ↔ `~` の置換が DuckDB URL に効くため、G5 のアサーションは置換後文字列で行う。
6. **`#download-as` 必須** — URL を入れただけでは有効化されない（`as.length > 0` チェックあり）。E2 のテストで両方入れる必要あり。
7. **ファイルロック (`NoModificationAllowedError`)** — 実シナリオでは DuckDB タブが OPFS をロックする。J 系で意図的に発生させるのは難しいため P2 扱い。`alertErr` のメッセージ自体は他経路でカバーされる。
8. **`navigator.storage.getDirectory()` の persisted バケット警告** — 一部環境では `Persistent Storage` 確認を求める。テスト用に `context.grantPermissions(['persistent-storage'])` を呼ぶ。

---

## 8. 実装フェーズ計画（次ステップ）

| フェーズ | 内容 | 見積 |
|--------|------|------|
| 1 | `playwright.config.ts` 更新、`tests/fixtures/`・`tests/helpers/`・`tests/pages/` 雛形作成 | 0.5d |
| 2 | A・B グループ（初期表示 + ディレクトリ）実装 | 0.5d |
| 3 | C・D グループ（アップロード + エディタ）実装 | 0.5d |
| 4 | E・F グループ（URLダウンロード + 削除）実装 | 0.5d |
| 5 | G グループ（DuckDB 連携）実装 | 0.5d |
| 6 | H・I・J グループ（Save as + 再描画 + 履歴）実装 | 0.5d |
| 7 | CI ワークフロー更新（Go ビルド + nvgd 起動ステップ） | 0.25d |
| 8 | 全テスト安定化（flake 取り除き）、レビュー | 0.5d |

合計: 約 3.75 人日

---

## 9. 完了の定義

- 上記 P0 / P1 シナリオが Chromium + Firefox の両方でグリーン
- `npx playwright test` がローカル・CI 双方で実行可能
- `playwright-report/` が CI でアーティファクト保存される
- README または `tests/README.md` に「起動手順」を追記済み

---

## 付録 A: 主要セレクタ一覧

| 要素 | セレクタ |
|------|---------|
| パンくず | `#header` |
| ファイルテーブル | `#main > .directory` |
| 全選択 | `#toggle-selection-all` |
| 行のチェックボックス | `input.selectedFile[name="<filename>"]` |
| Reload | `#command-reload` |
| Delete | `#command-delete` |
| DuckDB | `#command-duckdb` |
| mkdir 名 | `#mkdir-name` |
| mkdir 実行 | `#mkdir-mkdir` |
| アップロード入力 | `#upload-file` |
| アップロード名 | `#upload-name` |
| アップロード実行 | `#upload-upload` |
| エディタ名 | `#editor-name` |
| エディタ本文 | `#editor-edit` |
| エディタ保存 | `#editor-save` |
| エディタクリア | `#editor-clear` |
| ダウンロード URL | `#download-url` |
| 保存名 | `#download-as` |
| ダウンロード実行 | `#download-download` |
| ダウンロードクリア | `#download-clear` |

## 付録 B: 参照ファイル

- UI: `protocol/opfs/assets/index.html`, `protocol/opfs/assets/main.js`
- 仕様: `doc/protocol-opfs.md`
- バックエンド: `protocol/opfs/opfs.go`, `protocol/opfs/opfs_test.go`
- 設定: `playwright.config.ts`, `package.json`, `.github/workflows/playwright.yml`
