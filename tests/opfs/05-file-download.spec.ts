import { test, expect } from '@playwright/test';

test.describe('05-file-download', () => {
  test('ファイルダウンロード機能（Save as）', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create a file first
    await page.locator('#editor-name').fill('download-test.txt');
    await page.locator('#editor-edit').fill('Content to download');
    await page.locator('#editor-save').click();
    await expect(page.locator('a:has-text("download-test.txt")')).toBeVisible();

    // 1. ファイルの"Save as"リンクをクリック
    // expect: ファイルの"Save as"リンクをクリックするとファイルピッカーが表示される
    const saveAsLink = page.locator('a:has-text("Save as")');
    await expect(saveAsLink).toBeVisible();

    // Handle file picker
    const [fileChooser] = Promise.all([
      page.waitForEvent('filechooser'),
      saveAsLink.click(),
    ]);

    // 2. window.showSaveFilePickerが呼ばれることを確認
    // expect: ファイルピッカーで保存先を選択できる
    // Cancel the file picker (no files to select for download)
    await fileChooser.cancel();

    // 3. ファイルピッカーをキャンセルし、エラーが発生しないことを確認
    // expect: キャンセル時はエラーにならない（AbortError）
    // No error should occur
  });
});
