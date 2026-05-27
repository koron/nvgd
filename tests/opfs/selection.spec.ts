// spec: specs/opfs-ui-test-plan.md
// TC-14: ファイルの選択と Delete/DuckDB ボタンの有効化
// TC-15: 全選択チェックボックス

import { test, expect } from '@playwright/test';
import { gotoOPFS, createOPFSFile, reloadListing } from './helpers';

const SKIP_WEBKIT = 'OPFS は WebKit の HTTP では利用不可（セキュアコンテキスト外）';

test.describe('TC-14: 選択状態に応じたボタン有効/無効の切り替え', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('ファイルを選択すると Delete・DuckDB ボタンが有効になる', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'file1.txt', 'a');
    await reloadListing(page);

    // 初期状態: 両ボタンが無効
    await expect(page.locator('#command-delete')).toBeDisabled();
    await expect(page.locator('#command-duckdb')).toBeDisabled();

    // チェックボックスをオン
    await page.locator('input.selectedFile[name="file1.txt"]').check();
    await expect(page.locator('#command-delete')).toBeEnabled();
    await expect(page.locator('#command-duckdb')).toBeEnabled();

    // チェックボックスをオフ → 再び無効
    await page.locator('input.selectedFile[name="file1.txt"]').uncheck();
    await expect(page.locator('#command-delete')).toBeDisabled();
    await expect(page.locator('#command-duckdb')).toBeDisabled();
  });
});

test.describe('TC-15: 全選択チェックボックス', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('全選択・全解除・不定状態が正しく動作する', async ({ page }) => {
    await gotoOPFS(page);
    await createOPFSFile(page, 'a.txt', 'a');
    await createOPFSFile(page, 'b.txt', 'b');
    await reloadListing(page);

    const toggleAll = page.locator('#toggle-selection-all');

    // 全選択チェックボックスをオン
    await toggleAll.check();
    await expect(page.locator('input.selectedFile:checked')).toHaveCount(2);
    await expect(page.locator('#command-delete')).toBeEnabled();

    // 全選択チェックボックスをオフ → 全解除
    await toggleAll.uncheck();
    await expect(page.locator('input.selectedFile:checked')).toHaveCount(0);
    await expect(page.locator('#command-delete')).toBeDisabled();

    // 1件だけ選択 → 不定状態（indeterminate）になる
    await page.locator('input.selectedFile').first().check();
    const isIndeterminate = await toggleAll.evaluate(
      (el) => (el as HTMLInputElement).indeterminate,
    );
    expect(isIndeterminate).toBe(true);
  });
});
