// spec: specs/opfs-ui-test-plan.md
// TC-16: ファイル・ディレクトリの削除
// TC-22: ディレクトリの再帰削除

import { test, expect } from '@playwright/test';
import { gotoOPFS, createOPFSFile, createOPFSFileInDir, reloadListing } from './helpers';

test.describe('TC-16: ファイル・ディレクトリの削除', () => {
  test('キャンセルすると削除されない', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'target.txt', 'delete me');
    await reloadListing(page);

    await page.locator('input.selectedFile[name="target.txt"]').check();
    const dialog1 = page.waitForEvent('dialog');
    await page.click('#command-delete');
    const dlg1 = await dialog1;
    expect(dlg1.message()).toContain('target.txt');
    await dlg1.dismiss(); // キャンセル

    // ファイルが残っている
    await expect(page.locator('.grid-row').filter({ hasText: 'target.txt' })).toBeVisible();
  });

  test('OK すると選択したファイルが削除される', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'target.txt', 'delete me');
    await reloadListing(page);

    await page.locator('input.selectedFile[name="target.txt"]').check();
    const dialogPromise = page.waitForEvent('dialog');
    await page.click('#command-delete');
    await (await dialogPromise).accept();

    // 行が消える
    await expect(page.locator('.grid-row')).toHaveCount(0);
  });
});

test.describe('TC-22: ディレクトリの再帰削除', () => {
  test('ファイルを含むディレクトリをまとめて削除できる', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFileInDir(page, 'testdir', 'child.txt', 'content');
    await reloadListing(page);

    // testdir/ のチェックボックスは name="testdir/" (スラッシュ付き)
    await page.locator('input.selectedFile[name="testdir/"]').check();
    const dialogPromise = page.waitForEvent('dialog');
    await page.click('#command-delete');
    const dlg = await dialogPromise;
    expect(dlg.message()).toContain('testdir/');
    await dlg.accept();

    // ディレクトリが消える
    await expect(page.locator('.grid-row')).toHaveCount(0);
  });
});
