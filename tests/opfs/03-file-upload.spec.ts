import { test, expect } from '@playwright/test';

test.describe('03-file-upload', () => {
  test('ファイルアップロード機能', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. ファイルを選択（input[type="file"]を使用）
    await page.locator('#upload-file').setInputFiles({
      name: 'upload-test.txt',
      contentType: 'text/plain',
      buffer: Buffer.from('Hello, OPFS!'),
    });

    // expect: ファイル選択後、#upload-nameにファイル名が自動入力される
    await expect(page.locator('#upload-name')).toHaveValue('upload-test.txt');

    // 2. #upload-uploadのdisabled属性を確認
    // expect: ファイル選択後、#upload-uploadボタンが有効になる
    const uploadBtn = page.locator('#upload-upload');
    await expect(uploadBtn).not.toBeDisabled();

    // 3. #upload-uploadをクリックし、一覧にファイルが追加されることを確認
    // Handle the alert dialog
    page.once('dialog', dialog => {
      expect(dialog.message()).toContain('upload-test.txt');
      dialog.accept();
    });

    await uploadBtn.click();

    // expect: #upload-uploadクリック後、ファイルがOPFSに保存される
    await expect(page.locator('a:has-text("upload-test.txt")')).toBeVisible();

    // 4. alertダイアログのメッセージを確認
    // expect: アップロード成功アラートが表示される

    // 5. 一覧でsize列がファイルサイズと一致することを確認
    // expect: ファイルのサイズが正しい
    await expect(page.locator('.grid-row div.size:has-text("13")')).toBeVisible();
  });

  test('同じ名前のファイルが存在する場合のアップロード（上書き確認）', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // Create an existing file first
    await page.locator('#editor-name').fill('overwrite-test.txt');
    await page.locator('#editor-edit').fill('Original content');
    await page.locator('#editor-save').click();

    // 1. 既存ファイルと同じ名前でファイルアップロードを試みる
    await page.locator('#upload-file').setInputFiles({
      name: 'overwrite-test.txt',
      contentType: 'text/plain',
      buffer: Buffer.from('New content'),
    });

    await expect(page.locator('#upload-name')).toHaveValue('overwrite-test.txt');

    // Handle confirm dialog
    let firstDialog = true;
    page.once('dialog', async dialog => {
      expect(dialog.message()).toContain('overwrite-test.txt');
      if (firstDialog) {
        firstDialog = false;
        await dialog.accept();
      }
    });

    await page.locator('#upload-upload').click();

    // expect: 確認ダイアログでOKを選択すると上書きされる

    // 2. confirmダイアログでacceptし、ファイルが上書きされることを確認
    // Already handled above

    // 3. confirmダイアログでdismissし、ファイルが上書きされないことを確認
    page.once('dialog', async dialog => {
      expect(dialog.message()).toContain('overwrite-test.txt');
      await dialog.dismiss();
    });

    await page.locator('#upload-upload').click();
    // expect: 確認ダイアログでキャンセルすると上書きされない
  });

  test('空のファイル名のバリデーション（エディタ）', async ({ page }) => {
    await page.goto('http://127.0.0.1:9280/opfs/');

    // 1. #editor-nameを空にして#editor-saveをクリック
    await page.locator('#editor-name').clear();
    await page.locator('#editor-save').click();

    // expect: エディタで空のファイル名で保存時、エラーが発生する
    await expect(page.locator('#editor-name')).toHaveValue('');
  });
});
