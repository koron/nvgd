// spec: specs/opfs-ui-test-plan.md
// TC-09: シンプルエディタでファイルを新規作成
// TC-10: シンプルエディタ — ファイル名が空のとき
// TC-11: シンプルエディタ — Tab キーでタブ文字を挿入
// TC-12: ファイルの編集（Edit アクション）
// TC-13: エディタのクリアボタン

import { test, expect } from '@playwright/test';
import { gotoOPFS, createOPFSFile, reloadListing } from './helpers';

const SKIP_WEBKIT = 'OPFS は WebKit の HTTP では利用不可（セキュアコンテキスト外）';

test.describe('TC-09: シンプルエディタでファイルを新規作成', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('エディタからファイルが正常に作成される', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#editor-name', 'note.txt');
    await page.fill('#editor-edit', 'Hello, OPFS!');
    await page.click('#editor-save');

    await expect(page.locator('.grid-row').filter({ hasText: 'note.txt' })).toBeVisible();

    // 入力欄がクリアされる
    await expect(page.locator('#editor-name')).toHaveValue('');
    await expect(page.locator('#editor-edit')).toHaveValue('');
  });
});

test.describe('TC-10: シンプルエディタ — ファイル名が空のとき', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('バリデーションアラートが表示されファイルは作成されない', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#editor-edit', 'some content');

    // page.once でクリック中に発火するアラートをハンドルする（waitForEvent+await click のデッドロック回避）
    let capturedMessage = '';
    page.once('dialog', async (dialog) => {
      capturedMessage = dialog.message();
      await dialog.accept();
    });
    await page.click('#editor-save');

    expect(capturedMessage).toBe('Need file name');
    await expect(page.locator('.grid-row')).toHaveCount(0);
  });
});

test.describe('TC-11: シンプルエディタ — Tab キーでタブ文字を挿入', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('Tab キーでタブ文字が挿入されフォーカスが移動しない', async ({ page }) => {
    await gotoOPFS(page);

    await page.locator('#editor-edit').focus();
    await page.keyboard.press('Tab');

    const value = await page.locator('#editor-edit').inputValue();
    expect(value).toContain('\t');

    await expect(page.locator('#editor-edit')).toBeFocused();
  });
});

test.describe('TC-12: ファイルの編集（Edit アクション）', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('Edit リンクでファイル内容がエディタにロードされ上書き保存できる', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'note.txt', 'Hello, OPFS!');
    await reloadListing(page);

    await page.locator('.grid-row').filter({ hasText: 'note.txt' }).getByText('Edit').click();

    await expect(page.locator('#editor-name')).toHaveValue('note.txt');
    await expect(page.locator('#editor-edit')).toHaveValue('Hello, OPFS!');

    await page.fill('#editor-edit', 'Updated content');
    await page.click('#editor-save');

    await expect(page.locator('.grid-row')).toHaveCount(1);
    await expect(page.locator('.grid-row').filter({ hasText: 'note.txt' })).toBeVisible();
  });
});

test.describe('TC-13: エディタのクリアボタン', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('Clear ボタンで Name 欄とテキストエリアがリセットされる', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#editor-name', 'temp.txt');
    await page.fill('#editor-edit', 'some text');
    await page.click('#editor-clear');

    await expect(page.locator('#editor-name')).toHaveValue('');
    await expect(page.locator('#editor-edit')).toHaveValue('');
  });
});
