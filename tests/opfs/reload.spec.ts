// spec: specs/opfs-ui-test-plan.md
// TC-17: Reload ボタン

import { test, expect } from '@playwright/test';
import { gotoOPFS, createOPFSFile } from './helpers';

test.describe('TC-17: Reload ボタン', () => {
  test('Reload ボタンでエラーなく一覧が再描画される', async ({ page }) => {
    await gotoOPFS(page);

    // 初期状態で Reload — エラーなし・行数が変わらない
    await page.click('#command-reload');
    await expect(page.locator('.grid-row')).toHaveCount(0);

    // OPFS に直接ファイルを作成してから Reload → 反映される
    await createOPFSFile(page, 'new.txt', 'reload test');
    await page.click('#command-reload');
    await expect(page.locator('.grid-row').filter({ hasText: 'new.txt' })).toBeVisible();
  });
});
