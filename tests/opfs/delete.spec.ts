// spec: specs/opfs-ui-test-plan.md
// TC-16: ファイル・ディレクトリの削除
// TC-22: ディレクトリの再帰削除

import { test, expect } from '@playwright/test';
import { gotoOPFS, createOPFSFile, createOPFSFileInDir, reloadListing } from './helpers';

const SKIP_WEBKIT = 'OPFS は WebKit の HTTP では利用不可（セキュアコンテキスト外）';

test.describe('TC-16: ファイル・ディレクトリの削除', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('キャンセルすると削除されない', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'target.txt', 'delete me');
    await reloadListing(page);

    await page.locator('input.selectedFile[name="target.txt"]').check();

    // page.once でクリック中に発火する confirm をハンドルする（waitForEvent+await click のデッドロック回避）
    page.once('dialog', (dialog) => dialog.dismiss());
    await page.click('#command-delete');

    await expect(page.locator('.grid-row').filter({ hasText: 'target.txt' })).toBeVisible();
  });

  test('OK すると選択したファイルが削除される', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'target.txt', 'delete me');
    await reloadListing(page);

    await page.locator('input.selectedFile[name="target.txt"]').check();

    let capturedMessage = '';
    page.once('dialog', async (dialog) => {
      capturedMessage = dialog.message();
      await dialog.accept();
    });
    await page.click('#command-delete');

    expect(capturedMessage).toContain('target.txt');
    await expect(page.locator('.grid-row')).toHaveCount(0);
  });
});

test.describe('TC-22: ディレクトリの再帰削除', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('ファイルを含むディレクトリをまとめて削除できる', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFileInDir(page, 'testdir', 'child.txt', 'content');
    await reloadListing(page);

    // testdir/ のチェックボックスは name="testdir/" (スラッシュ付き)
    await page.locator('input.selectedFile[name="testdir/"]').check();

    let capturedMessage = '';
    page.once('dialog', async (dialog) => {
      capturedMessage = dialog.message();
      await dialog.accept();
    });
    await page.click('#command-delete');

    expect(capturedMessage).toContain('testdir/');
    await expect(page.locator('.grid-row')).toHaveCount(0);
  });
});
