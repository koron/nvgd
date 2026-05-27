// spec: specs/opfs-ui-test-plan.md
// TC-02: ディレクトリの作成
// TC-03: ディレクトリ作成 — 名前が空のとき

import { test, expect } from '@playwright/test';
import { gotoOPFS } from './helpers';

test.describe('TC-02: ディレクトリの作成', () => {
  test('新しいディレクトリが一覧に反映される', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#mkdir-name', 'testdir');
    await page.click('#mkdir-mkdir');

    // ディレクトリ行が表示される
    const row = page.locator('.grid-row').filter({ hasText: 'testdir/' });
    await expect(row).toBeVisible();
    await expect(row).toContainText('dir');
    await expect(row).toContainText('(N/A)'); // Size と Modified At

    // 入力欄がクリアされる
    await expect(page.locator('#mkdir-name')).toHaveValue('');
  });
});

test.describe('TC-03: ディレクトリ作成 — 名前が空のとき', () => {
  test('バリデーションアラートが表示されディレクトリは作成されない', async ({ page }) => {
    await gotoOPFS(page);

    const dialogPromise = page.waitForEvent('dialog');
    await page.click('#mkdir-mkdir');

    const dialog = await dialogPromise;
    expect(dialog.message()).toBe('Need directory name');
    await dialog.accept();

    // ファイル行が増えていない
    await expect(page.locator('.grid-row')).toHaveCount(0);
  });
});
