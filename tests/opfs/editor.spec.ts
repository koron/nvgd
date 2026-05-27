// spec: specs/opfs-ui-test-plan.md
// TC-09: シンプルエディタでファイルを新規作成
// TC-10: シンプルエディタ — ファイル名が空のとき
// TC-11: シンプルエディタ — Tab キーでタブ文字を挿入
// TC-12: ファイルの編集（Edit アクション）
// TC-13: エディタのクリアボタン

import { test, expect } from '@playwright/test';
import { gotoOPFS, createOPFSFile, reloadListing } from './helpers';

test.describe('TC-09: シンプルエディタでファイルを新規作成', () => {
  test('エディタからファイルが正常に作成される', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#editor-name', 'note.txt');
    await page.fill('#editor-edit', 'Hello, OPFS!');
    await page.click('#editor-save');

    // ファイル一覧に note.txt が現れる
    await expect(page.locator('.grid-row').filter({ hasText: 'note.txt' })).toBeVisible();

    // 入力欄がクリアされる
    await expect(page.locator('#editor-name')).toHaveValue('');
    await expect(page.locator('#editor-edit')).toHaveValue('');
  });
});

test.describe('TC-10: シンプルエディタ — ファイル名が空のとき', () => {
  test('バリデーションアラートが表示されファイルは作成されない', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#editor-edit', 'some content');

    const dialogPromise = page.waitForEvent('dialog');
    await page.click('#editor-save');
    const dialog = await dialogPromise;
    expect(dialog.message()).toBe('Need file name');
    await dialog.accept();

    await expect(page.locator('.grid-row')).toHaveCount(0);
  });
});

test.describe('TC-11: シンプルエディタ — Tab キーでタブ文字を挿入', () => {
  test('Tab キーでタブ文字が挿入されフォーカスが移動しない', async ({ page }) => {
    await gotoOPFS(page);

    await page.locator('#editor-edit').focus();
    await page.keyboard.press('Tab');

    // タブ文字が挿入される
    const value = await page.locator('#editor-edit').inputValue();
    expect(value).toContain('\t');

    // フォーカスがエディタに残っている
    await expect(page.locator('#editor-edit')).toBeFocused();
  });
});

test.describe('TC-12: ファイルの編集（Edit アクション）', () => {
  test('Edit リンクでファイル内容がエディタにロードされ上書き保存できる', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'note.txt', 'Hello, OPFS!');
    await reloadListing(page);

    // Edit リンクをクリック
    await page.locator('.grid-row').filter({ hasText: 'note.txt' }).getByText('Edit').click();

    // エディタに読み込まれる
    await expect(page.locator('#editor-name')).toHaveValue('note.txt');
    await expect(page.locator('#editor-edit')).toHaveValue('Hello, OPFS!');

    // 内容を書き換えて保存
    await page.fill('#editor-edit', 'Updated content');
    await page.click('#editor-save');

    // 一覧に note.txt が残っている（重複せず）
    await expect(page.locator('.grid-row')).toHaveCount(1);
    await expect(page.locator('.grid-row').filter({ hasText: 'note.txt' })).toBeVisible();
  });
});

test.describe('TC-13: エディタのクリアボタン', () => {
  test('Clear ボタンで Name 欄とテキストエリアがリセットされる', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#editor-name', 'temp.txt');
    await page.fill('#editor-edit', 'some text');
    await page.click('#editor-clear');

    await expect(page.locator('#editor-name')).toHaveValue('');
    await expect(page.locator('#editor-edit')).toHaveValue('');
  });
});
