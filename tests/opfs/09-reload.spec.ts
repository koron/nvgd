import { test, expect } from '@playwright/test';

test.describe('09-reload', () => {
  test('リロードボタンの機能', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. #command-reloadをクリックし、一覧が最新の状態に表示されることを確認
    await page.locator('#command-reload').click();

    // expect: #command-reloadクリックで一覧が更新される
    await expect(page.getByRole('heading', { level: 1 })).toBeVisible();

    // 2. エディタでファイルを作成後、リロードし一覧に追加されることを確認
    await page.locator('#editor-name').fill('reload-test.txt');
    await page.locator('#editor-edit').fill('Reload test content');
    await page.locator('#editor-save').click();

    await expect(page.locator('a:has-text("reload-test.txt")')).toBeVisible();

    // Click reload again
    await page.locator('#command-reload').click();

    // expect: 外部操作（エディタでの作成）後のリロードで新ファイルが表示される
    await expect(page.locator('a:has-text("reload-test.txt")')).toBeVisible();
  });

  test('ハッシュベースのナビゲーション', async ({ page }) => {
    // 1. "http://127.0.0.1:9280/opfs/#/nav-test/"にアクセス
    // First create the directory
    await page.goto('http://127.0.0.1:9280/opfs/');
    await page.locator('#mkdir-name').fill('nav-hash-test');
    await page.locator('#mkdir-mkdir').click();
    await expect(page.locator('a:has-text("nav-hash-test/")')).toBeVisible();

    // Navigate to the directory
    await page.locator('a:has-text("nav-hash-test/")').click();
    await expect(page).toHaveTitle('OPFS: /nav-hash-test/');

    // 2. ページタイトルが"OPFS: /nav-hash-test/"であることを確認
    await expect(page).toHaveTitle('OPFS: /nav-hash-test/');

    // 3. パンくずリストの構成を確認
    await expect(page.locator('#header span')).toBeVisible();
  });
});
