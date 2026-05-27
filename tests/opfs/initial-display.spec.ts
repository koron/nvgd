// spec: specs/opfs-ui-test-plan.md
// TC-01: 初期表示

import { test, expect } from '@playwright/test';
import { gotoOPFS } from './helpers';

test.describe('TC-01: 初期表示', () => {
  test('空の OPFS ルートが正しく表示される', async ({ page }) => {
    await gotoOPFS(page);

    // タイトル（render() は d.dir.name を使うので root は空文字 → 'OPFS: /'）
    await expect(page).toHaveTitle('OPFS: /');

    // パンくずリストに (Root) が表示される
    await expect(page.locator('#header')).toContainText('(Root)');

    // テーブルヘッダーが揃っている
    const table = page.locator('#main .directory');
    await expect(table).toContainText('Name');
    await expect(table).toContainText('Type');
    await expect(table).toContainText('Size');
    await expect(table).toContainText('Modified At');
    await expect(table).toContainText('Actions');

    // ファイル行が 0 件
    await expect(page.locator('.grid-row')).toHaveCount(0);

    // ボタン状態
    await expect(page.locator('#command-reload')).toBeEnabled();
    await expect(page.locator('#command-delete')).toBeDisabled();
    await expect(page.locator('#command-duckdb')).toBeDisabled();
  });
});
