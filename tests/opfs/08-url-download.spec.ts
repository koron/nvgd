import { test, expect } from '@playwright/test';

test.describe('08-url-download', () => {
  test('URLダウンロード機能', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. #download-urlに"http://127.0.0.1:9280/help/"を入力
    await page.locator('#download-url').fill('http://127.0.0.1:9280/help/');

    // expect: #download-urlに有効なURL（http:またはhttps:）を入力すると#download-downloadが有効になる
    await expect(page.locator('#download-download')).not.toBeDisabled();

    // 2. #download-asに名前を入力し、#download-downloadのdisabledがfalseであることを確認
    await page.locator('#download-as').fill('help-download.txt');

    // expect: #download-asにファイル名を入力すると#download-downloadが有効になる
    await expect(page.locator('#download-download')).not.toBeDisabled();

    // 3. #download-downloadをクリックし、一覧にファイルが追加されることを確認
    // Handle alert dialog
    page.once('dialog', dialog => {
      expect(dialog.message()).toContain('help-download.txt');
      dialog.accept();
    });

    await page.locator('#download-download').click();

    // expect: #download-downloadクリックでファイルがOPFSに保存される
    await expect(page.locator('a:has-text("help-download.txt")')).toBeVisible();

    // 4. 既存ファイルと同じ名前でURLダウンロードを試みる
    await page.locator('#download-url').fill('http://127.0.0.1:9280/help/');
    await page.locator('#download-as').fill('help-download.txt');

    // expect: 既存ファイルの上書き確認ダイアログが表示される

    // 5. confirmでacceptし、ファイルが上書きされることを確認
    page.once('dialog', async dialog => {
      expect(dialog.message()).toContain('help-download.txt');
      await dialog.accept();
    });

    await page.locator('#download-download').click();

    // expect: 上書き確認でOKを選択すると上書きされる
    await expect(page.locator('a:has-text("help-download.txt")')).toBeVisible();
  });

  test('URLのバリデーション', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. "ftp://..."を入力し、#download-downloadのdisabledがtrueであることを確認
    await page.locator('#download-url').fill('ftp://example.com/file.txt');

    // expect: http:またはhttps:で始まらないURLでは#download-downloadが無効
    await expect(page.locator('#download-download')).toBeDisabled();

    // 2. #download-urlを空にして、#download-downloadのdisabledがtrueであることを確認
    await page.locator('#download-url').clear();

    // expect: #download-urlを空にすると#download-downloadが無効になる
    await expect(page.locator('#download-download')).toBeDisabled();

    // 3. #download-asを空にして、#download-downloadのdisabledがtrueであることを確認
    await page.locator('#download-url').fill('http://127.0.0.1:9280/help/');
    await page.locator('#download-as').fill('test.txt');
    await expect(page.locator('#download-download')).not.toBeDisabled();

    await page.locator('#download-as').clear();

    // expect: #download-asを空にすると#download-downloadが無効になる
    await expect(page.locator('#download-download')).toBeDisabled();

    // 4. #download-clearをクリックし、両フィールドが空であることを確認
    await page.locator('#download-clear').click();

    // expect: #download-clearクリックで両フィールドが空になる
    await expect(page.locator('#download-url')).toHaveValue('');
    await expect(page.locator('#download-as')).toHaveValue('');
  });
});
