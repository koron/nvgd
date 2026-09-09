# OPFS E2E Test Plan

## Application Overview

OPFS（Origin Private File System）プロトコルのE2Eテスト計画。NVGDサーバー（http://127.0.0.1:9280）の/opfs/パスで提供されるファイルシステムUIのテスト。OPFSはブラウザのnavigator.storage.getDirectory() APIを使用し、Secure Context（localhost/HTTPS）で動作する。テストはChromium、Firefox、WebKitの3ブラウザで並列実行。

## Test Scenarios

### 1. 01-navigation

**Seed:** `tests/seed.spec.ts`

#### 1.1. OPFSページが正常に読み込まれる

**File:** `tests/opfs/01-navigation.spec.ts`

**Steps:**
  1. http://127.0.0.1:9280/opfs/ にアクセス
    - expect: タイトルが"OPFS: /"である
    - expect: h1に"OPFS: Origin Private File System"が表示される
    - expect: パンくずリストに"(Root)"が表示される
  2. ページタイトルを確認
    - expect: ページタイトルが"OPFS: /"である
  3. h1要素の内容を確認
    - expect: h1要素に"OPFS: Origin Private File System"が含まれる
  4. パンくずリストの内容を確認
    - expect: パンくずリストに"(Root)"が表示される

#### 1.2. サブディレクトリのナビゲーション（クリック）

**File:** `tests/opfs/01-navigation.spec.ts`

**Steps:**
  1. ディレクトリ作成: #mkdir-nameに"nav-test"を入力し、#mkdir-mkdirをクリック
    - expect: サブディレクトリが作成されている
  2. "nav-test/"のリンクをクリック
    - expect: ページタイトルが"OPFS: /nav-test/"である
  3. location.hashが"#/nav-test/"であることを確認
    - expect: URLハッシュが"#/nav-test/"である
  4. ページタイトルを確認
    - expect: タイトルが"OPFS: /nav-test/"である
  5. パンくずリストを確認
    - expect: パンくずリストに"(Root)"と"nav-test"が表示される

#### 1.3. パンくずリストでの親ディレクトリ移動

**File:** `tests/opfs/01-navigation.spec.ts`

**Steps:**
  1. サブディレクトリ内にいる状態で、パンくずの"(Root)"をクリック
    - expect: パンくずの"(Root)"リンクをクリックするとルートに戻る
  2. ページタイトルが"OPFS: /"であることを確認
    - expect: ページタイトルが"OPFS: /"である
  3. location.hashを確認
    - expect: URLハッシュが"."または空である

#### 1.4. 操作後のリレンダリング

**File:** `tests/opfs/01-navigation.spec.ts`

**Steps:**
  1. ディレクトリを作成し、一覧に追加されることを確認
    - expect: ディレクトリ作成後、一覧に新ディレクトリが表示される
  2. エディタでファイルを作成し、一覧に追加されることを確認
    - expect: ファイル作成後、一覧に新ファイルが表示される
  3. ファイルを削除し、一覧から消えることを確認
    - expect: ファイル削除後、一覧から削除される

### 2. 02-directory-ops

**Seed:** `tests/seed.spec.ts`

#### 2.1. ディレクトリ作成機能

**File:** `tests/opfs/02-directory-ops.spec.ts`

**Steps:**
  1. #mkdir-nameに"test-dir"を入力し、#mkdir-mkdirをクリック
    - expect: ディレクトリが一覧に追加される
  2. #mkdir-nameの値が空であることを確認
    - expect: 入力フィールドがクリアされる
  3. 一覧でtype列が"dir"であることを確認
    - expect: 作成されたディレクトリの型が"dir"である
  4. 一覧でsizeとmodifiedAtが"(N/A)"であることを確認
    - expect: サイズとModified Atが"(N/A)"である

#### 2.2. ディレクトリ名のバリデーション（空入力）

**File:** `tests/opfs/02-directory-ops.spec.ts`

**Steps:**
  1. #mkdir-nameを空にして#mkdir-mkdirをクリック
    - expect: 空の入力でアラート"Need directory name"が表示される
  2. 一覧に新しいディレクトリが追加されないことを確認
    - expect: ディレクトリが作成されない

#### 2.3. ディレクトリ名に特殊文字

**File:** `tests/opfs/02-directory-ops.spec.ts`

**Steps:**
  1. #mkdir-nameに"dir with spaces"を入力し作成
    - expect: スペースを含むディレクトリ名が作成される
  2. #mkdir-nameに"my-dir_name"を入力し作成
    - expect: ハイフンを含むディレクトリ名が作成される
  3. #mkdir-nameに"dir.special;name"を入力し作成
    - expect: 特殊文字を含むディレクトリ名が作成される

#### 2.4. アップロード名の変更

**File:** `tests/opfs/02-directory-ops.spec.ts`

**Steps:**
  1. ファイルを選択し、#upload-nameの値を確認
    - expect: ファイル選択後、#upload-nameにファイル名が自動入力される
  2. #upload-nameの値を別の名前に変更
    - expect: #upload-nameの値を編集できる

### 3. 03-file-upload

**Seed:** `tests/seed.spec.ts`

#### 3.1. ファイルアップロード機能

**File:** `tests/opfs/03-file-upload.spec.ts`

**Steps:**
  1. ファイルを選択（evaluateでDataTransferを使用）
    - expect: ファイル選択後、#upload-nameにファイル名が自動入力される
  2. #upload-uploadのdisabled属性を確認
    - expect: ファイル選択後、#upload-uploadボタンが有効になる
  3. #upload-uploadをクリックし、一覧にファイルが追加されることを確認
    - expect: #upload-uploadクリック後、ファイルがOPFSに保存される
  4. alertダイアログのメッセージを確認
    - expect: アップロード成功アラートが表示される
  5. 一覧でsize列がファイルサイズと一致することを確認
    - expect: ファイルのサイズが正しい

#### 3.2. 同じ名前のファイルが存在する場合のアップロード（上書き確認）

**File:** `tests/opfs/03-file-upload.spec.ts`

**Steps:**
  1. 既存ファイルと同じ名前でファイルアップロードを試みる
    - expect: 既存ファイルと同じ名前でアップロード時、確認ダイアログが表示される
  2. confirmダイアログでacceptし、ファイルが上書きされることを確認
    - expect: 確認ダイアログでOKを選択すると上書きされる
  3. confirmダイアログでdismissし、ファイルが上書きされないことを確認
    - expect: 確認ダイアログでキャンセルすると上書きされない

#### 3.3. 空のファイル名のバリデーション（エディタ）

**File:** `tests/opfs/03-file-upload.spec.ts`

**Steps:**
  1. #editor-nameを空にして#editor-saveをクリック
    - expect: エディタで空のファイル名で保存時、エラーが発生する

### 4. 04-file-edit

**Seed:** `tests/seed.spec.ts`

#### 4.1. ファイル編集機能（64KiB未満）

**File:** `tests/opfs/04-file-edit.spec.ts`

**Steps:**
  1. 小ファイルを作成し、一覧で"Edit"リンクが表示されることを確認
    - expect: 64KiB未満のファイルに"Edit"アクションが表示される
  2. "Edit"リンクをクリックし、#editor-nameと#editor-editの値を確認
    - expect: "Edit"クリックでエディタにファイル名と内容が読み込まれる
  3. #editor-editの値を変更し、#editor-saveをクリック
    - expect: エディタで内容を変更して保存できる
  4. ファイルを再度Editで読み込み、更新された内容を確認
    - expect: 保存後、ファイルの内容が更新される

#### 4.2. エディタのClear機能

**File:** `tests/opfs/04-file-edit.spec.ts`

**Steps:**
  1. エディタに値を入力後、#editor-clearをクリック
    - expect: エディタに値が入力された状態で#editor-clearをクリックすると、両フィールドが空になる
  2. #editor-name.inputValue()が空であることを確認
    - expect: #editor-nameが空になる
  3. #editor-edit.inputValue()が空であることを確認
    - expect: #editor-editが空になる

#### 4.3. 64KiB超のファイルにはEditが表示されない

**File:** `tests/opfs/04-file-edit.spec.ts`

**Steps:**
  1. 64KiB超のファイルを作成し、一覧で"Edit"リンクが表示されないことを確認
    - expect: 64KiB超のファイルには"Edit"リンクが表示されない
  2. 一覧で"Save as"リンクが表示されることを確認
    - expect: 64KiB超のファイルには"Save as"リンクが表示される

### 5. 05-file-download

**Seed:** `tests/seed.spec.ts`

#### 5.1. ファイルダウンロード機能（Save as）

**File:** `tests/opfs/05-file-download.spec.ts`

**Steps:**
  1. ファイルの"Save as"リンクをクリック
    - expect: ファイルの"Save as"リンクをクリックするとファイルピッカーが表示される
  2. window.showSaveFilePickerが呼ばれることを確認
    - expect: ファイルピッカーで保存先を選択できる
  3. ファイルピッカーをキャンセルし、エラーが発生しないことを確認
    - expect: キャンセル時はエラーにならない（AbortError）

### 6. 06-file-delete

**Seed:** `tests/seed.spec.ts`

#### 6.1. 単体ファイル削除

**File:** `tests/opfs/06-file-delete.spec.ts`

**Steps:**
  1. ファイルのチェックボックスをチェックし、#command-deleteのdisabled状態を確認
    - expect: ファイルのチェックボックスを選択後、Deleteボタンが有効になる
  2. #command-deleteをクリックし、confirmダイアログが表示されることを確認
    - expect: Deleteクリックで確認ダイアログが表示される
  3. confirmダイアログでacceptし、一覧からファイルが消えることを確認
    - expect: 確認でOKを選択するとファイルが削除される
  4. confirmダイアログでdismissし、一覧に残ることを確認
    - expect: 確認でキャンセルするとファイルが残る

#### 6.2. 複数ファイルの削除

**File:** `tests/opfs/06-file-delete.spec.ts`

**Steps:**
  1. 複数のファイルをチェック
    - expect: 複数のファイルを選択後、Deleteボタンが有効になる
  2. confirmダイアログのメッセージに全ファイル名が含まれることを確認
    - expect: 確認ダイアログに選択したファイル名が全て表示される
  3. confirmでacceptし、全ファイルが消えることを確認
    - expect: 確認でOKを選択すると全ファイルが削除される

#### 6.3. ディレクトリ選択時の再帰削除

**File:** `tests/opfs/06-file-delete.spec.ts`

**Steps:**
  1. サブディレクトリ（中にファイルあり）を選択し削除
    - expect: ディレクトリを選択してDeleteすると、内容を含むディレクトリが削除される
  2. ディレクトリのcheckbox nameが"dir/"形式であることを確認
    - expect: ディレクトリ名は"dir/"形式で表示される

### 7. 07-multi-select

**Seed:** `tests/seed.spec.ts`

#### 7.1. 全選択チェックボックス

**File:** `tests/opfs/07-multi-select.spec.ts`

**Steps:**
  1. #toggle-selection-allをクリックし、全てのcheckboxがcheckedであることを確認
    - expect: 全選択チェックボックスをクリックすると全アイテムが選択される
  2. #toggle-selection-allを再度クリックし、全てのcheckboxがuncheckedであることを確認
    - expect: 再度クリックすると全選択解除される
  3. 一部のアイテムをチェックし、#toggle-selection-allのindeterminate状態を確認
    - expect: 一部のアイテムが選択されている状態ではindeterminateになる

#### 7.2. 選択状態とボタンの連動

**File:** `tests/opfs/07-multi-select.spec.ts`

**Steps:**
  1. アイテムをチェックし、#command-deleteのdisabledがfalseであることを確認
    - expect: アイテムを選択するとDeleteボタンが有効になる
  2. アイテムをチェックし、#command-duckdbのdisabledがfalseであることを確認
    - expect: アイテムを選択するとDuckDBボタンが有効になる
  3. 全選択解除し、両ボタンのdisabledがtrueであることを確認
    - expect: 全選択解除すると両ボタンが無効になる

#### 7.3. 全選択解除後の自動アンチェック

**File:** `tests/opfs/07-multi-select.spec.ts`

**Steps:**
  1. サブディレクトリに移動し、チェックボックスが全てuncheckedであることを確認
    - expect: ディレクトリ移動時に選択が解除される

### 8. 08-url-download

**Seed:** `tests/seed.spec.ts`

#### 8.1. URLダウンロード機能

**File:** `tests/opfs/08-url-download.spec.ts`

**Steps:**
  1. #download-urlに"http://127.0.0.1:9280/help/"を入力
    - expect: #download-urlに有効なURL（http:またはhttps:）を入力すると#download-downloadが有効になる
  2. #download-asに名前を入力し、#download-downloadのdisabledがfalseであることを確認
    - expect: #download-asにファイル名を入力すると#download-downloadが有効になる
  3. #download-downloadをクリックし、一覧にファイルが追加されることを確認
    - expect: #download-downloadクリックでファイルがOPFSに保存される
  4. 既存ファイルと同じ名前でURLダウンロードを試みる
    - expect: 既存ファイルの上書き確認ダイアログが表示される
  5. confirmでacceptし、ファイルが上書きされることを確認
    - expect: 上書き確認でOKを選択すると上書きされる

#### 8.2. URLのバリデーション

**File:** `tests/opfs/08-url-download.spec.ts`

**Steps:**
  1. "ftp://..."を入力し、#download-downloadのdisabledがtrueであることを確認
    - expect: http:またはhttps:で始まらないURLでは#download-downloadが無効
  2. #download-urlを空にして、#download-downloadのdisabledがtrueであることを確認
    - expect: #download-urlを空にすると#download-downloadが無効になる
  3. #download-asを空にして、#download-downloadのdisabledがtrueであることを確認
    - expect: #download-asを空にすると#download-downloadが無効になる
  4. #download-clearをクリックし、両フィールドが空であることを確認
    - expect: #download-clearクリックで両フィールドが空になる

### 9. 09-reload

**Seed:** `tests/seed.spec.ts`

#### 9.1. リロードボタンの機能

**File:** `tests/opfs/09-reload.spec.ts`

**Steps:**
  1. #command-reloadをクリックし、一覧が最新の状態に表示されることを確認
    - expect: #command-reloadクリックで一覧が更新される
  2. エディタでファイルを作成後、リロードし一覧に追加されることを確認
    - expect: 外部操作（エディタでの作成）後のリロードで新ファイルが表示される

#### 9.2. ハッシュベースのナビゲーション

**File:** `tests/opfs/09-reload.spec.ts`

**Steps:**
  1. "http://127.0.0.1:9280/opfs/#/nav-test/"にアクセス
    - expect: URLハッシュでサブディレクトリに直接アクセスできる
  2. ページタイトルが"OPFS: /nav-test/"であることを確認
    - expect: ハッシュでアクセスしたディレクトリのタイトルが更新される
  3. パンくずリストの構成を確認
    - expect: パンくずリストがハッシュのパスに一致する
