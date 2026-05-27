// spec: specs/opfs-ui-test-plan.md
// TC-18: URLからのファイルダウンロード
// TC-19: ダウンロード — 無効な URL

import { test, expect } from '@playwright/test';
import { gotoOPFS } from './helpers';

const SKIP_WEBKIT = 'OPFS は WebKit の HTTP では利用不可（セキュアコンテキスト外）';

const MOCK_URL = 'http://127.0.0.1:9280/__test-download__.txt';
const MOCK_BODY = 'Downloaded file content';

test.describe('TC-18: URL からのファイルダウンロード', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('URL を指定してファイルを OPFS に保存できる', async ({ page }) => {
    // fetch リクエストをインターセプトしてモックレスポンスを返す
    await page.route(MOCK_URL, (route) =>
      route.fulfill({
        status: 200,
        contentType: 'text/plain',
        body: MOCK_BODY,
      }),
    );

    await gotoOPFS(page);

    // URL が空の段階では Download ボタンは無効
    await expect(page.locator('#download-download')).toBeDisabled();

    // URL と保存名を入力
    await page.fill('#download-url', MOCK_URL);
    await page.fill('#download-as', 'downloaded.txt');

    // URL と名前が揃うと Download ボタンが有効になる
    await expect(page.locator('#download-download')).toBeEnabled();

    // ダウンロード実行
    await page.click('#download-download');

    // ファイルが一覧に現れる
    await expect(page.locator('.grid-row').filter({ hasText: 'downloaded.txt' })).toBeVisible();

    // Clear ボタンで入力欄がリセットされる
    await page.click('#download-clear');
    await expect(page.locator('#download-url')).toHaveValue('');
    await expect(page.locator('#download-as')).toHaveValue('');
  });
});

test.describe('TC-19: ダウンロード — 無効な URL', () => {
  test.skip(({ browserName }) => browserName === 'webkit', SKIP_WEBKIT);

  test('http/https 以外のプロトコルでは Download ボタンが無効のまま', async ({ page }) => {
    await gotoOPFS(page);

    // ftp スキーム → 無効
    await page.fill('#download-url', 'ftp://example.com/file.txt');
    await page.fill('#download-as', 'test.txt');
    await expect(page.locator('#download-download')).toBeDisabled();

    // https スキームに変更 → 有効になる
    await page.fill('#download-url', 'https://example.com/file.txt');
    await expect(page.locator('#download-download')).toBeEnabled();
  });

  test('URL が空の場合は Download ボタンが無効のまま', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#download-as', 'test.txt');
    await expect(page.locator('#download-download')).toBeDisabled();
  });

  test('保存名が空の場合は Download ボタンが無効のまま', async ({ page }) => {
    await gotoOPFS(page);

    await page.fill('#download-url', 'https://example.com/file.txt');
    await expect(page.locator('#download-download')).toBeDisabled();
  });
});
