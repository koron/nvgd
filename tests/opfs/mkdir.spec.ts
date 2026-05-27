// spec: specs/opfs-ui-test-plan.md
// TC-02: ディレクトリの作成
// TC-03: ディレクトリ作成 — 名前が空のとき

import { test, expect } from '@playwright/test';
import { gotoOPFS } from './helpers';

const SKIP_WEBKIT = 'OPFS は WebKit の HTTP では利用不可（セキュアコンテキスト外）';

test.describe('TC-02: ディレクトリの作成', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('新しいディレクトリが一覧に反映される', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#mkdir-name', 'testdir');
    await page.click('#mkdir-mkdir');

    // ディレクトリ行が表示される
    const row = page.locator('.grid-row').filter({ hasText: 'testdir/' });
    await expect(row).toBeVisible();
    await expect(row).toContainText('dir');
    await expect(row).toContainText('(N/A)');

    // 入力欄がクリアされる
    await expect(page.locator('#mkdir-name')).toHaveValue('');
  });
});

test.describe('TC-03: ディレクトリ作成 — 名前が空のとき', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('バリデーションアラートが表示されディレクトリは作成されない', async ({ page }) => {
    await gotoOPFS(page);

    // page.once でクリック中に発火するアラートをハンドルする（waitForEvent+await click のデッドロック回避）
    let capturedMessage = '';
    page.once('dialog', async (dialog) => {
      capturedMessage = dialog.message();
      await dialog.accept();
    });
    await page.click('#mkdir-mkdir');

    expect(capturedMessage).toBe('Need directory name');
    await expect(page.locator('.grid-row')).toHaveCount(0);
  });
});
